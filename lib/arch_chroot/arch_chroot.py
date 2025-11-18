import subprocess

# Mount point of "/"
_root_mount_point = "/mnt"


class ArchChrootExecutionError(Exception):
    """
    A class that represent an error when arch-chroot was executed.

    ...

    Attributes
    ----------
    return_code: int
        Arch-chroot return code
    error_message: str
        The string returned on STDERR by arch-chroot

    Methods
    -------
    to_string(): str
        Returns a string that contains the error code and the error message
    """

    def __init__(self, return_code: int, error_message: str):
        self.return_code = return_code
        self.error_message = error_message

    def to_string(self) -> str:
        return f"{self.return_code}: {self.error_message}"


def run(command: str) -> None:
    """Will execute the given command in arch-chroot.

    If command is null or empty, the function will just return without executing the command.

    Parameters
    ----------
    command: str
        The command that will be ran using arch-chroot

    Raises
    ------
    ArchChrootExecutionError
        If arch-chroot fails to execute or returns a none 0 code
    """

    if not command:
        return

    try:
        subprocess.run(
            ["arch-chroot", _root_mount_point, command],
            check=True,
            text=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
    except subprocess.CalledProcessError as e:
        raise ArchChrootExecutionError(e.returncode, e.stderr)
