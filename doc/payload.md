# JSON payload for main parser

```json
{
  "drives" : [
    {
      "path": "/dev/xyz",
      "append": true/false
      "partitions": [
        {
          "size": {
            "amount": 1234,
            "unit": "MiB/GiB/etc.",
            "takeRemaining": true/false
          },
          "fileSystem": "btrfs/ext4"
          "partitionType": "gpt partition type (guid)",
          "mountPoint": "/absolute/path/to/directory"
        }
      ],
    }
  ]
  "users": [
    {
      "username": "[username]",
      "password": "[password]",
      "homepath": "[path to home]",
      "sudoer": true,
    }
  ],
  "timezone": "[user timezone]",
  "locale": "[user locale]",
  "hostname": "[user hostname]",
  "rootpassword": "[root password]"
}
```
