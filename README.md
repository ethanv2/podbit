# Podbit - **Podboat Improved**

![Podbit Player](https://www.ethanjmarshall.co.uk/wp-content/uploads/2022/04/screenshot-220408-2219-04.png)

[![AUR Release](https://img.shields.io/aur/version/podbit?color=1793d1&label=podbit&logo=arch-linux)](https://aur.archlinux.org/packages/podbit/)
[![Makefile CI](https://github.com/ethanv2/podbit/actions/workflows/makefile.yml/badge.svg)](https://github.com/ethanv2/podbit/actions/workflows/makefile.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ethanv2/podbit)](https://goreportcard.com/report/github.com/ethanv2/podbit)

Podbit is a replacement for ``newsboat``'s standard podboat tool for listening to podcasts. It is minimal, performant and abides by the UNIX-philosophy, with an ncurses terminal user interface and vi-like keybindings.

Podbit runs on Linux and the BSDs.

## Features

* Reads the standard ``newsboat`` queue file to integrate seamlessly
* Automatic podcast downloading, including in parallel
* Podcast playing using ``mpv``
* Podcast caching and automatic deletion once finished
* Vi-like "hjkl" to navigate the interface

## Requirements

Podbit is written in Go. So, to build, you will need a copy of the Go command line tool. In addition, you will need:

* *ncurses* development libraries and headers, including wide character support (*libncusesw*)
* A normal install of ``mpv``
* A copy of GNU Make
* Newsboat to enqueue podcasts - *(optional)*
* A YouTube downloader tool, such as ``youtube-dl`` or ``yt-dlp``, to download YouTube podcasts - *(optional)*

Because of security issues in the Go tool, the provided Makefile must be used instead of simply ``go build``.
