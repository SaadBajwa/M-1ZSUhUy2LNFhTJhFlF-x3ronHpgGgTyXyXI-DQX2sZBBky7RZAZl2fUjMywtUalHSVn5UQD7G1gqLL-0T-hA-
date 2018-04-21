package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getFiles(ext string) []string {
	pathS, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	var files []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}

func connectToServer(msgchan chan net.Conn) {
	fmt.Println("Connecting to server:")

	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	fmt.Println(conn)
	if err != nil {
		// handle error
	} else {
		fmt.Println("Connected to Server ")

		var files_index string
		files_index = ""

		files := getFiles(".txt")

		for i := 0; i < len(files)-1; i++ {

			file_index := strings.Split(files[i], ".")

			//files_index = append(files_index, file_index[0])
			//files_index = append(files_index, ":")

			var buf bytes.Buffer
			buf.WriteString(files_index)
			buf.WriteString(file_index[0])
			buf.WriteString(":")
			files_index = buf.String()

		}
		fmt.Println(files_index)
		conn.Write([]byte(files_index))
	}
}

func main() {

	msgchan := make(chan net.Conn)

	go connectToServer(msgchan)
	cond := 1
	for cond == 1 {
	}

	//files := getFiles(".txt")
	//fmt.Println(files)
}
