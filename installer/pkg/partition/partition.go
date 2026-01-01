package partition

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Sets up the partitions for a list of Drive:
// 1. Checks compatibility
// 2. Creates the partitions
// 3. Formats and mounts each partition
//
// Can return one type of error: SetupPartitionsError
// if:
// The compatibility check failed/found an incompatibility
// or
// The partitions creation failed
// or
// A partition's formatting failed
// or
// A partition's mounting failed
func SetupPartitions(drives []Drive) error {
	if err := checkCompatibility(drives); err != nil {
		return err
	}
	newPartitionsMappings, err := createPartitions(drives)
	if err != nil {
		return err
	}

	for _, mapping := range newPartitionsMappings {
		for partition, sfdiskPartition := range mapping {
			if err = formatPartition(partition, sfdiskPartition.Node); err != nil {
				return err
			}
			if err = mountPartition(partition, sfdiskPartition.Node); err != nil {
				return err
			}
		}
	}

	return nil
}

// Checks the compatibility of a list of Drives
// A drive needs the GPT partition table to be compatible
//
// Can return one type of error: SetupPartitionsError
// if:
// stdout couldn't be piped
// or
// stderr couldn't be piped
// or
// the partition table type couldn't be fetched using lsblk
// or
// stdout couldn't be read
// or
// a drive is not compatible
func checkCompatibility(drives []Drive) error {
	for _, drive := range drives {
		cmd := exec.Command("lsblk", drive.Path, "-dno", "pttype")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return &SetupPartitionsError{
				Err: fmt.Errorf("error piping stdout: error=%s", err.Error()),
			}
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return &SetupPartitionsError{
				Err: fmt.Errorf("error piping stderr: error=%s", err.Error()),
			}
		}
		if err := cmd.Start(); err != nil {
			stderrOutput, _ := io.ReadAll(stderr)
			return &SetupPartitionsError{
				Err: fmt.Errorf("error getting partition table type for drive '%s': error=%s", drive.Path, string(stderrOutput)),
			}
		}
		var stdoutOutput []byte
		if stdoutOutput, err = io.ReadAll(stdout); err != nil {
			return &SetupPartitionsError{
				Err: fmt.Errorf("error reading stdout: error=%s", err.Error()),
			}
		}
		if err := cmd.Wait(); err != nil {
			return &SetupPartitionsError{
				Err: fmt.Errorf("error reading stdout: error=%s", err.Error()),
			}
		}
		if string(stdoutOutput) != "gpt\n" {
			return &SetupPartitionsError{
				Err: fmt.Errorf("drive '%s' is not compatible: partition table must be GPT", drive.Path),
			}
		}
	}
	return nil
}

// Create Partitions from a list of Drives using sfdisk
//
// Returns a mapping of the Partition and its corresponding SfdiskJsonPartition
// to map the Partition object to the partition created on the system
//
// Can return one type of error:
//   - SetupPartitionsError:
//     when the creation of partitioning files failed
//     or
//     when the creation of partitions failed using sfdisk
//     or
//     stderr couldn't be piped
//     or
//     getting drive's states using sfdisk
func createPartitions(drives []Drive) ([]map[Partition]SfdiskJsonPartition, error) {
	partitioningFiles, err := createPartitioningFiles(drives)
	if err != nil {
		return nil, err
	}

	var mappings []map[Partition]SfdiskJsonPartition

	for drive, fileName := range partitioningFiles {
		sfdiskCommand := ""
		var initialState *SfdiskJsonDrive

		if drive.Append {
			initialState, err = getDriveStateWithSfdisk(drive.Path)
			if err != nil {
				return nil, &SetupPartitionsError{
					Err: fmt.Errorf("error getting initial state of drive '%s': error=%s", drive.Path, err.Error()),
				}
			}
			sfdiskCommand = fmt.Sprintf("sfdisk -a %s < %s", drive.Path, fileName)
		} else {
			sfdiskCommand = fmt.Sprintf("sfdisk %s < %s", drive.Path, fileName)
		}

		cmd := exec.Command("/bin/bash", "-c", sfdiskCommand)
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return nil, &SetupPartitionsError{
				Err: fmt.Errorf("error piping stderr: error=%s", err.Error()),
			}
		}

		if err := cmd.Run(); err != nil {
			stderrOutput, _ := io.ReadAll(stderr)
			return nil, &SetupPartitionsError{
				Err: fmt.Errorf("error creating partitions on drive '%s' with file '%s' using sfdisk: error=%s", drive.Path, fileName, string(stderrOutput)),
			}
		}

		stateAfterCreatingPartitions, err := getDriveStateWithSfdisk(drive.Path)
		var newPartitions []SfdiskJsonPartition
		if err != nil {
			return nil, &SetupPartitionsError{
				Err: fmt.Errorf("error getting state after partitions creation of drive '%s': error=%s", drive.Path, err.Error()),
			}
		}
		if initialState != nil {
			newPartitions = stateAfterCreatingPartitions.PartitionTable.Partitions[len(initialState.PartitionTable.Partitions):]
		} else {
			newPartitions = stateAfterCreatingPartitions.PartitionTable.Partitions
		}

		partitionsMap := make(map[Partition]SfdiskJsonPartition)
		for i := 0; i < len(newPartitions) || i < len(drive.Partitions); i++ {
			partitionsMap[drive.Partitions[i]] = newPartitions[i]
		}

		mappings = append(mappings, partitionsMap)
	}

	return mappings, nil
}

