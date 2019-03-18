package views

const listen = `<!DOCTYPE html>
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

    <section>
      <h1>{{linkArtist .Artist}} / {{.Track}}</h1>
      <p>From the album {{linkAlbum .Artist .Album}}</p>
      <p>played <time class="unright" datetime="{{ datetime .Timestamp }}">{{ readable .Timestamp }}</time></p>
    </section>
  </body>
</html>`
