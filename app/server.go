package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"
)

var filesDir string

func respond(conn net.Conn, status int, stringStatus string, contentType string, body []byte) error {
	_, err := fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n", status, stringStatus)
	if err != nil {
		return fmt.Errorf("can't write to conn: %w", err)
	}

	if len(body) == 0 {
		_, err = fmt.Fprint(conn, "\r\n")
		return err
	}

	if _, err := fmt.Fprintf(conn, "Content-Type: %s\r\nContent-Length: %d\r\n\r\n%s", contentType, len(body), body); err != nil {
		return fmt.Errorf("can't write body to conn: %w", err)
	}

	return nil
}

func handleRequest(conn net.Conn) error {
	defer conn.Close()

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

	switch {
	case split[1] == "/":
		return respond(conn, 200, "OK", "text/plain", nil)

	case split[1] == "/user-agent":
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
				return respond(conn, 200, "OK", "text/plain", []byte(headerVal))
			}
		}

	case strings.HasPrefix(split[1], "/files/"):
		p := path.Join(filesDir, strings.TrimPrefix(split[1], "/files/"))
		f, err := os.Open(p)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return respond(conn, 404, "Not Found", "application/octet-stream", nil)
			}
			return fmt.Errorf("can't open file: %w", err)
		}

		content, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("can't read from file: %w", err)
		}

		return respond(conn, 200, "OK", "application/octet-stream", content)

	case strings.HasPrefix(split[1], "/echo/"):
		return respond(conn, 200, "OK", "text/plain", []byte(strings.TrimPrefix(split[1], "/echo/")))

	default:
		return respond(conn, 404, "Not Found", "text/plain", nil)
	}
}

func main() {
	if len(os.Args) >= 3 && os.Args[1] == "--directory" {
		filesDir = os.Args[2]
	}

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func() {
			if err := handleRequest(conn); err != nil {
				fmt.Printf("can't handle request: %v", err)
			}
		}()
	}
}
