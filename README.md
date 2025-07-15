# Tengo

Tengo is a complete reimplementation of the osu!arcade Tencho LAN server in Golang.

## Features

- [X] Login
- [X] Multiplayer
- [ ] Clean code

## Building and running

```
git clone https://github.com/ascenttree/tengo
cd tengo
go mod tidy
go run .
```

## How to connect

Any osu!arcade clients that are launched on PCs connected to the same network as the server will connect automatically thanks to the UDP broadcast server.

## Credits

[BurntSushi/toml](https://github.com/BurntSushi/toml) - Config parser

[Lekuruu](https://github.com/lekuruu) - Logging, and stream utils

[peppy](https://github.com/ppy) - Making osu!arcade
