package routes

import (
	"main/views/components"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func Home() g.Node {
	return components.BasePage(
		"Login",
		Section(c.Classes{"container": true},
			g.Text("Hi, my name is Dave and I'm the CTO of a tech company."),
		),
	)
}
