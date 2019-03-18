package views

const played = `<!DOCTYPE html>
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

    <section class="recently-played">
      {{ range .Tracks }}
      <h2 style="margin-left: -.5rem">{{.Date.Day}} {{.Date.Month}} {{.Date.Year}}</h2>

      <table>
        {{ range .Tracks }}
        <tr>
          <td>
            {{linkArtist .Artist}}
            {{linkTrack .Artist .Album .Track}}
          </td>
          <td><a href="/listen/{{.Timestamp}}"><time datetime="{{.Timestamp | datetime}}">{{.Timestamp | kitchen}}</time></a></td>
        </tr>
        {{ end }}
      </table>
      {{ end }}

      <a class="more" href="/played?from={{.MoreTime}}">More</a>
    </section>
  </body>
</html>`
