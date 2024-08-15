package udpserver_test

import (
	"net"
	"testing"

	"github.com/sintine-game/liman/pkg/udpserver"
)

func TestNewUDPServer(t *testing.T) {
	server, err := udpserver.NewUDPServer(":0")
	if err != nil {
		t.Fatalf("Failed to create UDP server: %v", err)
	}
	defer server.Close()
}

func TestRegisterHandler(t *testing.T) {
	server, err := udpserver.NewUDPServer(":0")
	if err != nil {
		t.Fatalf("Failed to create UDP server: %v", err)
	}
	defer server.Close()

	handler := func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {}
	server.RegisterHandler("test", handler)
}

func TestListen(t *testing.T) {
	server, err := udpserver.NewUDPServer(":0")
	if err != nil {
		t.Fatalf("Failed to create UDP server: %v", err)
	}
	defer server.Close()

	server.RegisterHandler("test", func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
		if message != "test" {
			t.Errorf("Expected message to be 'test', got '%s'", message)
		}
	})

	go server.Listen()

	conn, err := net.Dial("udp", server.LocalAddr().String())
	if err != nil {
		t.Fatalf("Failed to dial UDP server: %v", err)
	}

	_, err = conn.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Failed to write to UDP server: %v", err)
	}
}
