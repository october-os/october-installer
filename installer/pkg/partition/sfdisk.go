package partition

import (
	"encoding/json"
	"io"
	"os/exec"
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
}

// Gets a drive's state using 'sfdisk --json <device>'
// Useful to compare the state before and after creating partitions
//
// Decodes the JSON state into a SfdiskJsonDrive object and returns it
// Can return one type of error: SetupPartitionsError
func getDriveStateWithSfdisk(drive string) (*SfdiskJsonDrive, error) {
	cmd := exec.Command("sfdisk", "--json", drive)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	var stdoutOutput []byte
	if stdoutOutput, err = io.ReadAll(stdout); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	var sjd SfdiskJsonDrive
	if err = json.Unmarshal(stdoutOutput, &sjd); err != nil {
		return nil, err
	}
	return &sjd, nil
}
