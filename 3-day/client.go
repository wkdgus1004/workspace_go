package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	connect, err := net.Dial("tcp", "127.0.0.1:5000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// some event happens

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	myIp := ""
	for _, addr := range addrs {
		ipv4 := addr.To4().String()
		fmt.Println("IPv4: ", ipv4)
		myIp = ipv4

	}
	sendString := myIp + " : send hihi"

	connect.Write([]byte(sendString))
	fmt.Println("Send Data : ", sendString)

	data := make([]byte, 4096)

	for {
		n, err := connect.Read(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(data[:n]))

		time.Sleep(1 * time.Second)
	}

}
