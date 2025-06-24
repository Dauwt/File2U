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

	clientNameLine, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("Client lost connection\n", err)
		} else {
			fmt.Println("Error while reading:", err)
		}
		return
	}

	//fmt.Println(clientNameLine)
	clientName := strings.TrimSpace(strings.TrimPrefix(clientNameLine, "CLIENTNAME "))
	fmt.Println(clientName)

	filesNumLine, _ := reader.ReadString('\n')
	//fmt.Println(filesNumLine)
	filesNum, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(filesNumLine, "FILESNUM ")))
	fmt.Println(filesNum)

	totalSizeLine, _ := reader.ReadString('\n')
	//fmt.Println(totalSizeLine)
	totalSize, _ := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(totalSizeLine, "TOTALSIZE ")), 10, 64)
	fmt.Println(totalSize)

	for i := 0; i < filesNum; i++ {
		fileNameLine, _ := reader.ReadString('\n')
		//fmt.Println(fileNameLine)
		fileName := strings.TrimSpace(strings.TrimPrefix(fileNameLine, "FILENAME "))
		fmt.Println(fileName)

		fileSizeLine, _ := reader.ReadString('\n')
		//fmt.Println(fileSizeLine)
		fileSize, _ := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(fileSizeLine, "FILESIZE ")), 10, 64)
		fmt.Println(fileSize)

		err := os.MkdirAll(filepath.Dir(filepath.Join("out", fileName)), 0755)
		if err != nil {
			panic(err)
		}

		outFile, err := os.Create(filepath.Join("out", fileName))
		if err != nil {
			panic(err)
		}

		written, err := io.CopyN(outFile, reader, fileSize)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client lost connection while sending files\n", err)
			} else {
				fmt.Println("Error occured while receiving binary file data:", err)
			}
			return
		}
		outFile.Close()
		fmt.Printf("Received %d bytes and saved it on %s\n", written, fileName)
	}
}
