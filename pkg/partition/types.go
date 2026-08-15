package partition

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"

	"github.com/october-os/october-installer/pkg/utils"
)

const (
	fileSystemExt4  string = "ext4"
	fileSystemBtrfs string = "btrfs"
)

var supportedFileSystems = []string{
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

var supportedGptPartitionTypes = []string{
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

var supportedPartitionSizeUnits = []string{
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

func (d *Drive) isCompatible() error {
	cmd := commandExecutor("lsblk", d.Path, "-dno", "pttype")
	stdoutOutput, err := cmd.Output()
	if err != nil {
		return err
	}
	if string(stdoutOutput) != "gpt\n" {
		return fmt.Errorf("drive '%s' is not compatible: partition table must be GPT", d.Path)
	}
	return nil
}

// Validate validates the attributes of a Drive struct
// Returns a PartitionError if validation fails
func (d *Drive) Validate() error {
	if !strings.HasPrefix(d.Path, "/dev/") {
		return PartitionError{
			err: errors.New("Drive validation: error=Path is in the wrong format: should start by '/dev/'"),
		}
	}
	for _, partition := range d.Partitions {
		if err := partition.Validate(); err != nil {
			return err
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

// toSfdiskFormat transforms a partition into its sfdisk format
// Returns a string
//
// Example:
// "type=C12A7328-F81F-11D2-BA4B-00A0C93EC93B, size=1GiB"
func (p *Partition) toSfdiskFormat() string {
	partitionString := fmt.Sprintf("type=%s", p.PartitionType)
	if p.Size.TakeRemaining {
		partitionString += ", size=+"
	} else {
		partitionString += fmt.Sprintf(", size=%d%s", p.Size.Amount, p.Size.Unit)
	}
	return partitionString
}

// isSystemPartition determines whether the partition is a system partition or not
func (p *Partition) isSystemPartition() bool {
	return p.PartitionType == gptPartitionTypeEfi || p.PartitionType == gptPartitionTypeSwap || p.PartitionType == gptPartitionTypeRoot
}

// Validate validates the attributes of a Partition struct
// Returns a PartitionError if validation fails
func (p *Partition) Validate() error {
	if p.MountPoint != "" {
		if !strings.HasPrefix(p.MountPoint, "/") {
			return PartitionError{
				err: errors.New("Partition validation: error=MountPoint is in the wrong format: should start by '/'"),
			}
		}
	}
	if !slices.Contains(supportedGptPartitionTypes, p.PartitionType) {
		return PartitionError{
			err: errors.New("Partition validation: error=specified PartitionType is not supported"),
		}
	}
	if p.FileSystem != "" {
		if !slices.Contains(supportedFileSystems, p.FileSystem) {
			return PartitionError{
				err: errors.New("Partition validation: error=specified FileSystem is not supported"),
			}
		}
	}

	if p.FileSystem == "" {
		if p.PartitionType != gptPartitionTypeEfi && p.PartitionType != gptPartitionTypeSwap {
			return PartitionError{
				err: errors.New("Partition validation: error=Filesystem is not defined, but the partition type needs a file system"),
			}
		}
	}

	if p.MountPoint == "" {
		if p.PartitionType != gptPartitionTypeEfi &&
			p.PartitionType != gptPartitionTypeSwap &&
			p.PartitionType != gptPartitionTypeRoot &&
			p.PartitionType != gptPartitionTypeHome {
			return PartitionError{
				err: errors.New("Partition validation: error=MountPoint is not defined, but the partition type needs a mount point"),
			}
		}
	}

	return p.Size.Validate()
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

// Validate validates the attributes of a PartitionSize struct
// Returns a PartitionError if validation fails
func (p *PartitionSize) Validate() error {
	if p.TakeRemaining == false && (p.Amount == 0 || p.Unit == "") {
		return PartitionError{
			err: errors.New("PartitionSize validation: error=TakeRemaining is false but Amount and/or Unit are not defined"),
		}
	}

	if p.Amount != 0 {
		if p.Amount < 1 {
			return PartitionError{
				err: errors.New("PartitionSize validation: error=Amount must be greater or equal 1"),
			}
		}
	}

	if p.Unit != "" {
		if !slices.Contains(supportedPartitionSizeUnits, p.Unit) {
			return PartitionError{
				err: errors.New("PartitionSize validation: error=specified Unit is not supported"),
			}
		}
	}

	return nil
}

// CreatedPartition represents a Partition that has been created on the system.
// Attributes:
// Partition: the Partition that was created
// SfdiskJsonPartition: the SfdiskJsonPartition that represents the createdPartition on the system
type CreatedPartition struct {
	Partition           Partition
	SfdiskJsonPartition SfdiskJsonPartition
}

// format formats the partition
func (p *CreatedPartition) format() error {
	var cmd utils.ICommandExecutor

	switch p.Partition.PartitionType {
	case gptPartitionTypeEfi:
		cmd = commandExecutor("mkfs.fat", "-F", "32", p.SfdiskJsonPartition.Node)
	case gptPartitionTypeSwap:
		cmd = commandExecutor("mkswap", p.SfdiskJsonPartition.Node)
	case gptPartitionTypeRoot, gptPartitionTypeHome, gptPartitionTypeFileSystem:
		switch p.Partition.FileSystem {
		case fileSystemExt4:
			cmd = commandExecutor("mkfs.ext4", "-F", p.SfdiskJsonPartition.Node)
		case fileSystemBtrfs:
			cmd = commandExecutor("mkfs.btrfs", "-f", p.SfdiskJsonPartition.Node)
		}
	}

	if cmd == nil {
		return errors.New("error choosing a formatting command: unsupported file system or partition type")
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// mount Mounts the partition
func (p *CreatedPartition) mount() error {
	switch p.Partition.PartitionType {
	case gptPartitionTypeEfi:
		if err := os.MkdirAll("/mnt/boot", 0755); err != nil {
			return fmt.Errorf("mkdir /mnt/boot %s", err.Error())
		}

		return syscall.Mount(p.SfdiskJsonPartition.Node, "/mnt/boot", "vfat", 0, "")
	case gptPartitionTypeSwap:
		cmd := exec.Command("swapon", p.SfdiskJsonPartition.Node)
		return cmd.Run()
	case gptPartitionTypeRoot:
		return syscall.Mount(p.SfdiskJsonPartition.Node, "/mnt", p.Partition.FileSystem, 0, "")
	case gptPartitionTypeHome:
		return syscall.Mount(p.SfdiskJsonPartition.Node, "/mnt/home", p.Partition.FileSystem, 0, "")
	case gptPartitionTypeFileSystem:
		if err := os.MkdirAll(fmt.Sprintf("/mnt%s", p.Partition.MountPoint), 0755); err != nil {
			return fmt.Errorf("mkdir %s %s", p.Partition.MountPoint, err.Error())
		}
		return syscall.Mount(p.SfdiskJsonPartition.Node, fmt.Sprintf("/mnt%s", p.Partition.MountPoint), p.Partition.FileSystem, 0, "")
	}

	return errors.New("error choosing a mounting syscall: unsupported partition type")
}
