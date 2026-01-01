package partition

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	fileSystemExt4  string = "ext4"
	fileSystemBtrfs string = "btrfs"
)

var supportedFileSystems []string = []string{
	fileSystemExt4,
	fileSystemBtrfs,
}

const (
	gptPartitionTypeEfi        string = "C12A7328-F81F-11D2-BA4B-00A0C93EC93B"
	gptPartitionTypeSwap       string = "0657FD6D-A4AB-43C4-84E5-0933C84B4F4F"
	gptPartitionTypeRoot       string = "4F68BCE3-E8CD-4DB1-96E7-FBCAF984B709"
	gptPartitionTypeFileSystem string = "0FC63DAF-8483-4772-8E79-3D69D8477DE4"
	gptPartitionTypeHome       string = "933AC7E1-2EB4-4F13-B844-0E14E2AEF915"
)

var supportedGptPartitionTypes []string = []string{
	gptPartitionTypeEfi,
	gptPartitionTypeSwap,
	gptPartitionTypeRoot,
	gptPartitionTypeFileSystem,
	gptPartitionTypeHome,
}

const (
	partitionSizeUnitKiB string = "KiB"
	partitionSizeUnitMiB string = "MiB"
	partitionSizeUnitGiB string = "GiB"
	partitionSizeUnitTiB string = "TiB"
	partitionSizeUnitPiB string = "PiB"
	partitionSizeUnitEiB string = "EiB"
	partitionSizeUnitZiB string = "ZiB"
	partitionSizeUnitYiB string = "YiB"
)

var supportedPartitionSizeUnits []string = []string{
	partitionSizeUnitKiB,
	partitionSizeUnitMiB,
	partitionSizeUnitGiB,
	partitionSizeUnitTiB,
	partitionSizeUnitPiB,
	partitionSizeUnitEiB,
	partitionSizeUnitZiB,
	partitionSizeUnitYiB,
}

// Drive represents a drive that needs to have partitions added to it
// Possible attributes values:
// - Path: the full path of to drive (starting with '/dev/')
type Drive struct {
	Path       string      `json:"path"`
	Append     bool        `json:"append"`
	Partitions []Partition `json:"partitions"`
}

// Validates the attributes of a Drive struct
// Returns a ValidationError if validation fails
func (d *Drive) Validate() error {
	if !strings.HasPrefix(d.Path, "/dev/") {
		return &ValidationError{
			Err: errors.New("Drive validation: error=Path is in the wrong format: should start by '/dev/'"),
		}
	}
	return nil
}

// Partition represents a drive/disk partition that needs to be created
// Possible attributes values:
// - FileSystem: A file system present in the supportedFileSystems slice above, or default string value
// - PartitionType: a GPT partition type present in the supportedGptPartitionTypes slice above
// - MountPoint: an absolute Linux filesystem path, or string default value
type Partition struct {
	Size          PartitionSize `json:"size"`
	FileSystem    string        `json:"fileSystem"`
	PartitionType string        `json:"partitionType"`
	MountPoint    string        `json:"mountPoint"`
}

// Transforms a partition into its sfdisk format
// Returns a string
//
// Example:
// "type=C12A7328-F81F-11D2-BA4B-00A0C93EC93B, size=1GiB"
func (p *Partition) toSfdiskFormat() string {
	partition_string := fmt.Sprintf("type=%s", p.PartitionType)
	if p.Size.TakeRemaining {
		partition_string += ", size=+"
	} else {
		partition_string += fmt.Sprintf(", size=%d%s", p.Size.Amount, p.Size.Unit)
	}
	return partition_string
}

// Validates the attributes of a Partition struct
// Returns a ValidationError if validation fails
func (p *Partition) Validate() error {
	if p.MountPoint != "" {
		if !strings.HasPrefix(p.MountPoint, "/") {
			return &ValidationError{
				Err: errors.New("Partition validation: error=MountPoint is in the wrong format: should start by '/'"),
			}
		}
	}
	if !slices.Contains(supportedGptPartitionTypes, p.PartitionType) {
		return &ValidationError{
			Err: errors.New("Partition validation: error=specified PartitionType is not supported"),
		}
	}
	if p.FileSystem != "" {
		if !slices.Contains(supportedFileSystems, p.FileSystem) {
			return &ValidationError{
				Err: errors.New("Partition validation: error=specified FileSystem is not supported"),
			}
		}
	}

	if p.FileSystem == "" {
		if p.PartitionType != gptPartitionTypeEfi && p.PartitionType != gptPartitionTypeSwap {
			return &ValidationError{
				Err: errors.New("Partition validation: error=Filesystem is not defined, but the partition type needs a file system"),
			}
		}
	}

	if p.MountPoint == "" {
		if p.PartitionType == gptPartitionTypeEfi || p.PartitionType == gptPartitionTypeSwap || p.PartitionType == gptPartitionTypeRoot {
			return &ValidationError{
				Err: errors.New("Partition validation: error=MountPoint is not defined, but the partition type needs a mount point"),
			}
		}
	}

	return nil
}

// PartitionSize represents the size of a Partition
// Possible attributes values:
// Amount: any positive integer greater or equal 1, or int default value
// Unit: a partition size unit present in the supportedPartitionSizeUnits slice above, or string default value
// TakeRemaining: true/false, if false: Amount and Unit must not be default int/string values
type PartitionSize struct {
	Amount        int    `json:"amount"`
	Unit          string `json:"unit"`
	TakeRemaining bool   `json:"takeRemaining"`
}

// Validates the attributes of a PartitionSize struct
// Returns a ValidationError if validation fails
func (p *PartitionSize) Validate() error {
	if p.TakeRemaining == false && (p.Amount == 0 || p.Unit == "") {
		return &ValidationError{
			Err: errors.New("PartitionSize validation: error=TakeRemaining is false but Amount and/or Unit are not defined"),
		}
	}

	if p.Amount != 0 {
		if p.Amount < 1 {
			return &ValidationError{
				Err: errors.New("PartitionSize validation: error=Amount must be greater or equal 1"),
			}
		}
	}

	if p.Unit != "" {
		if !slices.Contains(supportedPartitionSizeUnits, p.Unit) {
			return &ValidationError{
				Err: errors.New("PartitionSize validation: error=specified Unit is not supported"),
			}
		}
	}

	return nil
}
