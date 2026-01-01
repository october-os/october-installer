package partition

import (
	"encoding/json"
	"fmt"
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
// # Decodes the JSON state into a SfdiskJsonDrive object and returns it
//
// Can return one type of error: SetupPartitionsError
// if:
// stdout couldn't be piped
// or
// stderr couldn't be piped
// or
// the state couldn't be fetched using sfdisk
// or
// stdout couldn't be read
// or
// JSON decoding failed
func getDriveStateWithSfdisk(drive string) (*SfdiskJsonDrive, error) {
	cmd := exec.Command("sfdisk", "--json", drive)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, &SetupPartitionsError{
			Err: fmt.Errorf("error piping stdout: error=%s", err.Error()),
		}
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, &SetupPartitionsError{
			Err: fmt.Errorf("error piping stderr: error=%s", err.Error()),
		}
	}
	if err := cmd.Start(); err != nil {
		stderrOutput, _ := io.ReadAll(stderr)
		return nil, &SetupPartitionsError{
			Err: fmt.Errorf("error getting drive state as JSON using sfdisk: error=%s", string(stderrOutput)),
		}
	}
	var stdoutOutput []byte
	if stdoutOutput, err = io.ReadAll(stdout); err != nil {
		return nil, &SetupPartitionsError{
			Err: fmt.Errorf("error reading stdout: error=%s", err.Error()),
		}
	}
	if err := cmd.Wait(); err != nil {
		return nil, &SetupPartitionsError{
			Err: fmt.Errorf("error reading stdout: error=%s", err.Error()),
		}
	}
	var sjd SfdiskJsonDrive
	if err = json.Unmarshal(stdoutOutput, &sjd); err != nil {
		return nil, &SetupPartitionsError{
			Err: fmt.Errorf("error decoding JSON drive state coming from stdout of sfdisk: error=%s", err.Error()),
		}
	}
	return &sjd, nil
}
