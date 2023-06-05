import subprocess
from pathlib import Path


def protogen(setup_kwargs):
    subprocess.run(["./proto-gen.sh"], cwd=Path(__file__).parent, shell=True, check=True)
    return setup_kwargs


if __name__ == "__main__":
    protogen({})
