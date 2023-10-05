package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func respond(conn net.Conn, status int, stringStatus string) error {
	if _, err := fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n\r\n", status, stringStatus); err != nil {
		return fmt.Errorf("can't write to conn: %w", err)
	}

	return nil
}

func checkStartLine(conn net.Conn) error {
	buf := make([]byte, 1024)
	if _, err := conn.Read(buf); err != nil {
		return fmt.Errorf("can't read from conn: %w", err)
	}

	bb := bytes.NewBuffer(buf)
	startLine, err := bb.ReadString('\n')
	if err != nil {
		return fmt.Errorf("can' read start line: %w", err)
	}

	split := strings.Split(startLine, " ")
	if len(split) != 3 {
		return errors.New("malformed start line")
	}

	status, stringStatus := 200, "OK"

	if split[1] != "/" {
		status, stringStatus = 404, "Not Found"
	}

	return respond(conn, status, stringStatus)
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

	if err := checkStartLine(conn); err != nil {
		fmt.Printf("can't check start line: %v", err)
		os.Exit(1)
	}
}
