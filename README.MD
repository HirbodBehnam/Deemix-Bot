# Deemix Bot
A simple bot to upload downloaded musics from Deemix to Telegram.

## Features
* Full deemix features
* Builtin Deezer search
* Very lightweight bot

## Installing
At first make sure that you have `deemix` and `ffprobe` in your path.
Then either clone and build this bot, or download one of the executables from releases.
Then place the `config.json` file next to the executable program and run it.

### Building
Just use `go build` to build the bot.

### Config
The config file contains two fields:
1. `bot_token`: Which is simply your Telegram bot token which must be acquired from @BotFather
2. `users`: An array of int which is the user id of people who you want to allow working with this bot. You can get your ID using @myidbot. If you leave this array empty, the bot will be made public and anyone can access it.

#### Making Spotify links work
If you want to make deemix accept spotify links, you have to create this file:
`~/.config/deemix/spotify/settings.json`. The contents of the file must look like this:
```json
{
  "clientId": "",
  "clientSecret": "",
  "fallbackSearch": false
}
```

## How it works?
This bot does not download anything from deezer or spotify itself; Instead, it simply executes `deemix` in temp folder to download the links which you give them to it, and uploads the downloaded files to Telegram.
It also uses `ffprobe` to get the length of musics which are larger than 10 MB to send them to Telegram.