package metric

import "github.com/tidwall/redcon"

func ServiceWithContext(conn redcon.Conn, manager *Service) {
	conn.SetContext(manager)
}

func ServiceFromContext(conn redcon.Conn) *Service {
	return conn.Context().(*Service)
}
