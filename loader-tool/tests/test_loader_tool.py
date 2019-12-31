import argparse
import pytest

from loader_tool import __version__, Loader, TX_BUFFER_TEMPLATE
from pathlib import Path


# Default and utilities.
BASIC_ARGS = argparse.Namespace(encoding=None)
DATAFILE_CONTENT = """hello|2
(noun)|greetings
(verb)|hullo
"""


@pytest.fixture(scope="session")
def create_data_file(tmpdir_factory):
    filename = Path(tmpdir_factory.mktemp("data").join("data.dat"))

    with filename.open("w") as f:
        f.write(DATAFILE_CONTENT)

    return filename


def test_version():
    assert __version__ == "0.1.0"


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


def test_reader(create_data_file):
    args = argparse.Namespace(encoding=None, file=create_data_file)
    l = Loader(args)

    assert list(l.read()) == ["hello|2", "(noun)|greetings", "(verb)|hullo"]


def test_metadata_reader(create_data_file):
    args = argparse.Namespace(encoding=None, file=create_data_file)
    l = Loader(args)

    l.read_word_metadata(l.read)

    full, buffer = l.buffer

    assert full == False
    assert buffer == {
        "noun": {"hello": ["greetings"]},
        "verb": {"hello": ["hullo"]},
        "adj": {},
        "adv": {},
    }
