package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:39998")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	n, err := conn.Write([]byte(`{"key": "abc"}`))
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("recv data: %s", string(buf[:n]))
}
