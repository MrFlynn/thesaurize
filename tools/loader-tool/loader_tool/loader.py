import aioredis
import argparse
import codecs
import logging
import os
import typing

from copy import deepcopy
from progress.bar import Bar

DEFAULT_ENCODING = codecs.lookup("ISO8859-1")

BUFF_T = typing.TypeVar("BUFF_T", bound=typing.Dict[str, typing.Dict[str, typing.Any]])
TX_BUFFER_MAX_SIZE = 10000
TX_BUFFER_TEMPLATE = {"noun": {}, "verb": {}, "adj": {}, "adv": {}}

log = logging.getLogger("loader-tool")


class Loader:
    def __init__(self, args: argparse.Namespace):
        self.args = args

        self._encoding: codecs.CodecInfo = DEFAULT_ENCODING
        if args.encoding:
            self.encoding = args.encoding

        self._redis: aioredis.Redis = None

        self._tx_buffer: BUFF_T = deepcopy(TX_BUFFER_TEMPLATE)
        self._tx_buffer_size = 0
        self._terminate = False

        self._progress = Bar("Loading Data", suffix="%(percent).1f%% - %(eta)ds", max=100)

    @property
    def encoding(self) -> str:
        return self._encoding.name

    @encoding.setter
    def encoding(self, encoding_name: str) -> None:
        self._encoding = codecs.lookup(encoding_name)

    @property
    def buffer(self) -> typing.Tuple[bool, BUFF_T]:
        return (self._tx_buffer_size == TX_BUFFER_MAX_SIZE), self._tx_buffer

    def push_to_buffer(self, word: str, section: str, *args) -> None:
        if word not in self._tx_buffer[section]:
            self._tx_buffer[section][word] = list(args)

            self._tx_buffer_size += 1
        else:
            self._tx_buffer[section][word].extend(list(args))

    async def push_redis(self, push_remaining: bool = False) -> None:
        if not self._redis:
            log.error("Redis connection not setup.")
            return

        full, _ = self.buffer
        if full or push_remaining:
            pipeline = self._redis.pipeline()

            for lexograph, section in self._tx_buffer.items():
                for word, synonyms in section.items():
                    pipeline.sadd(f"{lexograph}:{word}", *synonyms)

            await pipeline.execute()
            self._tx_buffer = deepcopy(TX_BUFFER_TEMPLATE)

    def read(self) -> typing.Iterator[str]:
        if not self.args.file.exists() or not self.args.file.is_file():
            log.error(f"File {self.args.file.as_posix} does not exists or is not a file")
            return

        with self.args.file.open("rb") as f:
            # Get the size of the file.
            file_size = os.fstat(f.fileno()).st_size
            self._progress.max = file_size

            # Create buffered reader.
            reader = self._encoding.streamreader(f)

            size = 0
            while line := reader.readline(keepends=False):
                # Keep reading lines until current position exceeds file size.
                if size >= file_size:
                    break

                self._progress.next(n=(f.tell() - self._progress.index))

                size += len(line.encode()) + 1  # Account for line endings.
                yield line

    def read_word_metadata(self, reader: typing.Iterator[str]) -> None:
        try:
            word, num_sections = next(reader).split("|")

            for _ in range(int(num_sections)):
                items = next(reader).split("|")
                section = items[0][1:-1]  # Remove parentheses.

                self.push_to_buffer(word, section, *items[1:])
        except StopIteration:
            self._terminate = True

    async def run(self) -> None:
        self._redis = await aioredis.create_redis(self.args.connection)
        log.debug("Database connected.")

        # Check if first line of encoding matches first line of file.
        reader = self.read()
        if (encoding := next(reader).lower()) != self.encoding:
            log.error(f"Encoding mismatch! Expected {self.encoding}, got {encoding}.")
            return

        while not self._terminate:
            self.read_word_metadata(reader)
            await self.push_redis(push_remaining=self._terminate)

        log.debug("Database finished loading.")

        # Send command to `status` pubsub channel that data has been loaded.
        await self._redis.publish("status", "ready")
        log.debug("Pubsub channel `status` has been updated with `ready` message.")

        self._redis.close()

        self._progress.finish()
