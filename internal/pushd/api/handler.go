package api

import (
	"strconv"

	"github.com/kismia/pushd/internal/pushd/label"
	"github.com/kismia/pushd/internal/pushd/metric"
	"github.com/tidwall/redcon"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CounterAdd(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 3 {
		errWrongArgsNumber(conn, cmd.Args[0])
		return
	}

	value, err := strconv.ParseFloat(string(cmd.Args[2]), 64)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	labels, err := label.FromByteSlices(cmd.Args[3:])
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	counter, err := metric.ManagerFromContext(conn).Counter(string(cmd.Args[1]), labels)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	counter.Add(value)

	conn.WriteInt(1)
}

func (h *Handler) CounterInc(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		errWrongArgsNumber(conn, cmd.Args[0])
		return
	}

	labels, err := label.FromByteSlices(cmd.Args[2:])
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	counter, err := metric.ManagerFromContext(conn).Counter(string(cmd.Args[1]), labels)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	counter.Inc()

	conn.WriteInt(1)
}

func (h *Handler) GaugeSet(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 3 {
		errWrongArgsNumber(conn, cmd.Args[0])
		return
	}

	value, err := strconv.ParseFloat(string(cmd.Args[2]), 64)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	labels, err := label.FromByteSlices(cmd.Args[3:])
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	gauge, err := metric.ManagerFromContext(conn).Gauge(string(cmd.Args[1]), labels)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	gauge.Set(value)

	conn.WriteInt(1)
}

func (h *Handler) GaugeAdd(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 3 {
		errWrongArgsNumber(conn, cmd.Args[0])
		return
	}

	value, err := strconv.ParseFloat(string(cmd.Args[2]), 64)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	labels, err := label.FromByteSlices(cmd.Args[3:])
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	gauge, err := metric.ManagerFromContext(conn).Gauge(string(cmd.Args[1]), labels)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	gauge.Add(value)

	conn.WriteInt(1)
}

func (h *Handler) GaugeInc(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		errWrongArgsNumber(conn, cmd.Args[0])
		return
	}

	labels, err := label.FromByteSlices(cmd.Args[2:])
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	gauge, err := metric.ManagerFromContext(conn).Gauge(string(cmd.Args[1]), labels)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	gauge.Inc()

	conn.WriteInt(1)
}

func (h *Handler) GaugeDec(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		errWrongArgsNumber(conn, cmd.Args[0])
		return
	}

	labels, err := label.FromByteSlices(cmd.Args[2:])
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	gauge, err := metric.ManagerFromContext(conn).Gauge(string(cmd.Args[1]), labels)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	gauge.Dec()

	conn.WriteInt(1)
}

func (h *Handler) GaugeSub(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 3 {
		errWrongArgsNumber(conn, cmd.Args[0])
		return
	}

	value, err := strconv.ParseFloat(string(cmd.Args[2]), 64)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	labels, err := label.FromByteSlices(cmd.Args[3:])
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	gauge, err := metric.ManagerFromContext(conn).Gauge(string(cmd.Args[1]), labels)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	gauge.Add(value)

	conn.WriteInt(1)
}

func (h *Handler) HistogramObserve(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 3 {
		errWrongArgsNumber(conn, cmd.Args[0])
		return
	}

	value, err := strconv.ParseFloat(string(cmd.Args[2]), 64)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	labels, err := label.FromByteSlices(cmd.Args[3:])
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	observer, err := metric.ManagerFromContext(conn).Histogram(string(cmd.Args[1]), labels)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	observer.Observe(value)

	conn.WriteInt(1)
}

func (h *Handler) SummaryObserve(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 3 {
		errWrongArgsNumber(conn, cmd.Args[0])
		return
	}

	value, err := strconv.ParseFloat(string(cmd.Args[2]), 64)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	labels, err := label.FromByteSlices(cmd.Args[3:])
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	observer, err := metric.ManagerFromContext(conn).Summary(string(cmd.Args[1]), labels)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	observer.Observe(value)

	conn.WriteInt(1)
}

func (h *Handler) Ping(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("PONG")
}

func (h *Handler) Quit(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("OK")
	conn.Close()
}

func errWrongArgsNumber(conn redcon.Conn, command []byte) {
	conn.WriteError("ERR wrong number of arguments for '" + string(command) + "' command")
}
