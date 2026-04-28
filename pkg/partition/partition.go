package partition

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/october-os/october-installer/pkg/utils"
)

var partitions []CreatedPartition
var commandExecutor = utils.NewCommandExecutor

// CreatePartitions creates the partitions from a list of Drive
// Can return one type of error: PartitionError
func CreatePartitions(drives []Drive) error {
	if err := checkCompatibility(drives); err != nil {
		return PartitionError{err: err}
	}
	if err := createPartitions(drives); err != nil {
		return PartitionError{err: err}
	}
	return nil
}

// FormatPartitions formats the partitions
// Must be run after CreatePartitions
// Can return one type of error: PartitionError
func FormatPartitions() error {
	if partitions == nil || len(partitions) == 0 {
		return PartitionError{err: fmt.Errorf("no created partitions")}
	}

	for _, p := range partitions {
		if err := p.format(); err != nil {
			return err
		}
	}
	return nil
}

// MountSystemPartitions mounts the system partitions
// Must be run after CreatePartitions
// Can return one type of error: PartitionError
func MountSystemPartitions() error {
	if partitions == nil || len(partitions) == 0 {
		return PartitionError{err: fmt.Errorf("no created partitions")}
	}

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
		if err := p.mount(); err != nil {
			return err
		}
	}

	return nil
}

// MountAdditionalPartitions mounts the additional partitions
// Must be run after CreatePartitions
// Can return one type of error: PartitionError
func MountAdditionalPartitions() error {
	if partitions == nil || len(partitions) == 0 {
		return PartitionError{err: fmt.Errorf("no created partitions")}
	}

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
		if err := p.mount(); err != nil {
			return err
		}
	}

	return nil
}

// checkCompatibility Checks the compatibility of a list of Drive
// A drive needs the GPT partition table to be compatible
func checkCompatibility(drives []Drive) error {
	for _, drive := range drives {
		if err := drive.isCompatible(); err != nil {
			return err
		}
	}
	return nil
}

// createPartitions create partitions from a list of Drive using sfdisk
func createPartitions(drives []Drive) error {
	partitioningFiles, err := createPartitioningFiles(drives)
	if err != nil {
		return err
	}

	var createdPartitions []CreatedPartition

	for drive, fileName := range partitioningFiles {
		var initialState *SfdiskJsonDrive
		if drive.Append {
			initialState, err = getDriveStateWithSfdisk(drive.Path)
			if err != nil {
				return err
			}
		}

		if err := createPartitionsFromFile(*drive, fileName); err != nil {
			return err
		}

		stateAfterCreatingPartitions, err := getDriveStateWithSfdisk(drive.Path)
		if err != nil {
			return err
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

	partitions = createdPartitions

	return nil
}

// createPartitionsFromFile creates partitions from a sfdisk named-fields syntax file with sfdisk
func createPartitionsFromFile(drive Drive, fileName string) error {
	sfdiskCommand := ""
	if drive.Append {
		sfdiskCommand = fmt.Sprintf("sfdisk -a %s < %s", drive.Path, fileName)
	} else {
		sfdiskCommand = fmt.Sprintf("sfdisk %s < %s", drive.Path, fileName)
	}

	cmd := commandExecutor("/bin/bash", "-c", sfdiskCommand)

	return cmd.Run()
}

// createPartitioningFiles creates one file per drive containing its partitions in sfdisk named-fields syntax
// from a list of Drive
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
