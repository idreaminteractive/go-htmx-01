package components

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func BasePage(title string, children ...g.Node) g.Node {
	return Doctype(
		HTML(
			Lang("en"),
			Head(
				TitleEl(g.Text(title)),
				Link(Rel("stylesheet"), Href("/static/css/pico.min.css")),
				Link(Rel("stylesheet"), Href("/static/css/custom.css")),
			),
			Body(children...),
		),
	)
}
