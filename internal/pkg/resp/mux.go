package resp

import (
	"strings"

	"github.com/tidwall/redcon"
)

type Handler interface {
	ServeRESP(conn redcon.Conn, cmd redcon.Command)
}

type HandlerFunc func(conn redcon.Conn, cmd redcon.Command)

func (f HandlerFunc) ServeRESP(conn redcon.Conn, cmd redcon.Command) {
	f(conn, cmd)
}

type ServeMux struct {
	handlers map[string]Handler
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		handlers: make(map[string]Handler),
	}
}

func (m *ServeMux) HandleFunc(command string, handler func(conn redcon.Conn, cmd redcon.Command)) {
	m.Handle(command, HandlerFunc(handler))
}

func (m *ServeMux) Handle(command string, handler Handler) {
	m.handlers[command] = handler
}

func (m *ServeMux) ServeRESP(conn redcon.Conn, cmd redcon.Command) {
	command := strings.ToLower(string(cmd.Args[0]))

	if handler, ok := m.handlers[command]; ok {
		handler.ServeRESP(conn, cmd)
	} else {
		conn.WriteError("ERR unknown command '" + command + "'")
	}
}
