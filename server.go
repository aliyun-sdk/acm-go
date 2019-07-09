package acm

// server 服务地址
type server struct {
	host string
	port uint16
}

func newServer(host string, port uint16) server {
	return server{host: host, port: port}
}
