import argparse
import codecs
import logging
import pytest
import requests

from pathlib import Path
from thesaurize_loader import __version__, Loader, FileProtocol, HTTPProtocol, TX_BUFFER_TEMPLATE
from zipfile import ZipFile


# Default and utilities.
BASIC_ARGS = argparse.Namespace(encoding=None)
DATAFILE_CONTENT = """hello|2
(noun)|greetings
(verb)|hullo
"""


class RequestMock:
    def __init__(self, content: Path):
        with content.open("rb") as f:
            self.content = f.read()

        self.ok = True


@pytest.fixture(scope="session")
def create_data_file(tmpdir_factory):
    data_file = Path(tmpdir_factory.mktemp("data").join("data.dat"))
    zip_file = Path(tmpdir_factory.mktemp("archive").join("archive.zip"))

    with data_file.open("w") as f:
        f.write(DATAFILE_CONTENT)

    zipper = ZipFile(zip_file, "w")
    zipper.write(data_file)
    zipper.close()

    return data_file, zip_file


def test_version():
    assert __version__ == "0.2.0"


def test_default_encoding():
    l = Loader(BASIC_ARGS)

    assert l.encoding == "iso8859-1"


def test_set_encoding():
    l = Loader(BASIC_ARGS)
    l.encoding = "utf-8"

    assert l.encoding == "utf-8"


def test_invalid_encoding():
    l = Loader(BASIC_ARGS)

    with pytest.raises(LookupError):
        l.encoding = "invalid-encoding"


def test_encoding_on_create():
    l = Loader(argparse.Namespace(encoding="utf-8"))

    assert l.encoding == "utf-8"


def test_empty_buffer():
    l = Loader(BASIC_ARGS)
    full, buffer = l.buffer

    assert full == False
    assert buffer == TX_BUFFER_TEMPLATE


def test_buffer_push_single():
    l = Loader(BASIC_ARGS)
    l.push_to_buffer("hello", "noun", "world")

    full, buffer = l.buffer

    assert full == False
    assert buffer == {"noun": {"hello": ["world"]}, "verb": {}, "adj": {}, "adv": {}}


def test_buffer_push_multiple():
    l = Loader(BASIC_ARGS)
    l.push_to_buffer("this", "verb", "is", "a", "test")

    full, buffer = l.buffer

    assert full == False
    assert buffer == {
        "noun": {},
        "verb": {"this": ["is", "a", "test"]},
        "adj": {},
        "adv": {},
    }


def test_buffer_push_to_existing_key():
    l = Loader(BASIC_ARGS)
    l.push_to_buffer("this", "noun", "that")
    l.push_to_buffer("this", "noun", "here")

    full, buffer = l.buffer

    assert full == False
    assert buffer == {
        "noun": {"this": ["that", "here"]},
        "verb": {},
        "adj": {},
        "adv": {},
    }


def test_metadata_reader(create_data_file):
    l = Loader(BASIC_ARGS)

    reader = iter(FileProtocol(create_data_file[0], l._encoding, None))
    l.read_word_metadata(reader)

    full, buffer = l.buffer

    assert full == False
    assert buffer == {
        "noun": {"hello": ["greetings"]},
        "verb": {"hello": ["hullo"]},
        "adj": {},
        "adv": {},
    }


def test_file_protocol_contains():
    assert FileProtocol.supports_protocol("file://")


def test_file_protocol_size(create_data_file):
    proto = FileProtocol(create_data_file[0], None, None)

    assert proto.size == 38


def test_file_protocol(create_data_file):
    proto = FileProtocol(create_data_file[0], codecs.lookup("ISO8859-1"), None)

    iterator = iter(proto)
    assert next(iterator) == (8, "hello|2")
    assert next(iterator) == (17, "(noun)|greetings")
    assert next(iterator) == (13, "(verb)|hullo")

    with pytest.raises(StopIteration):
        next(iterator)


def test_http_protocol_contains():
    assert HTTPProtocol.supports_protocol("http://")
    assert HTTPProtocol.supports_protocol("https://")


def test_http_protocol_size(create_data_file, monkeypatch):
    log = logging.getLogger("test")

    def mockrequests(*args, **kwargs):
        return RequestMock(create_data_file[1])

    monkeypatch.setattr(requests, "get", mockrequests)

    proto = HTTPProtocol("", codecs.lookup("ISO8859-1"), log)
    assert proto.size == 38


def test_http_protocol(create_data_file, monkeypatch):
    log = logging.getLogger("test")

    def mockrequests(*args, **kwargs):
        return RequestMock(create_data_file[1])

    monkeypatch.setattr(requests, "get", mockrequests)

    proto = HTTPProtocol("", codecs.lookup("ISO8859-1"), log)

    iterator = iter(proto)
    assert next(iterator) == (8, "hello|2")
    assert next(iterator) == (17, "(noun)|greetings")
    assert next(iterator) == (13, "(verb)|hullo")

    with pytest.raises(StopIteration):
        next(iterator)
