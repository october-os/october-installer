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

> [!CAUTION]
> **Running this as super user will wipe your computer.**
>
> Do not run the installer directly on your main machine for development. Use a virtual machine.

To use October installer, you can either give it a JSON file or string.

To install from a file, use the `-json` flag:
```shell
# october-installer -json [file path]
```

To install from a string, use the `-json-raw` flag:
```shell
# october-installer -json-raw [json string]
```

## Development

For development, the installer is using two modes that completes each other. This allows safe, fast iterations.

### TDD
For normal development, we use test-driven development to allow fast testing and regression testing. 

Running an installer at each new changes dramatically slows down testing, so we recommend creating and using tests during
the core development of a feature or a bug fix. 

This command will run all tests:
```bash
go test ./...
```

### Virtual machine
When you're done with developping, you **must** at least have one working run in a virtual machine.

The directory `test-json` contains JSON configurations that are known to work and we encourage you to use them
for testing.

The setup we currently have is for QEMU. You need a valid QEMU virtual machine with two things:
- Shared memory enabled
- A `virtiofs` shared file system between the host machine and the VM.

The shared file system should look like this in XML:
```xml
<filesystem type="mount" accessmode="passthrough">
  <driver type="virtiofs"/>
  <source dir="[host machine installer dir]"/>
  <target dir="installer"/>
  <address type="pci" domain="0x0000" bus="0x01" slot="0x00" function="0x0"/>
</filesystem>
```

To generate a new October Linux ISO for testing, you can check the instructions in [the ISO repository](https://github.com/october-os/october-iso).


## Configuration documentation

All the documentation for the JSON configuration is in the `doc` directory or on [the October Linux website](https://octoberlinux.org/docs/).
