# Thesaurize Loader
This utility transforms [OpenOffice](https://openoffice.org) thesaurus data
files (based on Princeton's WordNet) into Redis protocol data streams. This
utility essentially mass-inserts thesaurus data into a Redis instance for use
with the [thesaurize bot](https://github.com/MrFlynn/thesaurize) for
Discord.

You can read more and download the OpenOffice thesaurus
[here](https://www.openoffice.org/lingucomponent/thesaurus.html).

## Usage
You will need to install this utility with pip(x) and have Redis installed and
running. Then run the utility with the following arguments.

```bash
$ pipx install thesaurize-loader
$ thesaurize-loader \
    --file=https://www.openoffice.org/lingucomponent/MyThes-1.zip \
    --connection=redis://localhost:6379
```

Alternatively, you can download the thesaurus archive linked above and extract 
it. Then run the following command:

```bash
$ thesaurize-loader --file=file:///path/to/thesaurus.dat --connection=redis://localhost:6379
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
