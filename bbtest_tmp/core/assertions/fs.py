import os


def assertDirectoryExists(path: str) -> None:
    if not os.path.isdir(path):
        raise AssertionError(
            f"Failed to find directory.\n\tExpected: '{path}' to exist"
        )


def assertFileExists(path: str) -> None:
    if not os.path.isfile(path):
        raise AssertionError(f"Failed to find file.\n\tExpected: '{path}' to exist")


def assertFileContains(path: str, expected: str) -> None:
    with open(path, "r") as file:
        actual = file.read()
        if actual != expected:
            raise AssertionError(
                f"File content does not match.\n\tActual: '{actual}'\n\tExpected: '{expected}'"
            )
