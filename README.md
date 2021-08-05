# PoDownloader

[![GitHub license](https://img.shields.io/github/license/LGiki/PoDownloader?style=flat-square)](https://github.com/LGiki/PoDownloader) [![GitHub release (latest by date)](https://img.shields.io/github/v/release/LGiki/PoDownloader?style=flat-square)](https://github.com/LGiki/PoDownloader/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/LGiki/PoDownloader)](https://goreportcard.com/report/github.com/LGiki/PoDownloader)

🎙️⬇️ PoDownloader = **Po**dcast **Downloader**, a simple CLI tool to download podcasts.

This tool will download podcast RSS, podcast cover image, episode audio files, episode cover images and episode shownotes.

[中文说明](https://github.com/LGiki/PoDownloader/blob/master/README.zh_CN.md)

# Screenshot

![](https://raw.githubusercontent.com/LGiki/PoDownloader/master/screenshot/screenshot.png)

# Install

## Download from Releases

You can download latest release from [Releases page](https://github.com/LGiki/PoDownloader/releases) directly.

## Build from source code

Make sure [go](https://golang.org/) is installed on your system correctly.

```bash
git clone https://github.com/LGiki/PoDownloader.git
cd PoDownloader
go mod download
go build -o podownloader ./cmd
```

Then you can find the output binary file named `podownloader`.

# Usage

## Download podcasts from OPML

Most podcast apps support exporting podcast lists in [OPML format](https://en.wikipedia.org/wiki/OPML).

An example of an OPML file is as follows:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
    <head>
        <title>Example OPML</title>
    </head>
    <body>
        <outline text="Example Podcast 1" title="Example Podcast 1" type="rss" xmlUrl="https://exmaple.org/podcast1/rss.xml" />
        <outline text="Example Podcast 2" title="Example Podcast 2" type="rss" xmlUrl="https://exmaple.org/podcast2/rss.xml" />
        <outline text="Example Podcast 3" title="Example Podcast 3" type="rss" xmlUrl="https://exmaple.org/podcast3/rss.xml" />
    </body>
</opml>
```

This OPML file contains 3 podcasts (each `outline` tag is a podcast), and the `xmlUrl` attribute is the podcast RSS link.

Download podcasts from OPML file using:

```bash
podownloader download --opml /path/to/opml_file.xml
```

## Download podcasts from RSS links list

RSS links list file is a text file, one podcast RSS URL per line, for example:

```
https://exmaple.org/podcast1/rss.xml
https://example.org/podcast2/rss.xml
https://example.org/podcast3/rss.xml
```

Download podcasts from RSS links list file using:

```bash
podownloader download --list /path/to/rss_list_file.txt
```

## Download podcast from RSS link

```
podownloader download --rss https://example.org/podcast/rss.xml
```

# Download Options

Using `-h` or `--help` to view all options.

Use the `HTTP_PROXY` environment variable to set a HTTP or SOCSK5 proxy.

## Output directory

Using `-o` or `--output` to specify output directory, default output directory is `./podcast`.

## User agent

Using `-u` or `--ua` to specify user agent, default user agent is `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36`.

## Download threads

Using `-t` or `--thread` to specify download threads, default download threads is `3`.

## Log directory

You can specify `--log` parameter to set the log directory.

If you specify the `--log` parameter, the log file will be saved to the directory you specified.

If you leave the `--log` parameter empty, no log file will be generated.

Default value of `--log` is empty.

# Configuration file

If you don't want to specify parameters every time you run the program, you can save the parameters in a configuration file, the program will automatically load the parameters from the configuration file.

You can specify configuration file path using `-c` or `--config`:

```bash
podownloader --config ~/.podownloader.json
```

Default configuration file path is `$PWD/.podownloader`.

Supported configuration file formats: `json`, `toml`, `yaml`, `yml`, `properties`, `props`, `prop`, `hcl`, `dotenv`, `env`, `ini`.

An example of a configuration file in JSON format is as follows:

```json
{
    "list": "",
    "opml": "/path/to/opml_file.xml",
    "rss": "",
    "output": "podcast",
    "ua": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36",
    "thread": 3
}
```

You can find more configuration file templates in [config_template](https://github.com/LGiki/PoDownloader/tree/master/config_template) folder.

# Download folder structure

```
podcast
├─ podcast_1_title
│  ├─ episode_1_title
│  │  ├─ cover.jpg
│  │  ├─ episode_1_title.mp3
│  │  └─ shownotes.html
│  ├─ episode_2_title
│  │  ├─ cover.jpg
│  │  ├─ episode_2_title.mp3
│  │  └─ shownotes.html
│  ├─ ...
│  │  ├─ cover.jpg
│  │  ├─ *****.mp3
│  │  └─ shownotes.html
│  ├─ ...
│  ├─ cover.jpg
│  └─ rss.xml
└─ podcast_2_title
   ├─ episode_1_title
   │  ├─ cover.jpg
   │  ├─ episode_1_title.mp3
   │  └─ shownotes.html
   ├─ episode_2_title
   │  ├─ cover.jpg
   │  ├─ episode_2_title.mp3
   │  └─ shownotes.html
   ├─ ...
   │  ├─ cover.jpg
   │  ├─ *****.mp3
   │  └─ shownotes.html
   ├─ ...
   ├─ cover.jpg
   └─ rss.xml
```

# License

[Apache-2.0](https://github.com/LGiki/PoDownloader/blob/master/LICENSE)
