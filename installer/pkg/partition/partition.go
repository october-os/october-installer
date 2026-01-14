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
// Can return one type of error: PartitionError
func SetupPartitions(drives []Drive) error {
	if err := checkCompatibility(drives); err != nil {
		return &PartitionError{err: err}
	}
	newPartitionsMappings, err := createPartitions(drives)
	if err != nil {
		return &PartitionError{err: err}
	}

	for _, mapping := range newPartitionsMappings {
		for partition, sfdiskPartition := range mapping {
			if err = formatPartition(partition, sfdiskPartition.Node); err != nil {
				return &PartitionError{err: err}
			}
			if err = mountPartition(partition, sfdiskPartition.Node); err != nil {
				return &PartitionError{err: err}
			}
		}
	}

	return nil
}

// Checks the compatibility of a list of Drives
// A drive needs the GPT partition table to be compatible
func checkCompatibility(drives []Drive) error {
	for _, drive := range drives {
		cmd := exec.Command("lsblk", drive.Path, "-dno", "pttype")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		if err := cmd.Start(); err != nil {
			return err
		}
		var stdoutOutput []byte
		if stdoutOutput, err = io.ReadAll(stdout); err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}
		if string(stdoutOutput) != "gpt\n" {
			return fmt.Errorf("drive '%s' is not compatible: partition table must be GPT", drive.Path)
		}
	}
	return nil
}

// Create Partitions from a list of Drives using sfdisk
//
// Returns a mapping of the Partition and its corresponding SfdiskJsonPartition
// to map the Partition object to the partition created on the system
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
				return nil, err
			}
			sfdiskCommand = fmt.Sprintf("sfdisk -a %s < %s", drive.Path, fileName)
		} else {
			sfdiskCommand = fmt.Sprintf("sfdisk %s < %s", drive.Path, fileName)
		}

		cmd := exec.Command("/bin/bash", "-c", sfdiskCommand)

		if err := cmd.Run(); err != nil {
			return nil, err
		}

		stateAfterCreatingPartitions, err := getDriveStateWithSfdisk(drive.Path)
		var newPartitions []SfdiskJsonPartition
		if err != nil {
			return nil, err
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
//
// Returns a map of the drives and their files name
func createPartitioningFiles(drives []Drive) (map[*Drive]string, error) {
	drivePartitionsFiles := make(map[*Drive]string)
	for _, drive := range drives {
		fileName := strings.ReplaceAll(drive.Path, "/", "")
		drivePartitionsFiles[&drive] = fileName

		file, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		for _, partition := range drive.Partitions {
			partitionEntry := fmt.Sprintf("%s\n", partition.toSfdiskFormat())
			_, err := file.WriteString(partitionEntry)
			if err != nil {
				return nil, err
			}
		}
	}
	return drivePartitionsFiles, nil
}

// Formats a partition
func formatPartition(partition Partition, path string) error {
	cmd, err := partition.formatCommand(path)
	if err != nil {
		return err
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Mounts a partition
func mountPartition(partition Partition, path string) error {
	cmd, err := partition.mountCommand(path)
	if err != nil {
		return err
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
