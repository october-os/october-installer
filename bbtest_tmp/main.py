import subprocess

from test_suites.config import bench as config_bench

result = subprocess.run(
    ["./installer/main", "-json", "installer/test-json/json_struct_test.json"]
)

if result.returncode != 0:
    print(f"Installer failed with return code: {result.returncode}")
    exit(1)

print("\n------Starting tests------")

config_bench.run()
config_bench.print_failed_tests()
