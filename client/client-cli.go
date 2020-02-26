package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	list := flag.Bool("list", false, "get directory list")
	upload := flag.Bool("upload", false, "get directory list")
	download := flag.Bool("download", false, "get directory list")

	file, err := os.OpenFile("client-log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("can't start application client%e", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("can't close file %e", err)
		}
	}()

	log.SetOutput(file)
	flag.Parse()

	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Fatalf("can't connect to 0.0.0.0:9999 %v", err)
		fmt.Println("Сервер не найден")
		return
	}
	log.Print("Start application client")
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("can't close connection %v", err)
		}
	}()
	reader, writer := bufio.NewReader(conn), bufio.NewWriter(conn)
	if *list {
		_, err = writer.Write([]byte("LIST\n"))
		if err != nil {
			log.Printf("can't write to localhost:9999 %v", err)
			fmt.Println("Не удалось отправить запрос серверу")
			return
		}
		err = writer.Flush()
		if err != nil {
			log.Printf("can't write to localhost:9999 %v", err)
			fmt.Println("Не удалось отправить запрос серверу")
			return
		}
		counter := 0
		fmt.Println("files:")
		for {
			readString, err := reader.ReadString('\n')
			if err == io.EOF {
				fmt.Printf("Total: %d files", counter)
				return
			}
			if err != nil {
				log.Printf("can't read frocm loalhost %v", err)
				fmt.Println("Не удалось получить ответ")
				return
			}
			counter++
			fmt.Println(readString)
		}
	}else if *download {
		_, err = writer.Write([]byte("DOWNLOAD\n"))
		if err != nil {
			log.Printf("can't write to localhost:9999 %v", err)
			fmt.Println("Не удалось отправить запрос серверу")
			return
		}
		err = writer.Flush()
		fmt.Println("send download")
		if err != nil {
			log.Printf("can't write to localhost:9999 %v", err)
			fmt.Println("Не удалось отправить запрос серверу")
			return
		}

		_, err = writer.Write([]byte(os.Args[2]+"\n"))
		if err != nil {
			log.Printf("can't write to localhost:9999 %v", err)
			fmt.Println("Не удалось отправить запрос серверу")
			return
		}
		err = writer.Flush()

		if err != nil {
			log.Printf("can't write to localhost:9999 %v", err)
			fmt.Println("Не удалось отправить запрос серверу")
			return
		}
		fmt.Println("send file name")
		file, err := os.Create("client/"+os.Args[2])
		if err != nil {
			log.Printf("can't create file %v", err)
			fmt.Println("Не удалось создать файл")
			return
		}
		defer func() {
			err = file.Close()
			if err != nil {
				log.Printf("can't close file %v", err)
			}
		}()
		fmt.Println("start reading")
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Printf("can't write to file %v", err)
			fmt.Println("Не удалось записать в файл")
			return
		}

		_, err = file.Write(bytes)
		if err != nil {
			log.Printf("can't write to file %v", err)
			fmt.Println("Не удалось записать в файл")
			return
		}
	}else if *upload{
		_, err = writer.Write([]byte("UPLOAD\n"))
		if err != nil {
			log.Printf("can't write to localhost:9999 %v", err)
			fmt.Println("Не удалось отправить запрос серверу")
			return
		}
		err = writer.Flush()
		if err != nil {
			log.Printf("can't write to localhost:9999 %v", err)
			fmt.Println("Не удалось отправить запрос серверу")
			return
		}

		_, err = writer.Write([]byte(os.Args[2]+"\n"))
		if err != nil {
			log.Printf("can't write to localhost:9999 %v", err)
			fmt.Println("Не удалось отправить запрос серверу")
			return
		}
		err = writer.Flush()

		if err != nil {
			log.Printf("can't write to localhost:9999 %v", err)
			fmt.Println("Не удалось отправить запрос серверу")
			return
		}
		file, err := os.Open("client/" + os.Args[2])
		if err != nil {
			log.Printf("can't open file %v", err)
			fmt.Println("Не удалось открыть файл")
			return
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("can't read file %v", err)
			fmt.Println("Не удалось прочитать файл")
			return
		}

		_, err = writer.Write(bytes)
		if err != nil {
			log.Printf("can't write to localhost %v", err)
			fmt.Println("Не удалось отправить файл")
			return
		}
		err = writer.Flush()
		if err != nil {
			log.Printf("can't write to localhost %v", err)
			fmt.Println("Не удалось отправить файл")
			return
		}
	}

}

