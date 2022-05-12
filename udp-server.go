package breaker

import (
	"context"
	"net"
)

type UDPServer interface {
	Listen() error
}

type HandlerFunc func(ctx context.Context, request Request, data []byte)

type Request struct {
	FromAddr *net.UDPAddr
	Conn     *net.UDPConn
}

func NewUDPServer(addr string, port int, handlerFunc HandlerFunc) UDPServer {
	return &udpServer{
		udpAddr: &net.UDPAddr{
			IP:   net.ParseIP(addr),
			Port: port,
		},
		buffer:  make([]byte, 1024),
		handler: handlerFunc,
	}
}

type udpServer struct {
	udpAddr *net.UDPAddr
	buffer  []byte
	handler HandlerFunc
}

func (u *udpServer) Listen() error {
	conn, err := net.ListenUDP("udp", u.udpAddr)
	if err != nil {
		return err
	}

	return u.endlessRun(conn)
}

func (u *udpServer) endlessRun(conn *net.UDPConn) error {
	ctx := context.Background()
	for {
		n, addr, err := conn.ReadFromUDP(u.buffer)
		if err != nil {
			return err
		}
		go func() {
			u.handler(ctx, Request{
				FromAddr: addr,
				Conn:     conn,
			}, u.buffer[:n])
		}()
	}
}
