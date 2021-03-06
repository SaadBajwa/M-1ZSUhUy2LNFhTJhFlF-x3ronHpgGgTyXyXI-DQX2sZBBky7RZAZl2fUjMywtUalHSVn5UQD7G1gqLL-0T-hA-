package main

import (
	"bufio"
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

func connectToServer(s *string) {
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

			var buf bytes.Buffer
			buf.WriteString(files_index)
			buf.WriteString(file_index[0])
			buf.WriteString(":")
			files_index = buf.String()
		}
		fmt.Println(files_index)
		conn.Write([]byte(files_index))

		cond := 1
		count := 1
		var data []string

		for cond == 1 {
			fmt.Println("Waiting for HeartBeat")

			heartBeat := make([]byte, 3)
			n1, err1 := conn.Read(heartBeat)
			heartMsg := string(heartBeat[:n1])
			if err1 != nil {
				fmt.Println("Error while receiving Heart Beat Message")
			} else {
				if heartMsg == "002" {
					fmt.Println("Heart Beat Recv", heartMsg)
					conn.Write([]byte("alive"))
				}
			}

			data1 := make([]byte, 30)
			n, err := conn.Read(data1)
			if err != nil {
				fmt.Println("Error while receiving file index and data")
			}
			fileData := string(data1[:n])

			file_data := strings.Split(fileData, ":")

			fileIndex := file_data[0]
			dataToSearch := file_data[1]

			fmt.Println("File:Data = ", file_data)
			var fileToSearch string
			//			if count == 1 {
			fileToSearch = fileIndex + ".txt"
			//data = getData(fileToSearch)
			count = count + 1
			//			}
			//var data []string
			//go receiveMessage(conn, s)
			search(conn, fileToSearch, dataToSearch, data, s)
		}
	}
}

func receiveMessage(conn net.Conn, s *string) {

	cond := 1
	for cond == 1 {
		recvdMessage := make([]byte, 30)
		n, err := conn.Read(recvdMessage)
		if err != nil {
		}
		msg := string(recvdMessage[:n])
		//		msgchan <- msg
		*s = msg
	}
}

func search(conn net.Conn, fileToSearch string, dataToSearch string, data []string, s *string) { // []string {

	file := "./" + fileToSearch
	//var data []string
	fmt.Println("get Data")

	f, _ := os.Open(file)
	defer f.Close()
	i := 0

	var rcvmsg string

	msg := dataToSearch
	fmt.Println("Text to Search = ", msg)

	isFound := 0

	fmt.Println()
	fmt.Println()
	fmt.Println("Searching...")
	fmt.Println()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		i += 1
		fmt.Println(i)
		rcvmsg = *s
		fmt.Println(i)
		if rcvmsg == "001" {
			fmt.Println("Abort")
			break
		} else if line == msg {
			fmt.Println("Found")
			isFound = 1
			break
		} else if rcvmsg == "002" {
			conn.Write([]byte("alive"))
		}
		rcvmsg = ""
	}
	if isFound == 0 {
		fmt.Println("Not Found")
		conn.Write([]byte("not found"))
	} else {
		conn.Write([]byte("found"))
	}
	fmt.Println("Data has been found!!!")
}

func main() {
	s := ""
	connectToServer(&s)
	cond := 1
	for cond == 1 {
	}
}
