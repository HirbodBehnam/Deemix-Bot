# Deemix Bot

A simple bot to upload downloaded musics from Deemix to Telegram.

## Features

* Full deemix features
* Builtin Deezer search
* Very lightweight bot
* Partial ZSpotify support

## Installing

At first make sure that you have `deemix` and `ffprobe` in your path. Then either clone and build this bot, or download
one of the executables from releases. Then place the `config.json` file next to the executable program and run it.

You can also optionally install `zspotify` in your path if you want it's support too.

### Building

Just use `go build` to build the bot.

### Config

The config file contains two fields:

1. `bot_token`: Which is simply your Telegram bot token which must be acquired from @BotFather
2. `users`: An array of int which is the user id of people who you want to allow working with this bot. You can get your
   ID using @myidbot. If you leave this array empty, the bot will be made public and anyone can access it.
3. `zspotify_credentials`: If you want ZSpotify support, run ZSpotify once to register your account then it creates a
   file called `credentials.json`. Add the content of this file to this key in config.

#### Making Spotify links work with Deemix

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

Configuring ZSpotify at the moment is pretty fucking hard. So at first, you have to clone the ZSpotify repo. Then create
a file named `zspotify` and put these lines in it:

```bash
#!/bin/bash
python3.9 /opt/zspotify/zspotify "$@"
```

You can change the ZSpotify path or python version. Then make the script executable with `chmod +x zspotify` and place
it somewhere in your path. Next go to the cloned version of ZSpotify. There are two files which you need to
edit; `const.py` and `zspotify.py`. At first, open the `const.py` and change
the `CONFIG_FILE_PATH = '../zs_config.json'` line to `CONFIG_FILE_PATH = 'zs_config.json'`. After open `zspotify.py` and
look for the `load_config` function. Change the `true_config_file_path = os.path.join(app_dir, CONFIG_FILE_PATH)` file
to `true_config_file_path = CONFIG_FILE_PATH`. Save and now you can use this bot!

<sub><sup>This help was written for ZSpotify commit `76f4089`</sup></sub>

## How it works?

This bot does not download anything from deezer or spotify itself; Instead, it simply executes `deemix` or `zspotify` in
temp folder to download the links which you give them to it, and uploads the downloaded files to Telegram. It also
uses `ffprobe` to get the length of musics which are larger than 10 MB to send them to Telegram.