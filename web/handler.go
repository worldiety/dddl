package web

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/laher/mergefs"
	"github.com/vearutop/statigz"
	"github.com/worldiety/dddl/html"
	"github.com/worldiety/dddl/web/editor"
	"golang.org/x/exp/slog"
	"net/http"
	"time"
)

func StartServer(addr string, devMode bool, loader editor.Loader, saver editor.Saver, parser editor.Parser, linter editor.Linter) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	if devMode {
		slog.Info("devMode and hot-reload enabled")
		r.Get("/poll", html.LongPollHandler(60*time.Second))
	}

	r.Mount("/assets/", statigz.FileServer(mergefs.Merge(html.Assets, html.Tailwind).(mergefs.MergedFS), statigz.EncodeOnInit))
	r.HandleFunc("/", editor.Handler(devMode, loader, saver, parser, linter))

	if addr == "" {
		addr = "localhost:8080"
	}

	slog.Info("starting server on " + addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		return fmt.Errorf("cannot listen and server: %w", err)
	}

	return nil
}
