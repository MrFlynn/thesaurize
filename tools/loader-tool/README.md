# Thesaurize Loader
This utility transforms [OpenOffice](https://openoffice.org) thesaurus data
files (based on Princeton's WordNet) into Redis protocol data streams. This
utility essentially mass-inserts thesaurus data into a Redis instance for use
with the [thesaurize bot](https://github.com/MrFlynn/thesaurize-bot) for
Discord.

You can read more and download the OpenOffice thesaurus
[here](https://www.openoffice.org/lingucomponent/thesaurus.html).

## Usage
Download the thesaurus archive linked above and extract it. You also need to
have a running instance of Redis on your system. Then use the following 
commands to insert data into a Redis instance.

```bash
$ pip install thesaurize-loader
$ thesaurize-loader --file=/path/to/thesaurus.dat --connection=redis://localhost:6379
```

## License
[MIT](../../LICENSE)
