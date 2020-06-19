import asyncio
import argparse
import logging
import pathlib

from thesaurize_loader import Loader


def get_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        prog="loader", description="Loads thesaurus data file into Redis datastore."
    )

    parser.add_argument(
        "--file",
        "-f",
        required=True,
        type=pathlib.Path,
        help="Path to .dat file containing thesaurus.",
    )
    parser.add_argument(
        "--connection", "-c", required=True, type=str, help="Redis URI to connect to."
    )
    parser.add_argument("--encoding", "-e", required=False, type=str, help="Data file encoding.")

    return parser.parse_args()


def setup_logging() -> None:
    log = logging.getLogger("loader-tool")
    log.setLevel(logging.INFO)

    handler = logging.StreamHandler()
    handler.setFormatter(logging.Formatter("%(asctime)s %(levelname)s:%(message)s"))

    log.addHandler(handler)


def main() -> None:
    setup_logging()

    args = get_args()
    loader = Loader(args)

    loop = asyncio.get_event_loop()
    loop.run_until_complete(loop.create_task(loader.run()))


if __name__ == "__main__":
    main()
