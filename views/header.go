package views

const header = `<header>
  <h1><a href="/">{{.Title}}</a></h1>
  <div>{{.TotalPlays}} plays</div>
  {{ if .NowPlaying }}
  <div class="nowplaying">
    <span class="icon">&#x266B;</span>
    <span class="artist"><a href="/artist/{{.NowPlaying.Artist}}">{{.NowPlaying.Artist}}</a></span>
    <span class="track">{{.NowPlaying.Track}}</span>
  </div>
  {{ end }}
</header>`
