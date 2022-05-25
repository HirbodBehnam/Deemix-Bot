# Deemix Bot

A simple bot to upload downloaded musics from Deemix to Telegram.

## Features

* Full deemix features
* Builtin Deezer search
* Lightweight
* Partial ZSpotify support
* Custom downloader support

## Installing

At first make sure that you have `deemix` and `ffprobe` in your path. Then either clone and build this bot, or download
one of the executables from releases. Then place the `config.json` file next to the executable program and run it.

You can also optionally install `zspotify` in your path if you want it's support too.

### Building

Just use `go build` to build the bot.

### Config

The config file contains some fields:

1. `bot_token`: Which is simply your Telegram bot token which must be acquired from @BotFather
2. `users`: An array of int which is the user id of people who you want to allow working with this bot. You can get your
   ID using @myidbot. If you leave this array empty, the bot will be made public and anyone can access it.
3. `zspotify_credentials`: If you want ZSpotify support, run ZSpotify once to register your account then it creates a
   file called `credentials.json`. Add the content of this file to this key in config.
4. `custom_downloaders`: This is an array of objects which contain two keys: `name` and `url_prefix`. `name`
   is the application name which must be in path in order to run it. Also, bot will only send urls which are starting
   with `url_prefix` to the downloader and falls back to other downloaders if they don't match with it.
5. `direct_download_hosts`: You can also define some hosts in order to make the bot accept direct download links from
   them. These hosts can be defined as an array of strings in this key.

#### Making Spotify links work with Deemix

*Note:* This step is needed only if you wish to not use the ZSpotify.

If you want to make deemix accept spotify links, you have to create this file:
`~/.config/deemix/spotify/settings.json`. The contents of the file must look like this:

```json
{
  "clientId": "",
  "clientSecret": "",
  "fallbackSearch": false
}
```

#### ZSpotify Config

Configuring ZSpotify at the moment is kinda hard. So at first, you have to clone the ZSpotify repo. Then create a file
named `zspotify` and put these lines in it:

```bash
#!/bin/bash
python3.9 /opt/zspotify/zspotify $@
```

After, run `zspotify -s` and fill the account credentials for the first time. After that it should not ask for passwords
anymore.

<sub><sup>This help was written for ZSpotify commit `606125d`</sup></sub>

#### Adding a custom downloader

You can create your own downloaders to integrate it with this bot! To do so, create an executable and put it in your
path. The executable must accept an url as it's first argument. It should also save the music file in it's working
directory.

If an error occurs, it should return a non-zero status code.

## How it works?

This bot does not download anything from deezer or spotify itself; Instead, it simply executes `deemix` or `zspotify` in
temp folder to download the links which you give them to it, and uploads the downloaded files to Telegram. It also
uses `ffprobe` to get the length of musics which are larger than 10 MB to send them to Telegram.