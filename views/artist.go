package views

const artist = `<!DOCTYPE html>
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
      <h1>{{ .Name }}</h1>
      <ul class="graph">
        {{ range .Plays }}
        <li style="height: {{ percent . $.MaxPlays }}%;"></li>
        {{ end }}
      </ul>
    </section>

    <section>
      <h2>Top Tracks</h2>
      <table>
        {{ range .Tracks }}
        <tr>
          <td><span class="track">{{.Track}}</span></td>
          <td><span class="count">{{.Count}} plays</span></td>
        </tr>
        {{ end }}
      </table>
    </section>
  </body>
</html>`
