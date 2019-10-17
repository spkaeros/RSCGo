# Contibuting to RSCGo
Thanks for considering contributing to this wonderful FOSS project!

This project and everyone participating in it is governed by the [Code of Conduct](CODE_OF_CONDUCT.md).  By participating, you are expected to uphold this code.  Please report unacceptable behavior to aeros.storkpk@gmail.com

## What should I know before I get started?

RSCGo is written in Go, and development was started using Go 1.12.  It may work on 1.11, but I honestly am not sure if it does or doesn't and I'm not going to be offering support for anything under 1.12.  If you are compiling RSCGo, use the latest version of Go to do so to avoid any complications.

I make use of a number of free, open source Go libraries that I've found on Github and other places, in this project.
Currently, the libraries that RSCGo depends on are:
* [BurntSushi/toml v0.3.1](https://github.com/BurntSushi/toml) - For reading TOML configuration files
* [jessevdk/go-flags v1.4.0](https://github.com/jessevdk/go-flags) - For powerful command-line flag parsing
* [mattn/go-sqlite3 v1.11.0](https://github.com/mattn/go-sqlite3) - For reading/writing the player database and the game world database
* [golang.org/x/crypto v0.0.0](https://golang.org/x/crypto) - For password hashing, using the Argon2id functions provided.

Locally, my development environment runs amd64 Ubuntu 18.10 with Linux 5.0.0-23, so as a result the project is initially developed for and tested against amd64 Linux as the main target.  I also make and test builds for 386 and amd64 windows, using the mingw build tools, and 386 linux using a 32-bit GCC.  Any other build configurations may or may not work.  In practice, it should work fine anywhere that SQLite3 has native drivers.  If I dropped SQLite3 support for player and world data persistence, I would not need to use CGo at all and could support all of the BSDs and plan 9 and etc out of the box.  This is something that may eventually happen, but for now there is no plan to drop SQLite3 support.

## I can't seem to get my RSC client to connect to it properly!
The protocol which RSCGo is currently designed for is exactly identical to the one Jagex used in their 204 revision of RuneScape Classic.  I am making an effort to avoid making custom changes at the moment, even though I have thought of many improvements to various aspects of the networking protocol(e.g unsigned length for packet headers vs signed, null-terminated strings, etc..), I decided to keep it as identical to RSClassic 204 as possible, for software preservation purposes.  I may make a seperate branch with a variety of improvements later on, though.
You can get a working deobfuscated, refactored RSClassic 204 mudclient by looking through my repositories on GitHub.

## I want to contribute, but don't know how to write code!
That's okay!  I'll be making some utilities to allow less-savvy individuals to contribute to RSCGo, as well, via content updates, graphical contributions, etc.
I also would not object to donations.  After all, I have a family to raise, and this project takes a lot of time up for me.
Contact me through discord or [email](mailto:aeros.storkpk@gmail.com) to learn more.