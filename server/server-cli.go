package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	file, err := os.OpenFile("server-log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("can't start application server%e", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("can't close file %e", err)
		}
	}()

	log.SetOutput(file)
	log.Print("start application\n")
	host := "0.0.0.0"
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9999"
	}
	fmt.Println(port)

	err = startServer(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Printf("some error in server %e", err)
	}

}

func startServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("can't listen: %v", err)
	}
	defer func() {
		err = listener.Close()
		if err != nil {
			err = fmt.Errorf("can't close Listener %e", err)
		}
	}()
	fmt.Println("listen...")
	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept client %e", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("some connection...")
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("can't close connection %e", err)
		}
	}()

	reader, write := bufio.NewReader(conn), bufio.NewWriter(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("can't read command %e", err)
			return
		}
		dir, err := ioutil.ReadDir("downloads")
		if err != nil {
			log.Printf("can't read directory <downloads> %e", err)
			return
		}
		if str == "LIST\n" {
			for _, info := range dir {
				if !info.IsDir() {
					_, err = write.Write([]byte(info.Name() + "\n"))
					if err != nil {
						log.Printf("can't write to client %e", err)
						return
					}
					err = write.Flush()
					if err != nil {
						log.Printf("can't write to client %e", err)
						return
					}
				}
			}
			return
		} else if str == "DOWNLOAD\n" {
			fileName := ""
			fmt.Printf("download %s", fileName)
				fileName, err = reader.ReadString('\n')

				if err != nil {
					log.Printf("can't read command %e", err)
					return
				}
				if fileName == ""{
					return
				}

			fmt.Printf("download %s", fileName)
			for _, info := range dir {
				if !info.IsDir() && info.Name()+"\n" == fileName {
					fmt.Println("find")
					bytes, err := ioutil.ReadFile("downloads/"+info.Name())
					if err != nil {
						log.Printf("can't read file %e", err)
						return
					}
					_, err = write.Write(bytes)
					if err != nil {
						log.Printf("can't write to client %e", err)
						return
					}
					err = write.Flush()
					if err != nil {
						log.Printf("can't write to client %e", err)
						return
					}
					return
				}
			}
		} else if str == "UPLOAD\n" {

		}
		return
	}
}
