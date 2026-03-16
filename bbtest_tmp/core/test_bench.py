from core import Test
from core.assertions import AssertionError

RED = "\033[31m"
GREEN = "\033[32m"
RESET = "\033[0m"


class TestBench:
    def __init__(self, name: str) -> None:
        self._name = name
        self._tests: list[Test] = []
        self._failed_tests: dict[Test, AssertionError] = {}

    def print_failed_tests(self) -> None:
        for test, error in self._failed_tests.items():
            print(f"Test failed: {test.name} with error: {error}")

    def register(self, test: Test) -> None:
        self._tests.append(test)

    def run(self) -> bool:
        print(f"\nStarting test bench: {self._name}")

        successful_tests = 0
        for test in self._tests:
            try:
                test.run()
                successful_tests += 1
            except AssertionError as exception:
                self._failed_tests[test] = exception
            except Exception as e:
                print(f"{e}")
        if successful_tests == len(self._tests):
            print(
                f"{GREEN}O{RESET} {self._name} succeeded: {successful_tests}/{len(self._tests)}"
            )
            return True
        print(
            f"{RED}X{RESET} {self._name} failed: {successful_tests}/{len(self._tests)}"
        )
        return False
