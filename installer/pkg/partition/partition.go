package partition

import (
	"fmt"
	"os"
	"strings"
)

func createPartitioningFiles(partitions *[]Partition) (map[string]string, error) {
	partitionsPerDrive := partitionsPerDrive(partitions)
	drivePartitionsFiles := make(map[string]string)
	for drive, partitions := range *partitionsPerDrive {
		fileName := strings.ReplaceAll(drive, "/", "")
		fullFilename := fmt.Sprintf("%s.txt", fileName)
		drivePartitionsFiles[drive] = fullFilename
		file, err := os.Create(fullFilename)
		if err != nil {
			errorMessage := fmt.Sprintf("could not create file '%s' for '%s' drive partitioning", fullFilename, drive)
			return nil, &CreatePartitioningFileError{
				Message: errorMessage,
				Err:     err,
			}
		}
		defer file.Close()
		for _, p := range partitions {
			partitionEntry := fmt.Sprintf("%s\n", p.ToSfdiskFormat())
			_, err := file.WriteString(partitionEntry)
			if err != nil {
				errorMessage := fmt.Sprintf("could not edit file '%s' for partitioning", fullFilename)
				return nil, &CreatePartitioningFileError{
					Message: errorMessage,
					Err:     err,
				}
			}
		}
	}
	return drivePartitionsFiles, nil
}

// oublie pas /bin/bash -c -BAE

func doPartiotnafef() {
	// use createParitinFiles
	// iterate through
	// do : sfdisk /dev/sda [key] < partiton.txt [value]
}

func partitionsPerDrive(partitions *[]Partition) *map[string][]Partition {
	drives := make(map[string][]Partition)
	for _, p := range *partitions {
		drives[p.Drive] = append(drives[p.Drive], p)
	}
	return &drives
}
