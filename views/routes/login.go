package routes

import (
	"main/views/components"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func LoginPage() g.Node {
	return components.BasePage(
		"Login",
		components.Layout(Section(
			FormEl(
				Div(c.Classes{"grid": true},
					Label(
						For("firstname"),
						g.Text("First name"),
						Input(
							Type("text"),
							ID("firstname"),
							Name("firstname"),
							Placeholder("First name"),
							Required(),
						),
					),
					Label(
						For("lastname"),
						g.Text("Last name"),
						Input(
							Type("text"),
							ID("lastname"),
							Name("lastname"),
							Placeholder("Last name"),
							Required(),
						),
					),
				),
				Label(For("email"), g.Text("Email Address")),
				Input(Type("email"), ID("email"), Name("email"), Placeholder("Email address"), Required()),
				Small(g.Text("We'll never share your email with anyone else.")),
				Button(Type("submit"), g.Text("Submit")),
			),
		)),
	)
}
