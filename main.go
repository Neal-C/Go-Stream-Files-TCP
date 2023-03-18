//lint:file-ignore ST1006 heh...

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct {

}

func (self *FileServer) start(){
	listener, err := net.Listen("tcp", ":3000");
	if err != nil {
		log.Fatal(err);
	}

	for {

		connection, err := listener.Accept();

		if err != nil {
			log.Fatal(err);
		}

		go self.readLoop(connection);
	}
}

func (self *FileServer) readLoop(connection net.Conn){
	buffer := new(bytes.Buffer);

	for {
		var size int64;
		binary.Read(connection, binary.LittleEndian, &size);
		n, err := io.CopyN(buffer, connection, size);
		if err != nil {
			log.Fatal(err);
		}
		
		fmt.Println(buffer.Bytes());
		fmt.Printf("received %d bytes over the network \n", n )
	}
}

func sendFile(size int) error {
	
	file := make([]byte, size);

	_ , err := io.ReadFull(rand.Reader, file);

	if err != nil {
		return err;
	}

	connection, err := net.Dial("tcp", ":3000");

	if err != nil {
		return err;
	}

	binary.Write(connection, binary.LittleEndian, int64(size));
	n, err := io.CopyN(connection, bytes.NewReader(file), int64(size));

	if err != nil {
		return err;
	}

	// n, err := connection.Write(file);
	
	// if err != nil {
	// 	return err;
	// }

	fmt.Printf("written %d bytes over the network \n", n);

	return nil;
}


func main(){
	go func(){
		time.Sleep(4 * time.Second);
		sendFile(4000)
	}();
	server := &FileServer{};

	server.start();
}