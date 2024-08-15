package udpserver

import (
	"bytes"
	"log"
	"net"
	"strings"
	"sync"
)

type Conn net.UDPConn
type Addr net.UDPAddr
type Handler func(conn *Conn, addr *Addr, message string)

func (c *Conn) ReadFromUDP(b []byte) (int, *Addr, error) {
	n, addr, err := (*net.UDPConn)(c).ReadFromUDP(b)
	return n, (*Addr)(addr), err
}

func (c *Conn) WriteToUDP(b []byte, addr *Addr) (int, error) {
	n, err := (*net.UDPConn)(c).WriteToUDP(b, (*net.UDPAddr)(addr))
	return n, err
}

func (a *Addr) String() string {
	return (*net.UDPAddr)(a).String()
}

type UDPServer struct {
	conn     *Conn
	handlers map[string]Handler
}

func NewUDPServer(address string) (*UDPServer, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	return &UDPServer{
		conn:     (*Conn)(conn),
		handlers: make(map[string]Handler),
	}, nil
}

func (s UDPServer) LocalAddr() *Addr {
	return (*Addr)((*net.UDPConn)(s.conn).LocalAddr().(*net.UDPAddr))
}

func (s *UDPServer) RegisterHandler(prefix string, handler Handler) {
	s.handlers[prefix] = handler
}

func (s *UDPServer) handleRequest(addr *Addr, message string) {
	for prefix, handler := range s.handlers {
		if strings.HasPrefix(message, prefix) {
			handler(s.conn, addr, message)
			return
		}
	}
	log.Printf("No handler for message: %s", message)
	s.conn.WriteToUDP([]byte("no handler for message"), addr)
}

func (s *UDPServer) Listen() {
	bufPool := sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	for {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()
		tempBuf := make([]byte, 128)
		n, addr, err := s.conn.ReadFromUDP(tempBuf)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			bufPool.Put(buf)
			continue
		}

		buf.Write(tempBuf[:n])

		go func(buf *bytes.Buffer, addr *Addr) {
			defer bufPool.Put(buf)

			message := buf.String()
			log.Printf("Received from %s: %s", addr.String(), message)

			s.handleRequest(addr, message)
		}(buf, addr)
	}
}

func send(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	_, err := conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		log.Printf("Error sending UDP message: %v", err)
		return
	}
	log.Printf("Sent to %s: %s", addr.String(), msg)
}

func (s *UDPServer) Send(addr *Addr, msg string) {
	send((*net.UDPConn)(s.conn), (*net.UDPAddr)(addr), msg)
}

func (s *UDPServer) Close() {
	s.conn.Close()
}
