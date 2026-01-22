# User

Users contains an array of users that needs to be created after the installation.

- [Users configuration example](#Users-configuration-example)
- [User keys](#User-keys)

## Users configuration example
```json
{
  "username": "testuser",
  "password": "test",
  "homepath": "/home/testuser",
  "sudoer": true
}
```

## User keys

| Name | Type | Description | Needed |
| --- | --- | --- | --- |
| username | string | The username of the user. | Yes |
| password | string | The password of the user. | Yes |
| homepath | string | The absolute path to the user home folder. | No. Default: /home/[username] |
| sudoer | boolean | Is the user a sudoer on the new install. | Yes |
