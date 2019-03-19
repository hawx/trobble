package views

import (
	"fmt"
	"html/template"
	"io"
	"net/url"
	"time"
)

func Parse(glob string) (*template.Template, error) {
	return template.New("index").Funcs(template.FuncMap{
		"datetime":   datetime,
		"kitchen":    kitchen,
		"readable":   readable,
		"pair":       pair,
		"percent":    percent,
		"linkTrack":  linkTrack,
		"linkAlbum":  linkAlbum,
		"linkArtist": linkArtist,
	}).ParseGlob(glob)
}

func linkTrack(artist, album, track string) template.HTML {
	return template.HTML(fmt.Sprintf(`<a class="track" href="/artist/%s/%s/%s">%s</a>`,
		url.QueryEscape(artist),
		url.QueryEscape(album),
		url.QueryEscape(track),
		track))
}

func linkAlbum(artist, album string) template.HTML {
	return template.HTML(fmt.Sprintf(`<a class="album" href="/artist/%s/%s">%s</a>`,
		url.QueryEscape(artist),
		url.QueryEscape(album),
		album))
}

func linkArtist(artist string) template.HTML {
	return template.HTML(fmt.Sprintf(`<a class="artist" href="/artist/%s">%s</a>`,
		url.QueryEscape(artist),
		artist))
}

func datetime(t int64) string {
	return time.Unix(t, 0).UTC().Format(time.RFC3339)
}

func kitchen(t int64) string {
	return time.Unix(t, 0).UTC().Format(time.Kitchen)
}

func readable(t int64) string {
	n := time.Now().UTC()
	u := time.Unix(t, 0).UTC()
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

func percent(a, b int) int {
	return int(float64(a) / float64(b) * 100)
}

type Pair struct {
	Name string
	Hide bool
	Data interface{}
}

func pair(name string, show bool, data interface{}) *Pair {
	return &Pair{name, !show, data}
}

type wrappedTemplate struct {
	t *template.Template
	n string
}

func (w *wrappedTemplate) Execute(wr io.Writer, data interface{}) error {
	return w.t.ExecuteTemplate(wr, w.n, data)
}
