package views

const artistTab = `<li id="{{.Name}}" {{ if .Hide }}class="hide"{{ end }}>
  <table>
    {{ range .Data }}
    <tr>
      <td>{{linkArtist .Artist}}</td>
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
        {{linkArtist .Artist}}
        {{linkTrack .Artist .Album .Track}}
      </td>
      <td><span class="count">{{.Count}} plays</span></td>
    </tr>
    {{ end }}
  </table>
</li>`

const index = `<!DOCTYPE html>
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
      <h2>Recently played</h2>
      <table>
        {{ range .RecentlyPlayed }}
        <tr>
          <td>
            {{linkArtist .Artist}}
            {{linkTrack .Artist .Album .Track}}
          </td>
          <td><a href="/listen/{{.Timestamp}}"><time datetime="{{.Timestamp | datetime}}">{{.Timestamp | readable}}</time></a></td>
        </tr>
        {{ end }}
      </table>
      <a class="more" href="/played">More</a>
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
