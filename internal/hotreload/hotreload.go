package hotreload

import (
	"time"

	"github.com/alexandrevicenzi/go-sse"
	"github.com/go-chi/chi/v5"
)

func HotReload(r chi.Router) {
	s := sse.NewServer(nil)
	defer s.Shutdown()

	r.Mount("/hmr/", s)

	go func() {
		for {
			s.SendMessage("/value", sse.SimpleMessage(time.Now().Format("2006/02/01/ 15:04:05")))
			time.Sleep(5 * time.Second)
		}
	}()

}
