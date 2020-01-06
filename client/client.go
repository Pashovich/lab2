package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"../protector"
)

var (
	currentKey string
	protection *protector.SessionProtector
)

func get(conn net.Conn, login string) {
	defer recovery()
	buf := make([]byte, 256)
	for {
		size, _ := bufio.NewReader(conn).Read(buf)
		currentKey = protection.Next_session_key(currentKey)
		message := strings.Split(string(buf[:size]), "\n")
		if message[0] != currentKey {
			disconectHandler(fmt.Errorf("shutting down"), conn)
		}
		fmt.Println(message[1])
	}
}

func send(conn net.Conn, login string) {
	defer recovery()
	for {
		mess, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		currentKey = protection.Next_session_key(currentKey)
		conn.Write([]byte(currentKey + "\n" + fmt.Sprintf("%s:%s", login, mess)))
	}
}

func recovery() {
	if recv := recover(); recv != nil {
		os.Exit(1)
	}
}

func disconectHandler(err error, conn net.Conn) {
	if err != nil {
		conn.Close()
		panic(err)
	}
}

func StartClient(port int) {
	defer recovery()

	initHash := protector.Get_hash_str()
	currentKey = protector.Get_session_key()
	protection = protector.NewSessionProtector(initHash)

	buf := make([]byte, 256)
	addres, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "localhost", port))
	connection, err := net.DialTCP("tcp", nil, addres)

	disconectHandler(err, connection)
	connection.Write([]byte(initHash + "\n" + currentKey))

	var login string
	size, err := bufio.NewReader(connection).Read(buf)
	fmt.Println(string(buf[:size]))
	fmt.Scanln(&login)
	connection.Write([]byte(login))

	go send(connection, login)
	go get(connection, login)
}
