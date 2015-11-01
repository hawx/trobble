package views

import (
	"fmt"
	"html/template"
	"io"
	"time"
)

func datetime(t int64) string {
	return time.Unix(t, 0).Format(time.RFC3339)
}

func readable(t int64) string {
	n := time.Now()
	u := time.Unix(t, 0)
	d := n.Sub(u)

	if d.Hours() < 24 && n.Weekday() == u.Weekday() {
		if d.Hours() < 1 {
			return fmt.Sprintf("%d mintues ago", int(d.Minutes()))
		}
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	}

	if d.Hours() < 48 && n.Weekday() == u.Weekday()+1 {
		return u.Format("Yesterday at 15:04pm")
	}

	if n.Year() == u.Year() {
		return u.Format("02 Jan 15:04pm")
	}

	return u.Format("02 Jan 2006")
}

type Pair struct {
	Name string
	Hide bool
	Data interface{}
}

func pair(name string, show bool, data interface{}) *Pair {
	return &Pair{name, !show, data}
}

var Played interface {
	Execute(io.Writer, interface{}) error
}

func init() {
	var tmpl = template.Must(template.New("played").Funcs(template.FuncMap{
		"datetime": datetime,
		"readable": readable,
		"pair":     pair,
	}).Parse(played))
	tmpl = template.Must(tmpl.New("artistTab").Parse(artistTab))
	tmpl = template.Must(tmpl.New("trackTab").Parse(trackTab))

	Played = &wrappedTemplate{tmpl, "played"}
}

type wrappedTemplate struct {
	t *template.Template
	n string
}

func (w *wrappedTemplate) Execute(wr io.Writer, data interface{}) error {
	return w.t.ExecuteTemplate(wr, w.n, data)
}

const artistTab = `<li id="{{.Name}}" {{ if .Hide }}class="hide"{{ end }}>
  <table>
    {{ range .Data }}
    <tr>
      <td><span class="artist">{{.Artist}}</span></td>
      <td><span class="count">{{.Count}} plays</span></td>
    <tr/>
    {{ end }}
  </table>
</li>`

const trackTab = `<li id="{{.Name}}" {{ if .Hide }}class="hide"{{ end }}>
  <table>
    {{ range .Data }}
    <tr>
      <td>
        <span class="artist">{{.Artist}}</span>
        <span class="track">{{.Track}}</span>
      </td>
      <td><span class="count">{{.Count}} plays</span></td>
    </tr>
    {{ end }}
  </table>
</li>`

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
      <div>{{.TotalPlays}} plays</div>
      {{ if .NowPlaying }}
      <div class="nowplaying">
        <span class="icon">&#x266B;</span>
        <span class="artist">{{.NowPlaying.Artist}}</span>
        <span class="track">{{.NowPlaying.Track}}</span>
      </div>
      {{ end }}
    </header>

    <section class="recently-played">
      <h2>Recently played</h2>
      <table>
        {{ range .RecentlyPlayed }}
        <tr>
          <td>
            <span class="artist">{{.Artist}}</span>
            <span class="track">{{.Track}}</span>
          </td>
          <td><time datetime="{{.Timestamp | datetime}}">{{.Timestamp | readable}}</time></td>
        </tr>
        {{ end }}
      </table>
    </section>

    <section class="top-artists">
      <h2>Top Artists</h2>

      <nav class="tabs-choice">
        <h3 data-id="tab-artists-overall">Overall</h3>
        <h3 data-id="tab-artists-year">Year</h3>
        <h3 data-id="tab-artists-month" class="selected">Month</h3>
        <h3 data-id="tab-artists-week">Week</h3>
      </nav>
      <ul class="tabs-content">
        {{ template "artistTab" pair "tab-artists-overall" false .TopArtists.Overall }}
        {{ template "artistTab" pair "tab-artists-year"    false .TopArtists.Year    }}
        {{ template "artistTab" pair "tab-artists-month"   true  .TopArtists.Month   }}
        {{ template "artistTab" pair "tab-artists-week"    false .TopArtists.Week    }}
      </ul>
    </section>

    <section class="top-tracks">
      <h2>Top Tracks</h2>

      <nav class="tabs-choice">
        <h3 data-id="tab-tracks-overall">Overall</h3>
        <h3 data-id="tab-tracks-year">Year</h3>
        <h3 data-id="tab-tracks-month" class="selected">Month</h3>
        <h3 data-id="tab-tracks-week">Week</h3>
      </nav>
      <ul class="tabs-content">
        {{ template "trackTab" pair "tab-tracks-overall" false .TopTracks.Overall }}
        {{ template "trackTab" pair "tab-tracks-year"    false .TopTracks.Year    }}
        {{ template "trackTab" pair "tab-tracks-month"   true  .TopTracks.Month   }}
        {{ template "trackTab" pair "tab-tracks-week"    false .TopTracks.Week    }}
      </ul>
    </section>

    <script>
      function toggleClass(el, className, cond) {
          if (cond) {
              el.classList.add(className);
          } else{
              el.classList.remove(className);
          }
      }

      function activateTabsIn(el) {
          var tabsChoices = el.querySelectorAll('.tabs-choice h3');
          var tabsContents = el.querySelectorAll('.tabs-content > li');

          for (var choice of tabsChoices) {
              choice.onclick = (function(dataId) {
                  return function() {
                      for (var c of tabsChoices) {
                          toggleClass(c, "selected", c.dataset.id == dataId);
                      }

                      for (var content of tabsContents) {
                          toggleClass(content, "hide", content.id != dataId);
                      }
                  }
              })(choice.dataset.id);
          }
      }

      activateTabsIn(document.querySelector(".top-artists"));
      activateTabsIn(document.querySelector(".top-tracks"));
    </script>
  </body>
</html>`
