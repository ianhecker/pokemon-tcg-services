package networking

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
	url string
}

func MakeLogger(logger *zap.SugaredLogger) Logger {
	return Logger{logger, ""}
}

var w = func(s string) string { return fmt.Sprintf("client: %s", s) }

type elapsed = time.Duration

func (l *Logger) Set(url string) {
	l.url = url
}

func (l *Logger) Requesting() {
	l.Infow(w("new request"), "url", l.url)
}

func (l *Logger) RequestError(err error) {
	l.Errorw(w("request error"), "url", l.url, "err", err)
}

func (l *Logger) ContextIssue(e elapsed, err error) {
	l.Infow(w("context issue"), "url", l.url, "elapsed", e, "err", err)
}

func (l *Logger) ResponseError(e elapsed, err error) {
	l.Errorw(w("response error"), "url", l.url, "elapsed", e, "err", err)
}

func (l *Logger) UnexpectedStatus(e elapsed, status int, err error) {
	l.Warnw(w("unexpected status"), "url", l.url, "elapsed", e, "status", status, "err", err)
}

func (l *Logger) Success(e elapsed) {
	l.Infow(w("response success"), "url", l.url, "elapsed", e)
}
