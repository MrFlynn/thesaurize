import aioredis
import argparse
import codecs
import logging
import os
import re
import typing

from .protocols import Protocol, ProtocolFactory, FileProtocol, HTTPProtocol

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

    def read_word_metadata(self, reader: typing.Iterator[typing.Tuple[int, str]]) -> None:
        try:
            inc, header_line = next(reader)
            self._progress.next(n=inc)

            word, num_sections = header_line.split("|")
            for _ in range(int(num_sections)):
                inc, synonym_line = next(reader)
                self._progress.next(n=inc)

                items = synonym_line.split("|")
                section = items[0][1:-1]  # Remove parentheses.

                self.push_to_buffer(word, section, *items[1:])
        except StopIteration:
            self._terminate = True

    async def run(self) -> None:
        exit_code = 0

        try:
            self._redis = await aioredis.create_redis(self.args.connection)
            log.debug("Database connected.")

            # Search for valid protocol handler.
            proto_handler: Protocol
            matcher = re.compile("^[A-z]+://")

            if (proto_string := matcher.search(self.args.file)) != None:
                proto_handler = ProtocolFactory.create(proto_string.group(0))

                if isinstance(proto_handler, FileProtocol):
                    proto_handler = proto_handler(
                        matcher.sub("", self.args.file), self._encoding, log
                    )
                else:
                    proto_handler = proto_handler(self.args.file, self._encoding, log)
            else:
                raise RuntimeError("No protocol specified. You must specify one.")

            self._progress.max = proto_handler.size

            # Initializer reader and validate encoding.
            reader = iter(proto_handler)

            inc, encoding_header = next(reader)
            if (encoding := encoding_header.lower()) != self.encoding:
                RuntimeError(f"Encoding mismatch! Expected {self.encoding}, got {encoding}.")

            self._progress.next(n=inc)

            while not self._terminate:
                self.read_word_metadata(reader)
                await self.push_redis(push_remaining=self._terminate)

            log.debug("Database finished loading.")

            # Send command to `status` pubsub channel that data has been loaded.
            await self._redis.publish("status", "ready")
            log.debug("Pubsub channel `status` has been updated with `ready` message.")
        except Exception as e:
            log.error(e)
            exit_code = 1
        finally:
            self._redis.close()
            self._progress.finish()

            exit(exit_code)
