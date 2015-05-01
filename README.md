# trobble

Catches last.fm scrobble requests and stores them in a database instead.

This does not proxy requests! It receives a subset of requests of the last.fm
api:

- auth.getMobileSession
- track.updateNowPlaying
- track.scrobble (single track only)

and returns dummy ok responses.

Remember to set `--username`, `--api-key` and `--secret` as expected by the
client. `--api-key` and `--secret` are best set differently to your valid lastfm
details if possible.

``` bash
$ go get hawx.me/code/trobble
$ trobble --username "Me" --api-key "..." --secret "..."
...
```
