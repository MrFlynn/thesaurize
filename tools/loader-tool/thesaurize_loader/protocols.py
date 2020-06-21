import codecs
import io
import logging
import os
import re
import requests
import warnings
import zipfile

from abc import ABCMeta, abstractmethod
from typing import Iterator, List, Tuple

warnings.filterwarnings("ignore")


class Protocol(metaclass=ABCMeta):
    PROTOCOLS = set()

    def __init__(self, location: str, codec: codecs.CodecInfo, logger: logging.Logger,) -> None:
        self._logger = logger

        self._location = location
        self._codec = codec
        self._reader = None
        self._size = 0
        self._progress = 0

    @classmethod
    def supports_protocol(cls, proto: str) -> bool:
        return proto in cls.PROTOCOLS

    @property
    @abstractmethod
    def size(self) -> int:
        raise NotImplementedError

    @abstractmethod
    def __iter__(self) -> "Protocol":
        raise NotImplementedError

    def __next__(self) -> Iterator[Tuple[int, str]]:
        line = self._reader.readline(keepends=False)
        if (inc := len(line.encode()) + 1) + self._progress <= self._size:
            self._progress += inc
            return inc, line
        else:
            raise StopIteration


class FileProtocol(Protocol):
    PROTOCOLS = set(["file://"])

    def __init__(self, location: str, codec: codecs.CodecInfo, logger: logging.Logger) -> None:
        super(FileProtocol, self).__init__(location, codec, logger)
        self._file = None

    def _load_file(self) -> None:
        if not self._file:
            self._file = open(self._location, "rb")
            self._size = os.fstat(self._file.fileno()).st_size

    @property
    def size(self) -> int:
        self._load_file()
        return self._size

    def __iter__(self) -> "FileProtocol":
        self._load_file()
        self._reader = self._codec.streamreader(self._file)

        return self

    def __del__(self) -> None:
        self._file.close()


class HTTPProtocol(Protocol):
    PROTOCOLS = set(["http://", "https://"])

    def __init__(self, location: str, codec: codecs.CodecInfo, logger: logging.Logger) -> None:
        super(HTTPProtocol, self).__init__(location, codec, logger)

    def _find_data_file(self, content: io.BytesIO) -> None:
        remote_zip = zipfile.ZipFile(content)

        filename = None
        for f in remote_zip.namelist():
            if (m := re.search(".+\.dat", f)) != None:
                filename = m.group(0)

        if not filename:
            raise FileNotFoundError(f"Could not find .dat file in {self._location}")

        thesaurus = remote_zip.read(filename)
        self._size = len(thesaurus)
        self._reader = self._codec.streamreader(io.BytesIO(thesaurus))

    def _get(self) -> None:
        if self._reader:
            return

        try:
            r = requests.get(self._location, verify=False, stream=True)

            if r.ok:
                self._logger.debug(f"Got {self._location}")
                self._find_data_file(io.BytesIO(r.content))
        except requests.exceptions.RequestException:
            self._logger.error(f"Could not get resource {self._location}")

    @property
    def size(self) -> int:
        self._get()

        return self._size

    def __iter__(self) -> "HTTPProtocol":
        self._get()

        return self


class ProtocolFactory:
    @classmethod
    def create(cls, proto: str) -> Protocol:
        for c in Protocol.__subclasses__():
            if c.supports_protocol(proto):
                return c

        raise RuntimeError(f"Protocol {proto} not supported.")
