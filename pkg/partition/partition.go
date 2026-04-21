package partition

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"
)

func CreatePartitions(drives []Drive) ([]CreatedPartition, error) {
	cleanUpMounts()
	if err := checkCompatibility(drives); err != nil {
		return nil, PartitionError{err: err}
	}
	createdPartitions, err := createPartitions(drives)
	if err != nil {
		return nil, PartitionError{err: err}
	}
	return createdPartitions, nil
}

func FormatPartitions(partitions []CreatedPartition) error {
	for _, p := range partitions {
		if err := formatPartition(p.Partition, p.SfdiskJsonPartition.Node); err != nil {
			return err
		}
	}
	return nil
}

func MountSystemPartitions(partitions []CreatedPartition) error {
	var systemPartitions []CreatedPartition
	for _, p := range partitions {
		if p.Partition.isSystemPartition() {
			systemPartitions = append(systemPartitions, p)
		}
	}

	// only pushes root partition to front since its the only partition that really needs to be mounted first
	slices.SortFunc(systemPartitions, func(a, b CreatedPartition) int {
		if a.Partition.PartitionType == gptPartitionTypeRoot && b.Partition.PartitionType != gptPartitionTypeRoot {
			return -1
		}
		if a.Partition.PartitionType != gptPartitionTypeRoot && b.Partition.PartitionType == gptPartitionTypeRoot {
			return 1
		}
		return 0
	})

	for _, p := range systemPartitions {
		if err := mountPartition(p.Partition, p.SfdiskJsonPartition.Node); err != nil {
			return err
		}
	}

	return nil
}

func MountAdditionalPartitions(partitions []CreatedPartition) error {
	var additionalPartitions []CreatedPartition
	for _, p := range partitions {
		if !p.Partition.isSystemPartition() {
			additionalPartitions = append(additionalPartitions, p)
		}
	}
	slices.SortFunc(additionalPartitions, func(a, b CreatedPartition) int {
		if strings.HasPrefix(a.Partition.MountPoint, b.Partition.MountPoint) {
			return 1
		}
		if strings.HasPrefix(b.Partition.MountPoint, a.Partition.MountPoint) {
			return -1
		}
		return 0
	})

	for _, p := range additionalPartitions {
		if err := mountPartition(p.Partition, p.SfdiskJsonPartition.Node); err != nil {
			return err
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
func createPartitions(drives []Drive) ([]CreatedPartition, error) {
	partitioningFiles, err := createPartitioningFiles(drives)
	if err != nil {
		return nil, err
	}

	var createdPartitions []CreatedPartition

	for drive, fileName := range partitioningFiles {
		var initialState *SfdiskJsonDrive
		if drive.Append {
			initialState, err = getDriveStateWithSfdisk(drive.Path)
			if err != nil {
				return nil, err
			}
		}

		if err := createPartitionsFromFile(*drive, fileName); err != nil {
			return nil, err
		}

		stateAfterCreatingPartitions, err := getDriveStateWithSfdisk(drive.Path)
		if err != nil {
			return nil, err
		}

		var newPartitions []SfdiskJsonPartition
		if initialState != nil {
			newPartitions = stateAfterCreatingPartitions.PartitionTable.Partitions[len(initialState.PartitionTable.Partitions):]
		} else {
			newPartitions = stateAfterCreatingPartitions.PartitionTable.Partitions
		}

		for i := 0; i < len(newPartitions) || i < len(drive.Partitions); i++ {
			createdPartitions = append(createdPartitions, CreatedPartition{Partition: drive.Partitions[i], SfdiskJsonPartition: newPartitions[i]})
		}
	}

	return createdPartitions, nil
}

func createPartitionsFromFile(drive Drive, fileName string) error {
	sfdiskCommand := ""
	if drive.Append {
		sfdiskCommand = fmt.Sprintf("sfdisk -a %s < %s", drive.Path, fileName)
	} else {
		sfdiskCommand = fmt.Sprintf("sfdisk %s < %s", drive.Path, fileName)
	}

	cmd := exec.Command("/bin/bash", "-c", sfdiskCommand)

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
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
	err := partition.mount(path)
	if err != nil {
		return err
	}

	return nil
}

func cleanUpMounts() {
	syscall.Unmount("/mnt/boot", syscall.MNT_DETACH)
	syscall.Unmount("/mnt", syscall.MNT_DETACH)
}
