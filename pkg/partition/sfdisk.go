package partition

import (
	"encoding/json"
)

// SfdiskJsonDrive represents the JSON output from 'sfdisk --json <device>'
type SfdiskJsonDrive struct {
	PartitionTable SfdiskJsonPartitionTable `json:"partitiontable"`
}

// SfdiskJsonPartitionTable represents the 'partitiontable' field of SfdiskJsonDrive
type SfdiskJsonPartitionTable struct {
	Device     string                `json:"device"`
	Partitions []SfdiskJsonPartition `json:"partitions"`
}

// SfdiskJsonPartition represents one element of the 'partitions' field/array of SfdiskJsonPartitionTable
type SfdiskJsonPartition struct {
	Node string `json:"node"`
	Type string `json:"type"`
}

// Gets a drive's state using 'sfdisk --json <device>'
// Useful to compare the state before and after creating partitions
//
// Decodes the JSON state into a SfdiskJsonDrive object and returns it
func getDriveStateWithSfdisk(drive string) (*SfdiskJsonDrive, error) {
	cmd := commandExecutor("sfdisk", "--json", drive)
	stdoutOutput, err := cmd.Output()
	var sjd SfdiskJsonDrive
	if err = json.Unmarshal(stdoutOutput, &sjd); err != nil {
		return nil, err
	}
	return &sjd, nil
}
