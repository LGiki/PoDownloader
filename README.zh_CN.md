# PoDownloader

[![GitHub license](https://img.shields.io/github/license/LGiki/PoDownloader?style=flat-square)](https://github.com/LGiki/PoDownloader) [![GitHub release (latest by date)](https://img.shields.io/github/v/release/LGiki/PoDownloader?style=flat-square)](https://github.com/LGiki/PoDownloader/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/LGiki/PoDownloader)](https://goreportcard.com/report/github.com/LGiki/PoDownloader)

🎙️⬇️ PoDownloader = **Po**dcast **Downloader**, 一个用于下载播客的命令行工具.

这个工具会下载播客的RSS、播客封面图片、单集音频文件、单集封面图片和单集的Shownotes。

[English Version](https://github.com/LGiki/PoDownloader/blob/master/README.md)

# 截图

![](https://raw.githubusercontent.com/LGiki/PoDownloader/master/screenshot/screenshot.png)

# 安装

## 从Releases下载

你可以直接从 [Releases](https://github.com/LGiki/PoDownloader/releases) 页面下载最新版本。

## 从源码编译

确保 [go](https://golang.org/) 已经正确安装在系统中。

```bash
git clone https://github.com/LGiki/PoDownloader.git
cd PoDownloader
go mod download
go build -o podownloader ./cmd
```

执行完以上命令之后，你可以在目录下找到名为 `podownloader` 的二进制文件。

# 用法

## 从OPML文件下载播客

大多数播客APP都支持导出播客列表为[OPML格式](https://en.wikipedia.org/wiki/OPML)。

一个OPML文件的样例如下：

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

这个OPML文件包含3个播客（每个`outline`标签都是一个播客），其中的`xmlUrl`属性就是播客的RSS链接。

使用以下命令从OPML文件中下载播客：

```bash
podownloader download --opml /path/to/opml_file.xml
```

## 从RSS链接列表文件下载播客

RSS链接列表文件是一个文本文件，每一行一个播客的RSS链接，例如：

```
https://exmaple.org/podcast1/rss.xml
https://example.org/podcast2/rss.xml
https://example.org/podcast3/rss.xml
```

使用以下命令从RSS链接列表文件中下载播客：

```bash
podownloader download --list /path/to/rss_list_file.txt
```

## 从RSS链接下载播客

```
podownloader download --rss https://example.org/podcast/rss.xml
```

# 下载选项

通过`-h`或`--help`查看所有的选项及帮助信息。

## 输出文件夹

通过`-o`或`--output`来指定输出文件夹路径，默认输出目录是`./podcast`。

## 用户代理 (User agent)

通过`-u`或`--ua`来指定用户代理，默认的用户代理是`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36`。

## 下载线程数

通过`-t`或`--thread`来设定下载线程数，默认下载线程数为`3`。

## 日志文件夹

通过`--log`参数来指定日志文件夹，如果指定了`--log`参数，日志文件将会保存到指定的日志文件夹中；如果未指定`--log`参数，将不会生成日志文件。

`--log`参数默认为空，即不生成任何日志文件。

# 配置文件

如果你不想每次运行程序的时候都手动指定一堆参数，你可以将参数写入到配置文件中，程序将会自动从配置文件加载参数。

你可以通过`-c`或`--config`来指定配置文件的路径：

```bash
podownloader --config ~/.podownloader.json
```

默认配置文件路径是：`$PWD/.podownloader`。

支持的配置文件格式有：`json`、`toml`、`yaml`、`yml`、`properties`、`props`、`prop`、`hcl`、`dotenv`、`env`和`ini`。

一个JSON格式的配置文件样例如下：

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

你可以在 [config_template](https://github.com/LGiki/PoDownloader/tree/master/config_template) 目录下找到更多配置文件模板。

# 下载目录结构

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

# 许可

[Apache-2.0](https://github.com/LGiki/PoDownloader/blob/master/LICENSE)