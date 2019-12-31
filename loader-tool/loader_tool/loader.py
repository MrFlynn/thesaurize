import aioredis
import argparse
import codecs
import logging
import os
import typing

from copy import deepcopy

DEFAULT_ENCODING = codecs.lookup("ISO8859-1")

BUFF_T = typing.TypeVar("BUFF_T", bound=typing.Dict[str, typing.Dict[str, typing.Any]])
TX_BUFFER_MAX_SIZE = 100
TX_BUFFER_TEMPLATE = {"noun": {}, "verb": {}, "adj": {}, "adv": {}}

log = logging.getLogger(__name__)


class Loader:
    def __init__(self, args: argparse.Namespace):
        self.args = args

        self._encoding: codecs.CodecInfo = DEFAULT_ENCODING
        if args.encoding:
            self.encoding = args.encoding

        self._redis: aioredis.RedisConnection = None

        self._tx_buffer: BUFF_T = deepcopy(TX_BUFFER_TEMPLATE)
        self._tx_buffer_size = 0

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

    async def push_redis(self) -> None:
        if not self._redis:
            log.error("Redis connection not setup.")
            return

        if self._tx_buffer_size < TX_BUFFER_MAX_SIZE:
            return

        pipeline = self._redis.pipeline()
        for lexograph, section in self._tx_buffer.items():
            for word, synonyms in section.items():
                pipeline.sadd(f"{lexograph}:{word}", *synonyms)

        await self._redis.execute()
        self._tx_buffer = deepcopy(TX_BUFFER_TEMPLATE)

    def read(self) -> typing.Iterator[str]:
        if not self.args.file.exists() or not self.args.file.is_file():
            log.error(f"File {self.args.file.as_posix} does not exists or is not a file")
            return

        with self.args.file.open("rb") as f:
            # Get the size of the file.
            file_size = os.fstat(f.fileno()).st_size

            # Create buffered reader.
            reader = self._encoding.streamreader(f)

            # Keep reading lines until current position exceeds file size.
            while line := reader.readline(keepends=False):
                if f.tell() > file_size:
                    break

                yield line

    def read_word_metadata(self, reader: typing.Callable[[], typing.Iterator[str]]) -> None:
        gen = reader()
        word, num_sections = next(gen).split("|")

        for _ in range(int(num_sections)):
            items = next(gen).split("|")
            section = items[0][1:-1]  # Remove parentheses.

            self.push_to_buffer(word, section, *items[1:])

    async def run(self) -> None:
        self._redis = await aioredis.create_connection(self.args.connection)
        log.info("Database connected.")

        # Check if first line of encoding matches first line of file.
        reader = self.read()
        if (encoding := next(reader)) != self.encoding.lower():
            log.error(f"Encoding mismatch! Expected {self.encoding}, got {encoding}.")
            return

        while True:
            try:
                self.read_word_metadata(reader)
                await self.push_redis()
            except StopIteration:
                await self.push_redis()

                log.info("File imported to database.")
