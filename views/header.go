package views

const header = `<header>
  <h1><a href="/">{{.Title}}</a></h1>
  <div>{{.TotalPlays}} plays</div>
  {{ if .NowPlaying }}
  <div class="nowplaying">
    <span class="icon">&#x266B;</span>
    {{linkArtist .NowPlaying.Artist}}
    {{linkTrack .NowPlaying.Artist .NowPlaying.Album .NowPlaying.Track}}
  </div>
  {{ end }}
</header>`
