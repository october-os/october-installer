## General configuration

This part contains all the more "general" configuration of the installation.

- [General configuration example](#General-configuration-example)
- [General configuration keys](#General-configuration-keys)

## General configuration example
```json
{
  "mirrorCountries": ["Canada"],
  "timezone": "America/Montreal",
  "locale": "US.UTF-8",
  "hostname": "testhostname",
  "rootPassword": "test",
  "bestEffortGPU": false,
  "extraPackages": {
    "officialRepositories": ["package1", "package2"],
    "aur": ["package1", "package2"]
  }
}
```

## General configuration keys

| Name | Type | Description | Needed |
| --- | --- | --- | --- |
| mirrorCountries | array of strings | Names of the countries you want to use mirrors from. They can be seen on the [Arch Wiki](https://archlinux.org/mirrorlist/all/https/).| Yes |
| timezone | string | The timezone you want the system to be set to. | Yes |
| locale | string | The locale you want to set on the system. Only UTF-8 locales are supported. | Yes |
| hostname | string | The hostname of the new system. | Yes |
| rootPassword | string | The root password on the new install. | Yes |
| bestEffortGPU | boolean | If true, it will attempt to install the drivers for the systems GPU on the new install. Only Nvidia, AMD and Intel are supported. | Yes |
| extraPackages | object | The extra packages to be installed on the system. | No |
| extraPackages.officialRepositories | array of strings | The extra packages to be installed from the official Arch Linux repositories | No |
| extraPackages.aur | array of strings | The extra packages to be installed from the Arch User Repository (AUR) | No |
