package views

const played = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width,initial-scale=1" />
    <title>{{.Title}}</title>
    <link rel="alternate" type="application/rss+xml" href="/feed" />
    <style>
      body {
          font: 16px/1.3em Georgia;
          margin: 2rem;
      }

      header {
          margin: 1rem 0 2rem;
      }

      .nowplaying {
          margin-top: 1em;
      }

      section {
          margin: 1rem 0 2rem;
          border-top: 1px dotted;
      }

      h1 {
          font-size: 1.5rem;
      }

      h2 {
          font-size: 1rem;
          font-variant: small-caps;
      }

      .tabs-choice {
          margin-bottom: .5rem;
      }

      .tabs-choice h3 {
          margin: 0;
          color: #666;
      }

      .tabs-choice h3.selected {
          text-decoration: underline;
          color: black;
      }

      .tabs-content {
          margin: 0;
      }

      ol {
          list-style: none;
          padding-left: 0;
      }

      .tabs-choice {
          display: flex;
      }

      .tabs-choice h3 {
          cursor: pointer;
          font-size: 1rem;
      }

      .tabs-choice h3 + h3 {
          margin-left: 1em;
      }

      .tabs-content {
          list-style: none;
          padding: 0;
      }

      .tabs-content .hide {
          display: none;
      }

      section {
          max-width: 40rem;
      }

      .more {
          color: rgb(0, 20, 130);
          float: right;
      }

      table {
          border-collapse: collapse;
          width: 100%;
          table-layout: fixed;
      }

      td {
          padding: 0;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
      }

      tr td:last-child {
          width: 140px;
      }

      .artist {
          font-weight: bold;
      }

      .artist + .track {
          margin-left: .2rem;
      }

      time, .count {
          float: right;
          font-style: italic;
      }

      @media screen and (max-width: 30em) {
        table tr td:last-child { display: none; }
      }
    </style>
  </head>
  <body>
    <header>
      <h1>{{.Title}}</h1>
    </header>

    <section class="recently-played">
      {{ range .Tracks }}
      <h2 style="margin-left: -.5rem">{{.Date.Day}} {{.Date.Month}} {{.Date.Year}}</h2>

      <table>
        {{ range .Tracks }}
        <tr>
          <td>
            <span class="artist">{{.Artist}}</span>
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
