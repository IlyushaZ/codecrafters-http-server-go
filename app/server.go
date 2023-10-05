package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func respond(conn net.Conn, status int, stringStatus string, body []byte) error {
	_, err := fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n", status, stringStatus)
	if err != nil {
		return fmt.Errorf("can't write to conn: %w", err)
	}

	if len(body) == 0 {
		_, err = fmt.Fprint(conn, "\r\n")
		return err
	}

	if _, err := fmt.Fprintf(conn, "Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body); err != nil {
		return fmt.Errorf("can't write body to conn: %w", err)
	}

	return nil
}

func handleRequest(conn net.Conn) error {
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

	if split[1] == "/" {
		return respond(conn, 200, "OK", nil)
	}

	if split[1] == "/user-agent" {
		for {
			header, err := bb.ReadString('\n')
			if err != nil {
				return fmt.Errorf("can't read header: %w", err)
			}

			split := strings.SplitN(header, ":", 2)
			if len(split) != 2 {
				return errors.New("no user-agent given")
			}

			if strings.ToLower(split[0]) == "user-agent" {
				headerVal := strings.TrimSpace(split[1])
				return respond(conn, 200, "OK", []byte(headerVal))
			}
		}
	}

	splitPath := strings.SplitN(split[1], "/", 3)
	if len(splitPath) == 3 && splitPath[1] == "echo" {
		fmt.Println(splitPath[2])
		return respond(conn, 200, "OK", []byte(splitPath[2]))
	}

	return respond(conn, 404, "Not Found", nil)
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

	if err := handleRequest(conn); err != nil {
		fmt.Printf("can't check start line: %v", err)
		os.Exit(1)
	}
}
