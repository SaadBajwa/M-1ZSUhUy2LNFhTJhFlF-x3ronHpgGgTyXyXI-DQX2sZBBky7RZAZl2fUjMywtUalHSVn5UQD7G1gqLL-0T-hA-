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
	"time"
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

func connectToServer(msgchan chan string) {
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

		cond := 1
		for cond == 1 {
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
			time.Sleep(2)

			data := make([]byte, 30)
			n, err := conn.Read(data)
			if err != nil {
				fmt.Println("Error while receiving file index and data")
			}
			fileData := string(data[:n])

			file_data := strings.Split(fileData, ":")

			fileIndex := file_data[0]
			dataToSearch := file_data[1]

			fmt.Println("File:Data = ", file_data)

			go receiveMessage(conn, msgchan)
			search(conn, msgchan, fileIndex, dataToSearch)
		}
	}
}

func receiveMessage(conn net.Conn, msgchan chan string) {

	cond := 1
	for cond == 1 {
		recvdMessage := make([]byte, 30)
		n, err := conn.Read(recvdMessage)
		if err != nil {
		}
		msg := string(recvdMessage[:n])
		msgchan <- msg
		//		go search(conn, msgchan)
	}
}

func getData(fileToSearch string) []string {
	file := "./" + fileToSearch
	var data []string
	fmt.Println("get Data")
	f, _ := os.Open(file)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		data = append(data, line)
	}
	return data
}

func search(conn net.Conn, msgchan chan string, fileToSearch string, dataToSearch string) {
	fileToSearch = fileToSearch + ".txt"
	data := getData(fileToSearch)

	var rcvmsg string
	//	msg := <-msgchan
	//	rcvmsg := msg
	counter := 0
	msg := dataToSearch
	fmt.Println("Text to Search = ", msg)
	time.Sleep(1)
	for i := 0; i < len(data); i++ {
		if counter == 5 {
			msg1 := <-msgchan
			rcvmsg = msg1
			counter = 0
			fmt.Println("Message = ", rcvmsg)
			fmt.Println()
		}
		if rcvmsg == "001" {
			fmt.Println("Abort")
			break
		} else if data[i] == msg {
			fmt.Println("Found")
			conn.Write([]byte("found"))
			break
		} else if rcvmsg == "002" {
			conn.Write([]byte("alive"))
		}
		counter = counter + 1
	}
}

func main() {
	msgchan := make(chan string, 30)
	go connectToServer(msgchan)
	cond := 1
	for cond == 1 {
	}
}
