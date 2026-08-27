# The official October Linux installer

October installer is a fully working Arch Linux installer made in Go. It can install a working October Linux installation
from a JSON file or string containing the configuration.

## Quickstart

Compiling the installer only requires the [Go compiler](https://go.dev/). When it is installed, you can run:
```bash
git clone https://github.com/october-os/october-installer.git
cd october-installer
go build cmd/main.go -o october-installer
```

This produces a working installer called `october-installer` in the current directory.

To use October installer, you can either give it a JSON file or string.
> [!CAUTION]
> **Running this as super user will wipe your computer.**
>
> Do not run the installer directly on your main machine for development. Use a virtual machine.

To install from a file, use the `-json` flag:
```shell
# october-installer -json [jsonfile]
```

To install from a string, use the `-json-raw` flag:
```shell
# october-installer -json-raw [json string]
```

## Development



## Documentation
