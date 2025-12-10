package partition

import (
	"fmt"
	"strings"
)

type GptPartitionType string

const (
	GptPartitionTypeEfi        GptPartitionType = "C12A7328-F81F-11D2-BA4B-00A0C93EC93B"
	GptPartitionTypeSwap       GptPartitionType = "0657FD6D-A4AB-43C4-84E5-0933C84B4F4F"
	GptPartitionTypeRoot       GptPartitionType = "4F68BCE3-E8CD-4DB1-96E7-FBCAF984B709"
	GptPartitionTypeFileSystem GptPartitionType = "0FC63DAF-8483-4772-8E79-3D69D8477DE4"
	GptPartitionTypeHome       GptPartitionType = "933AC7E1-2EB4-4F13-B844-0E14E2AEF915"
)

type PartitionSizeUnit string

const (
	PartitionSizeUnitKiB PartitionSizeUnit = "KiB"
	PartitionSizeUnitMiB PartitionSizeUnit = "MiB"
	PartitionSizeUnitGiB PartitionSizeUnit = "GiB"
	PartitionSizeUnitTiB PartitionSizeUnit = "TiB"
	PartitionSizeUnitPiB PartitionSizeUnit = "PiB"
	PartitionSizeUnitEiB PartitionSizeUnit = "EiB"
	PartitionSizeUnitZiB PartitionSizeUnit = "ZiB"
	PartitionSizeUnitYiB PartitionSizeUnit = "YiB"
)

type Partition struct {
	Drive         string
	Size          PartitionSize
	PartitionType GptPartitionType
	MountPoint    *string
}

func (p Partition) ToSfdiskFormat() string {
	partition_string := fmt.Sprintf("uuid=%s", p.PartitionType)

	if p.Size.TakeRemaining == true {
		partition_string += ", size=+"
	} else {
		partition_string += fmt.Sprintf(", size=%d%s", p.Size.Amount, *p.Size.Unit)
	}
	return partition_string
}

func NewPartition(drive string, size *PartitionSize, partitionType GptPartitionType, mountPoint *string) (*Partition, error) {
	if strings.HasPrefix(drive, "/dev/") == false {
		return nil, NewPartitionError{
			Err: "parameter drive is in the wrong format: should start by '/dev/'",
		}
	}

	if strings.HasPrefix(*mountPoint, "/") == false {
		return nil, NewPartitionError{
			Err: "parameter mountPoint is in the wrong format: should start by '/'",
		}
	}

	return &Partition{
		Drive:         drive,
		Size:          *size,
		PartitionType: partitionType,
		MountPoint:    mountPoint,
	}, nil
}

type PartitionSize struct {
	Amount        *int
	Unit          *string
	TakeRemaining bool
}

func NewPartitionSize(amount *int, unit *string, takeRemaining *bool) (*PartitionSize, error) {
	var takeRemainingValue bool
	if takeRemaining == nil {
		takeRemainingValue = false
	} else {
		takeRemainingValue = *takeRemaining
	}

	if takeRemainingValue == false && (amount == nil || unit == nil) {
		return nil, &NewPartitionSizeError{
			Err: "parameter takeRemaining is false but parameters amount and/or unit are not specified.",
		}
	}

	return &PartitionSize{
		Amount:        amount,
		Unit:          unit,
		TakeRemaining: takeRemainingValue,
	}, nil
}
