package main

import (
	"fmt"
	"net"
	"os"
)

func respond(conn net.Conn, status int, stringStatus string) error {
	buf := make([]byte, 1024)
	if _, err := conn.Read(buf); err != nil {
		return fmt.Errorf("can't read from conn: %w", err)
	}

	if _, err := fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n\r\n", status, stringStatus); err != nil {
		return fmt.Errorf("can't write to conn: %w", err)
	}

	return nil
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	if err := respond(conn, 200, "OK"); err != nil {
		fmt.Println("Error responding: ", err.Error())
		os.Exit(1)
	}
}
