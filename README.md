# Thesaurize 
A fun bot for Discord inspired by OrionSuperman's 
[ThesaurizeThis](https://github.com/orionsuperman/ThesaurizeThis) Reddit bot.
Given a sentence, it replaces each word with a synonym of that word.

## Usage
You will need a Discord bot account and its corresponding token in order to
proceed. There are good tutorials on how to do this elsewhere like
[this one](https://discordpy.readthedocs.io/en/latest/discord.html).

### Using Kubernetes
This is the recommended method for deploying this application. Clone this 
repository and run the commands found below. Make sure to do them in the
following order.

1. First, add the namespace.

```bash
$ kubectl apply -f namespace.yml
```

2. Create a secret with your bot token. Make sure your secret is in 
your clipboard.
```bash
$ cd deployments/
$ pbpaste > discord-token.txt
$ kubectl -n=thesaurize-bot create secret discord-token --from-file=./discord-token.txt && \
    rm -f discord-token.txt; pbcopy ""
```

3. Next, add the volume.

```bash
$ kubectl apply -f volume.yml
```

4. Create the Redis service and load data into it.

```bash
$ kubectl apply -f redis.yaml
$ kubectl apply -f loader-job.yaml -f loader-scripts.yaml
```

5. Finally, run the bot.

```bash
$ kubectl apply -f thesaurize.yaml
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
