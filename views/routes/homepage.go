package routes

import (
	c "main/views/components"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func HomePage() g.Node {
	return c.BasePage(
		"Login",
		c.Layout(Section(ID("intro"),
			P(
				g.Text("Hi, my name is"), Span(g.Text(" Dave.")),
			),
			H2(g.Text("I'm the CTO of a tech company.")),
			P(g.Text("Currently, I'm learning Go + CSS")),
		)),
	)
}
