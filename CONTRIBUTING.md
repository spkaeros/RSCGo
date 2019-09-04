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
* [golang.org/x/crypto v0.0.0](https://golang.org/x/crypto) - For password hashing, using the SHA3/SHAKE256 functions provided.

Locally, my development environment runs amd64 Ubuntu 18.10 with Linux 5.0.0-23, so as a result the project is initially developed for and tested against amd64 Linux as the main target.  I also make and test builds for 386 and amd64 windows, using the mingw build tools, and 386 linux using a 32-bit GCC.  Any other build configurations may or may not work.  In practice, it should work fine anywhere that SQLite3 has native drivers.  If I dropped SQLite3 support for player and world data persistence, I would not need to use CGo at all and could support all of the BSDs and plan 9 and etc out of the box.  This is something that may eventually happen, but for now there is no plan to drop SQLite3 support.

## I can't seem to get my RSC client to connect to it properly!
The protocol which I'm implementing may not be exactly the same as your favorite RSClassic client.  I had taken a refactored 202 from eXemplar's old project from 2006 or 2007 and I modified it in some places to improve it in various respects.  This is the client which I test the server against.  You should be able to find it in another repository on my github account.  If it's not there, please let me know in discord or email!

## I want to contribute, but don't know how to write code!
That's okay!  I'll be making some utilities to allow less-savvy individuals to contribute to RSCGo, as well, via content updates, graphical contributions, etc.
I also would not object to donations.  After all, I have a family to raise, and this project takes a lot of time up for me.
Contact me through discord or [email](mailto:aeros.storkpk@gmail.com) to learn more.