// Creates one file per drive containing its partitions in sfdisk named-fields syntax
// from a list of Drives
// Returns a map of the drives and their files name
// Can return one type of error:
//   - SetupPartitionsError:
//     when a file couldn't be created
//     or
//     when a file couldn't be modified
func createPartitioningFiles(drives []Drive) (map[*Drive]string, error) {
	drivePartitionsFiles := make(map[*Drive]string)
	for _, drive := range drives {
		fileName := strings.ReplaceAll(drive.Path, "/", "")
		drivePartitionsFiles[&drive] = fileName

		file, err := os.Create(fileName)
		if err != nil {
			return nil, &SetupPartitionsError{
				Err: fmt.Errorf("could not create file '%s' for '%s' drive partitioning: error=%s", fileName, drive.Path, err.Error()),
			}
		}
		defer file.Close()

		for _, partition := range drive.Partitions {
			partitionEntry := fmt.Sprintf("%s\n", partition.toSfdiskFormat())
			_, err := file.WriteString(partitionEntry)
			if err != nil {
				return nil, &SetupPartitionsError{
					Err: fmt.Errorf("could not edit file '%s' for partitioning: error=%s", fileName, err.Error()),
				}
			}
		}
	}
	return drivePartitionsFiles, nil
}

// Formats a partition according to its settings
//
// Can return one type of error: SetupPartitionsError
// if:
// cmd is nil; a command couldn't be formed
// or
// the formatting failed
func formatPartition(partition Partition, path string) error {
	var cmd *exec.Cmd
	switch partition.PartitionType {
	case gptPartitionTypeEfi:
		cmd = exec.Command("mkfs.fat", "-F", "32", path)
	case gptPartitionTypeSwap:
		cmd = exec.Command("mkswap", path)
	case gptPartitionTypeRoot, gptPartitionTypeHome, gptPartitionTypeFileSystem:
		switch partition.FileSystem {
		case fileSystemExt4:
			cmd = exec.Command("mkfs.ext4", path)
		case fileSystemBtrfs:
			cmd = exec.Command("mkfs.btrfs", path)
		}
	}

	if cmd == nil {
		return &SetupPartitionsError{
			Err: fmt.Errorf("error formatting partition '%s': cmd is nil", path),
		}
	}

	if err := cmd.Run(); err != nil {
		return &SetupPartitionsError{
			Err: fmt.Errorf("error formatting partition '%s': error=%s", path, err.Error()),
		}
	}

	return nil
}

// Mounts a partition according to its settings
//
// Can return one type of error: SetupPartitionsError
// if:
// cmd is nil; a command couldn't be formed
// or
// the mounting failed
func mountPartition(partition Partition, path string) error {
	var cmd *exec.Cmd
	switch partition.PartitionType {
	case gptPartitionTypeEfi:
		cmd = exec.Command("mount", "--mkdir", path, "/mnt/boot")
	case gptPartitionTypeSwap:
		cmd = exec.Command("swapon", path)
	case gptPartitionTypeRoot:
		cmd = exec.Command("mount", path, "/mnt")
	case gptPartitionTypeHome, gptPartitionTypeFileSystem:
		cmd = exec.Command("mount", "--mkdir", path, partition.MountPoint)
	}

	if cmd == nil {
		return &SetupPartitionsError{
			Err: fmt.Errorf("error mounting partition %s", path),
		}
	}

	if err := cmd.Run(); err != nil {
		return &SetupPartitionsError{
			Err: fmt.Errorf("error mounting partition '%s': error=%s", path, err.Error()),
		}
	}

	return nil
}
