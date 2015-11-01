package views

import (
	"fmt"
	"html/template"
	"io"
	"time"
)

var Index interface {
	Execute(io.Writer, interface{}) error
}

var Played interface {
	Execute(io.Writer, interface{}) error
}

func init() {
	var tmpl = template.Must(template.New("index").Funcs(template.FuncMap{
		"datetime": datetime,
		"kitchen":  kitchen,
		"readable": readable,
		"pair":     pair,
	}).Parse(index))
	tmpl = template.Must(tmpl.New("played").Parse(played))
	tmpl = template.Must(tmpl.New("artistTab").Parse(artistTab))
	tmpl = template.Must(tmpl.New("trackTab").Parse(trackTab))

	Index = &wrappedTemplate{tmpl, "index"}
	Played = &wrappedTemplate{tmpl, "played"}
}

func datetime(t int64) string {
	return time.Unix(t, 0).Format(time.RFC3339)
}

func kitchen(t int64) string {
	return time.Unix(t, 0).Format(time.Kitchen)
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

type wrappedTemplate struct {
	t *template.Template
	n string
}

func (w *wrappedTemplate) Execute(wr io.Writer, data interface{}) error {
	return w.t.ExecuteTemplate(wr, w.n, data)
}
