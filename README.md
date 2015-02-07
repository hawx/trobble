# trobble

Catches last.fm scrobble requests and stores them in a database instead.

This does not proxy requests! It receives a subset of requests of the last.fm
api:

- auth.getMobileSession
- track.updateNowPlaying
- track.scrobble (single track only)

and returns dummy ok responses.

``` bash
$ go get github.com/hawx/trobble
$ trobble
...
```
