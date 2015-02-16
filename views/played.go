package views

import (
	"fmt"
	"html/template"
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
		if d.Hours() == 0 {
			return fmt.Sprintf("%d mintues ago", d.Minutes())
		}
		return fmt.Sprintf("%d hours ago", d.Hours())
	}

	if d.Hours() < 48 && n.Weekday() == u.Weekday()+1 {
		return u.Format("Yesterday at 15:04pm")
	}

	if n.Year() == u.Year() {
		return u.Format("02 Jan 15:04pm")
	}

	return u.Format("02 Jan 2006")
}

var Played = template.Must(template.New("played").Funcs(template.FuncMap{
	"datetime": datetime,
	"readable": readable,
}).Parse(played))

const played = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>trobble</title>
    <style>
      body {
          font: 16px/1.3em Georgia;
          margin: 2rem;
      }

      header {
          margin: 1rem 0 2rem;
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
    </style>
  </head>
  <body>
    <header>
      <h1>trobble</h1>
      <div>{{.TotalPlays}} plays</div>
    </header>

    <section class="recently-played">
      <h2>Recently played</h2>
      <ol>
        {{ range .RecentlyPlayed }}
        <li>
          <span class="artist">{{.Artist}}</span>
          <span class="track">{{.Track}}</span>
          <time datetime="{{.Timestamp | datetime}}">{{.Timestamp | readable}}</time>
        </li>
        {{ end }}
      </ol>
    </section>

    <section class="top-artists">
      <h2>Top Artists</h2>

      <nav class="tabs-choice">
        <h3 class="selected" data-id="tab-artists-overall">Overall</h3>
        <h3 data-id="tab-artists-year">Year</h3>
        <h3 data-id="tab-artists-month">Month</h3>
        <h3 data-id="tab-artists-week">Week</h3>
      </nav>
      <ul class="tabs-content">
        <li id="tab-artists-overall">
          <ol>
            {{ range .TopArtists.Overall }}
            <li>
              <span class="artist">{{.Artist}}</span>
              <span class="count">{{.Count}} plays</span>
            </li>
            {{ end }}
          </ol>
        </li>
        <li class="hide" id="tab-artists-year">
          <ol>
            {{ range .TopArtists.Year }}
            <li>
              <span class="artist">{{.Artist}}</span>
              <span class="count">{{.Count}} plays</span>
            </li>
            {{ end }}
          </ol>
        </li>
        <li class="hide" id="tab-artists-month">
          <ol>
            {{ range .TopArtists.Month }}
            <li>
              <span class="artist">{{.Artist}}</span>
              <span class="count">{{.Count}} plays</span>
            </li>
            {{ end }}
          </ol>
        </li>
        <li class="hide" id="tab-artists-week">
          <ol>
            {{ range .TopArtists.Week }}
            <li>
              <span class="artist">{{.Artist}}</span>
              <span class="count">{{.Count}} plays</span>
            </li>
            {{ end }}
          </ol>
        </li>
      </ul>
    </section>

    <section class="top-tracks">
      <h2>Top Tracks</h2>

      <nav class="tabs-choice">
        <h3 class="selected" data-id="tab-tracks-overall">Overall</h3>
        <h3 data-id="tab-tracks-year">Year</h3>
        <h3 data-id="tab-tracks-month">Month</h3>
        <h3 data-id="tab-tracks-week">Week</h3>
      </nav>
      <ul class="tabs-content">
        <li id="tab-tracks-overall">
          <ol>
            {{ range .TopTracks.Overall }}
            <li>
              <span class="artist">{{.Artist}}</span>
              <span class="track">{{.Track}}</span>
              <span class="count">{{.Count}} plays</span>
            </li>
            {{ end }}
          </ol>
        </li>
        <li class="hide" id="tab-tracks-year">
          <ol>
            {{ range .TopTracks.Year }}
            <li>
              <span class="artist">{{.Artist}}</span>
              <span class="track">{{.Track}}</span>
              <span class="count">{{.Count}} plays</span>
            </li>
            {{ end }}
          </ol>
        </li>
        <li class="hide" id="tab-tracks-month">
          <ol>
            {{ range .TopTracks.Month }}
            <li>
              <span class="artist">{{.Artist}}</span>
              <span class="track">{{.Track}}</span>
              <span class="count">{{.Count}} plays</span>
            </li>
            {{ end }}
          </ol>
        </li>
        <li class="hide" id="tab-tracks-week">
          <ol>
            {{ range .TopTracks.Week }}
            <li>
              <span class="artist">{{.Artist}}</span>
              <span class="track">{{.Track}}</span>
              <span class="count">{{.Count}} plays</span>
            </li>
            {{ end }}
          </ol>
        </li>
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
