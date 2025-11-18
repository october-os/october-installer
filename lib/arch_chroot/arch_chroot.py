import subprocess


def run(command: str):
    mount_point = "/mnt"

    subprocess.run(["arch-chroot", mount_point, command])
