## Drives

Drives contains an array of "drives" and those contain partitions. This section dictates what partitions to create on which drives.

- [Drive configuration example](#Drive-configuration-example)
- [Supported partition table](#Supported-partition-table)
- [Supported partition types](#Supported-partition-types)
- [Supported file systems](#Supported-file-systems)
- [Drive, partition and size keys](#Drive,-partition-and-size-keys)
  - [Drive keys](#Drive-keys)
  - [Partition keys](#Partition-keys)
  - [Size keys](#Size-keys)

## Drive configuration example

A basic drive configuration looks like this:
```json
"drives": [
  {
    "path": "/dev/sda",
    "append": false,
    "partitions": [
      {
        "size": {
          "amount": 1,
          "unit": "GiB"
        },
        "partitionType": "C12A7328-F81F-11D2-BA4B-00A0C93EC93B"
      },
      {
        "size": {
          "amount": 4,
          "unit": "GiB"
        },
        "partitionType": "0657FD6D-A4AB-43C4-84E5-0933C84B4F4F"
      },
      {
        "size": {
          "takeRemaining": true
        },
        "partitionType": "4F68BCE3-E8CD-4DB1-96E7-FBCAF984B709",
        "fileSystem": "ext4"
      }
    ]
  }
],
```

## Supported partition table

Currently, the installer **only supports GPT drives**, so make sure to format drives on which you'll create partitions to GPT
before running the installer.

This can be done by doing:
```shell
$ fdisk /dev/sda
Command (m for help): g
Command (m for help): w
```

## Supported partition types

The installer currently support all the listed partition types below. The ones with an asterisk are **needed** for the system to be functional.

| Name | GUID | Description |
| ----- | ---- | ---- |
| EFI* | C12A7328-F81F-11D2-BA4B-00A0C93EC93B | Used for the booloader to boot the system. |
| SWAP* | 0657FD6D-A4AB-43C4-84E5-0933C84B4F4F | Systems swap space. |
| Root* | 4F68BCE3-E8CD-4DB1-96E7-FBCAF984B709 | The partition that will have root (/). |
| File System | 0FC63DAF-8483-4772-8E79-3D69D8477DE4 | Partition with data on it. |
| Home | 933AC7E1-2EB4-4F13-B844-0E14E2AEF915 | Optional partition for home directory. |

## Supported file systems

The installer currently support two file systems:
- btrfs
- ext4

## Drive, partition and size keys

This is a table of all the keys and descriptions for each of them.
 
### Drive keys
| Name | type | Description | Needed |
| --- | --- | --- | --- |
| path | string | The absolute path to the drive. | Yes |
| append | boolean | Whether the new partitions will be appended to the existing partition table or replace it. | Yes |
| partitions | array of [objects](#Partition-keys) | Array of all the partitions that need to be created. | Yes |

### Partition keys 
| Name | type | Description | Needed |
| --- | --- | --- | --- |
| size | [object](#Size-keys) | The size of the new partition. | Yes |
| partitionType | string | The GUID of the partition type. | Yes |
| fileSystem | string | The partition file system. | If partitionType is not EFI or SWAP. |
| mountPoint | string | The mount point of the drive | If partitionType is not EFI, SWAP or Root. |


### Size keys
| Name | type | Description | Needed |
| --- | --- | --- | --- |
| amount | integer | The size of the new drive. | If takeRemaining is **false** or **not present**. |
| unit | string | The unit of the given amount. Units are in *iB (like GiB).| If amount is specified. |
| takeRemaining | boolean | If true, it will take the remaining space of the drive for this partition. | If no amount is specified. |
