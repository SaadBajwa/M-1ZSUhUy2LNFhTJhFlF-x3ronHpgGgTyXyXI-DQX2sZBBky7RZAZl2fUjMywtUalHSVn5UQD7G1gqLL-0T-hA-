package main

//dir.go
import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	//"os"
)

type slave struct {
	ip        string // ip address of the client
	port      string // port at slave is listening
	status    bool   // whether slave is busy in finding a number or not
	connected bool   // slave is connected or not
}

func handleSlaveConnection(port string, slaveData [][]slave) {

}

func handleSearchRequest(port string) {

}

func main() {
	//parsing command line argument
	port := os.Args[1] // port at which server will listen
	// declearing variables
	var slaveData [][]slave
	var maxIndex int = -1

	// Starting listeing at server
	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
	}
	fmt.Println("Listening!!!")
	for {
		conn, err := ln.Accept()
		buff := make([]byte, 1000)
		n, err := conn.Read(buff)
		if err != nil {
		}
		fmt.Println(string(buff[:n]))
		index, _ := strconv.Atoi(string(buff[:n]))
		if index+1 > maxIndex {
			for ; maxIndex < index; maxIndex++ {
				slaveData = append(slaveData, make([]slave, 0))
			}
		}
		fmt.Println(len(slaveData))
		slaveData[index] = append(slaveData[index], slave{ip: conn.RemoteAddr().String()})
		fmt.Println(slaveData[index])

		conn.Write([]byte("test"))

		conn.Write([]byte("def"))

		time.Sleep(2)
		counter := 0
		for {
			if counter == 5 {
				conn.Write([]byte("def"))
				counter = 0
			} else {
				counter = counter + 1
			}

			//			recvdMessage := make([]byte, 30)
			//			n, err := conn.Read(recvdMessage)
			//			if err != nil {
			//				fmt.Println(n)
			//			}

		}
	}

}
