# Thesaurize 
A fun bot for Discord inspired by OrionSuperman's 
[ThesaurizeThis](https://github.com/orionsuperman/ThesaurizeThis) Reddit bot.
Given a sentence, it replaces each word with a synonym of that word.

## Usage
You will need a Discord bot account and its corresponding token in order to
proceed. There are good tutorials on how to do this elsewhere like
[this one](https://discordpy.readthedocs.io/en/latest/discord.html).

### Using Helm
This is the recommended method for deploying this application. Clone this 
repository and run the commands found below. Make sure to do them in the
following order.

1. First, create a persistent volume for Redis DB persistence. This depends
on your cluster environment so however you wish to store data is up to
you. See [the Kubernetes
docs](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) on
storage classes for more detail or the docs for your Kubernetes cluster
provider (if they offer storage). If you need something simple, the
[`hostPath`](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath)
class type is a good choice.

2. Install the chart.
```bash
helm install thesaurize chart/ \
    --create-namespace --namespace thesaurize-bot \
    --set "thesaurize.apiToken.value=$DISCORD_API_TOKEN" \
    --set "redis.datastore.storageClassName=$STORAGE_CLASS_NAME"
```

That's it. The bot should be up and running within a few seconds
once the Redis DB has finished loading.

## License
[MIT](https://choosealicense.com/licenses/mit/)
