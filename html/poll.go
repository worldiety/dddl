package html

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/exp/slog"
	"net/http"
	"time"
)

func LongPollHandler(timeout time.Duration) http.HandlerFunc {
	var tmp [16]byte
	if _, err := rand.Read(tmp[:]); err != nil {
		panic(err)
	}

	instanceToken := hex.EncodeToString(tmp[:])
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("ping") {
			if _, err := w.Write([]byte(instanceToken)); err != nil {
				slog.Error("failed to write long poll token", slog.Any("err", err))
			}
			return
		}

		ctl := http.NewResponseController(w)
		if err := ctl.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
			slog.Error("longpoll: failed to set write deadline", slog.Any("err", err))
		}

		if err := ctl.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			slog.Error("longpoll: failed to set read deadline", slog.Any("err", err))
		}

		timer := time.NewTimer(timeout - time.Second)
		defer timer.Stop()

		select {
		case <-r.Context().Done():
			// timeout, cancellation etc
			slog.Info("long poll request interrupted")
		case <-timer.C:
			if _, err := w.Write([]byte(instanceToken)); err != nil {
				slog.Error("failed to write long poll token", slog.Any("err", err))
			}
			return
		}
	}
}
