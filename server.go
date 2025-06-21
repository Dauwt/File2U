package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
	if err != nil {
		panic(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	ln, err := tls.Listen("tcp", ":8443", config)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("Server listening on :8443")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting the connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	clientNameLine, _ := reader.ReadString('\n')
	//fmt.Println(clientNameLine)
	clientName := strings.TrimSpace(strings.Split(clientNameLine, " ")[1])
	fmt.Println(clientName)

	filesNumLine, _ := reader.ReadString('\n')
	//fmt.Println(filesNumLine)
	filesNum, _ := strconv.Atoi(strings.TrimSpace(strings.Split(filesNumLine, " ")[1]))
	fmt.Println(filesNum)

	totalSizeLine, _ := reader.ReadString('\n')
	//fmt.Println(totalSizeLine)
	totalSize, _ := strconv.ParseInt(strings.TrimSpace(strings.Split(totalSizeLine, " ")[1]), 10, 64)
	fmt.Println(totalSize)

	for i := 0; i < filesNum; i++ {
		fileNameLine, _ := reader.ReadString('\n')
		//fmt.Println(fileNameLine)
		fileName := strings.TrimSpace(strings.Split(fileNameLine, " ")[1])
		fmt.Println(fileName)

		fileSizeLine, _ := reader.ReadString('\n')
		//fmt.Println(fileSizeLine)
		fileSize, _ := strconv.ParseInt(strings.TrimSpace(strings.Split(fileSizeLine, " ")[1]), 10, 64)
		fmt.Println(fileSize)

		outFile, err := os.Create(filepath.Join("out", fileName))
		if err != nil {
			panic(err)
		}

		written, err := io.CopyN(outFile, reader, fileSize)
		if err != nil {
			panic(err)
		}
		outFile.Close()
		fmt.Printf("Received %d bytes and saved it on %s\n", written, fileName)
	}
}
