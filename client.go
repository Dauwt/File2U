package main

import (
	"File2U/utils"
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	config := &tls.Config{
		InsecureSkipVerify: true, // Validar o certificado em ambiente real
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:8443", config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Ligação estabelecida com sucesso.")

	fileNames, err := handleFiles()
	if err != nil {
		panic(err)
	}

	var protocol []string
	var protocolErr error
	protocol, protocolErr = handleInitialFileProtocol(fileNames)
	if protocolErr != nil {
		panic(protocolErr)
	}

	//fmt.Println(protocol)
	for _, line := range protocol {
		_, err := fmt.Fprintf(conn, "%s\n", line)
		if err != nil {
			fmt.Println("Error while sending file protocol:", err)
		}
		fmt.Println(line)
	}

	for _, fileName := range fileNames {
		fileProtocol, fileProtocolErr := handleFileProtocol(fileName)
		if fileProtocolErr != nil {
			panic(fileProtocolErr)
		}

		for _, line := range fileProtocol {
			_, err := fmt.Fprintf(conn, "%s\n", line)
			if err != nil {
				fmt.Println("Error while sending file protocol:", err)
			}
			fmt.Println(line)
		}
		sendFile(fileName, conn)
	}
}

func handleFiles() ([]string, error) {
	fmt.Println("Path of the files you want to send:")
	var files []string
	for {
		reader := bufio.NewReader(os.Stdin)
		inputFiles, _ := reader.ReadString('\n')
		inputFiles = strings.TrimSpace(inputFiles)

		fmt.Println("inputFiles:", inputFiles)
		if inputFiles == "!stop" {
			break
		} else {
			info, err := os.Stat(inputFiles)
			if err != nil {
				fmt.Println("Error:", err)
				return nil, err
			}

			if err != nil {
				fmt.Println(fmt.Sprintf("Couldn't find the specified file '%s'", inputFiles))
				continue
			}

			if !info.IsDir() {
				files = append(files, filepath.Base(inputFiles))
			} else {
				root := inputFiles
				err := filepath.Walk(root, func(path string, subInfo os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if path == root {
						return nil
					}

					relPath, err := filepath.Rel(filepath.Dir(root), path)
					if err != nil {
						return err
					}

					if !subInfo.IsDir() {
						files = append(files, relPath)
						fmt.Println("Relative Path(File):", relPath)
					} else {
						fmt.Println("Relative Path(Directory):", relPath)
					}

					return nil
				})

				if err != nil {
					return nil, fmt.Errorf("Error while walking through directory:", err)
				}
			}
		}
	}
	fmt.Println("Files:", files)

	return files, nil
}

func handleInitialFileProtocol(fileNames []string) ([]string, error) {
	var fileProtocol []string
	totalSize, err := utils.GetTotalFileSizeBYTES(fileNames)
	if err != nil {
		return fileProtocol, fmt.Errorf("error while getting total file size: %d", err)
	}

	//fileProtocol = append(fileProtocol, "CLIENTIP ") // Adicionar isto depois
	fileProtocol = append(fileProtocol, "CLIENTNAME dauwt")
	fileProtocol = append(fileProtocol, fmt.Sprintf("FILESNUM %d", len(fileNames)))
	fileProtocol = append(fileProtocol, fmt.Sprintf("TOTALSIZE %d", totalSize))

	return fileProtocol, nil
}

func handleFileProtocol(fileName string) ([]string, error) {
	var fileProtocol []string

	fileSize, fileErr := utils.GetFileSizeBYTES(fileName)
	if fileErr != nil {
		return fileProtocol, fmt.Errorf("Error while getting file size: %d", fileErr)
	}

	fileProtocol = append(fileProtocol, fmt.Sprintf("FILENAME %s", fileName))
	fileProtocol = append(fileProtocol, fmt.Sprintf("FILESIZE %d", fileSize))

	return fileProtocol, nil
}

func sendFile(fileName string, conn net.Conn) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	n, err := io.Copy(conn, file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("File sent (%d bytes).\n", n)
}
