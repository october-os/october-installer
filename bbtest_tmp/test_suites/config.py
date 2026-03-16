from core import Test, TestBench
from core.assertions import (
    assertDirectoryExists,
    assertFileExists,
)

bench = TestBench("October Configuration")


class DoesConfigurationFolderExist(Test):
    def __init__(self) -> None:
        super().__init__("Does configuration folder exist?")

    def run(self) -> None:
        assertDirectoryExists("/mnt/home/testuser/.config/october-config")
        assertDirectoryExists("/mnt/home/secondtestuser/.config/october-config")


class Penistest(Test):
    def __init__(self) -> None:
        super().__init__("Does configuration folder exist?")

    def run(self) -> None:
        assertDirectoryExists("/mnt/home/testuser/.config/october-config/penis")
        assertDirectoryExists("/mnt/home/secondtestuser/.config/october-config")


bench.register(DoesConfigurationFolderExist())
bench.register(Penistest())
