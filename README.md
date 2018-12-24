# Thesaurize Bot
A bot for Discord that makes sentences make 70% less sense. Heavily inspired
by OrionSuperman's [Reddit bot](https://www.reddit.com/r/ThesaurizeThis/). 
By replacing each word in a sentence with a random synonym you can produce some
pretty hilarious results.

## Running the Bot
1. Get an API key from the [Big Huge Thesaurus](https://words.bighugelabs.com/api.php).
    Their free plan supports up to 1000 words per day.
2. Create a new application on Discord's [developer site](https://discordapp.com/developers/applications/).
    Create a bot user and copy the bot's token. 
3. Enable _Send Messages_ and _Read Message History_ permissions and copy the 
    resulting interger.
4. Next, copy the client ID for the bot.
5. Clone this repository and `cd` into it.
6. Run the following commands to build and then run the bot.
    ```bash
    $ make build
    $ docker run -d \
        --env DISCORD_KEY="DISCORD_KEY_VALUE" \
       --env THESAURUS_KEY="THESAURUS_KEY_VALUE" \
        thesaurize-bot:latest
    ```
7. Using the integer and client ID from steps 3 and 4, respectively, into the 
    following URL stub:
    `https://discordapp.com/api/oauth2/authorize?client_id=<ID>&scope=bot&permissions=<INTEGER>`
8. Click the link to add the bot to your server.