package partition

func CreatePartitioningFiles(partitions *[]partition) error {
	partitionsPerDrive := partitionsPerDrive(partitions)
	// IN PROGRESS
	return nil
}

func partitionsPerDrive(partitions *[]partition) *map[string][]partition {
	drives := make(map[string][]partition)
	for _, p := range *partitions {
		drives[p.Drive] = append(drives[p.Drive], p)
	}
	return &drives
}
