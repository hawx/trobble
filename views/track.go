package views

const track = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width,initial-scale=1" />
    <title>{{.Title}}</title>
    <link rel="alternate" type="application/rss+xml" href="/feed" />
    <link rel="stylesheet" href="/styles.css" />
  </head>
  <body>
    {{ template "header" . }}

    <section class="hero">
      <h1>{{linkArtist .Artist}} / {{.Track}}</h1>
      <ul class="graph">
        {{ range .Plays }}
        <li style="height: {{ percent . $.MaxPlays }}%;"></li>
        {{ end }}
      </ul>
      <p>From the album {{linkAlbum .Artist .Album}}</p>
      <p>{{.TrackPlays}} plays</p>
    </section>
  </body>
</html>`
