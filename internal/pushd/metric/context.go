package metric

import "github.com/tidwall/redcon"

func ServiceWithContext(conn redcon.Conn, service *Service) {
	conn.SetContext(service)
}

func ServiceFromContext(conn redcon.Conn) *Service {
	return conn.Context().(*Service)
}
