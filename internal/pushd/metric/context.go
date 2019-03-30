package metric

import "github.com/tidwall/redcon"

func ManagerWithContext(conn redcon.Conn, manager *Manager) {
	conn.SetContext(manager)
}

func ManagerFromContext(conn redcon.Conn) *Manager {
	return conn.Context().(*Manager)
}
