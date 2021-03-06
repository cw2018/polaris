package context

import (
	"encoding/json"
	"github.com/siddontang/polaris/session"
	"net/http"
	"strconv"
)

type Env struct {
	Request *http.Request
	Status  int

	//context for current request
	Ctx Context

	//session for current request,
	Session *session.Session

	w        http.ResponseWriter
	finished bool
}

func NewEnv(w http.ResponseWriter, r *http.Request) *Env {
	e := new(Env)

	e.Request = r
	e.w = w

	e.Ctx = NewContext()

	e.finished = false
	e.Status = http.StatusOK

	return e
}

func (e *Env) Header() http.Header {
	return e.w.Header()
}

func (e *Env) SetContentType(tp string) {
	e.w.Header().Set("Content-type", tp)
}

func (e *Env) SetContentJson() {
	e.w.Header().Set("Content-type", "application/json; charset=utf-8")
}

func (e *Env) SetStatus(status int) {
	e.Status = status
}

func (e *Env) Write(v interface{}) {
	if e.finished {
		return
	}

	buf, err := json.Marshal(v)
	if err != nil {
		e.WriteError(http.StatusInternalServerError, err)
	} else {
		e.SetContentJson()

		e.write(buf)
	}
}

func (e *Env) WriteString(data string) {
	if len(e.Header().Get("Content-type")) == 0 {
		e.SetContentType("text/plain")
	}

	e.write([]byte(data))
}

func (e *Env) WriteBuffer(data []byte) {
	if len(e.Header().Get("Content-type")) == 0 {
		e.SetContentType("application/octet-stream")
	}

	e.write(data)
}

func (e *Env) write(data []byte) {
	if e.finished {
		return
	}

	e.finished = true

	e.w.Header().Set("Content-Length", strconv.Itoa(len(data)))

	e.w.WriteHeader(e.Status)
	e.w.Write(data)
}

func (e *Env) WriteError(status int, err error) {
	e.Status = status
	e.WriteString(err.Error())
}

func (e *Env) Redirect(url string, status int) {
	e.finished = true
	http.Redirect(e.w, e.Request, url, status)
}

func (e *Env) SetCookie(c *http.Cookie) {
	http.SetCookie(e.w, c)
}

func (e *Env) Finish() {
	if e.finished {
		return
	}

	e.finished = true

	e.w.WriteHeader(e.Status)
}

func (e *Env) IsFinished() bool {
	return e.finished
}
