package components

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func Layout(children ...g.Node) g.Node {
	return Div(children...)
}
