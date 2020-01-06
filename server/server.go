package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"../protector"
)

type Connection struct {
	conn       net.Conn
	user       string
	protection *protector.SessionProtector
	currentKey string
}

type Message struct {
	text string
	user string
}

var (
	connections []Connection
	message     []Message
)

func disconectHandler(err error, conn net.Conn) bool {
	if err != nil {
		conn.Close()
		return true
	} else {
		return false
	}
}

func get(index int) {
	buf := make([]byte, 256)
	for {

		size, err := bufio.NewReader(connections[index].conn).Read(buf)
		disconectHandler(err, connections[index].conn)
		connections[index].currentKey = connections[index].protection.Next_session_key(connections[index].currentKey)
		fmt.Println(connections[index].user + "'s key:" + strings.Split(string(buf[:size]), "\n")[0] + " : " + "server key:" + connections[index].currentKey)

		if connections[index].currentKey != strings.Split(string(buf[:size]), "\n")[0] {
			disconectHandler(fmt.Errorf("key error"), connections[index].conn)
			return
		}
		if disconectHandler(err, connections[index].conn) {
			for iter := 0; iter < len(connections); iter++ {
				if connections[iter] == connections[index] {
					connections = append(connections[:iter], connections[iter+1:]...)
					continue
				}
				connections[iter].conn.Write([]byte(fmt.Sprintf("%s disconected \n", connections[index].user)))
				fmt.Println(fmt.Sprintf("%s disconected\n", connections[index].user))
				return
			}
		} else {
			fmt.Println(strings.Split(string(buf[:size]), "\n")[1])
			message = append(message, Message{user: connections[index].user, text: strings.Split(string(buf[:size]), "\n")[1]})
		}
	}
}

func send() {
	for {
		for len(message) > 0 {
			deliver := message[0]
			for iter := 0; iter < len(connections); iter++ {
				if deliver.user != connections[iter].user {
					connections[iter].currentKey = connections[iter].protection.Next_session_key(connections[iter].currentKey)
					connections[iter].conn.Write([]byte(connections[iter].currentKey + "\n" + deliver.text))
				}
			}
			message = message[1:]
		}
	}
}
func bind(port int) net.Listener {
	temp, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", "localhost", port))
	return temp
}

func StartServer(port int, client_number int) {
	fmt.Println("SERVER HAS BEEN STARTED")
	listener := bind(port)
	go send()

	for len(connections) < client_number {
		loginBuf, buf := make([]byte, 256), make([]byte, 256)
		conn, _ := listener.Accept()
		size, _ := bufio.NewReader(conn).Read(buf)
		access := strings.Split(string(buf[:size]), "\n")

		fmt.Println(access)
		conn.Write([]byte("enter login : "))
		loginSize, err := bufio.NewReader(conn).Read(loginBuf)

		if !disconectHandler(err, conn) {
			fmt.Println(fmt.Sprintf("user %s has connected", string(loginBuf[:loginSize])))
			newConnection := Connection{conn: conn,
				user:       string(loginBuf[:loginSize]),
				protection: protector.NewSessionProtector(access[0]),
				currentKey: access[1]}
			connections = append(connections, newConnection)
			go get(len(connections) - 1)
		}
	}
}
