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
            <span class="artist"><a href="/artist/{{.Artist | urlquery }}">{{.Artist}}</a></span>
            <span class="track">{{.Track}}</span>
          </td>
          <td><time datetime="{{.Timestamp | datetime}}">{{.Timestamp | kitchen}}</time></td>
        </tr>
        {{ end }}
      </table>
      {{ end }}

      <a class="more" href="/played?from={{.MoreTime}}">More</a>
    </section>
  </body>
</html>`
