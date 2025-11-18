import subprocess

_root_mount_point = "/mnt"


class ArchChrootExecutionError(Exception):
    def __init__(self, return_code, error_message):
        self.return_code = return_code
        self.error_message = error_message


def run(command: str) -> None:
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
