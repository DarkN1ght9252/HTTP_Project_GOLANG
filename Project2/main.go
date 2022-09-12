package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	port       string
	bytesArray []byte
)

const (
	error404 string = "./www/404.html"
)

func main() {
	//Runs on port 8080
	port = "8080"

	//tries to connect to a Client
	listenTCP, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	//Defer Closes listen
	defer listenTCP.Close()
	//Runs concurrent connections
	for {
		//Listens for requests
		conn, err := listenTCP.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		//Creates a thread to deliver packets
		go connectionHandler(conn)
	}
}

func connectionHandler(conn net.Conn) {
	//Reads First Line from Request
	httpData, err := bufio.NewReader(conn).ReadString('\n')
	log.Println(httpData)

	if err != nil {
		fmt.Println(err)
		return
	}
	//Sends Request info to Build Packet either for HEAD or GET requests
	if strings.Contains(httpData, "HEAD") {
		bytesArray = []byte(handleHEAD(httpData, conn))
		log.Println("Sent HEAD")
	} else if strings.Contains(httpData, "GET") {
		bytesArray = []byte(handleGET(httpData, conn))
		// log.Println(string(bytesArray))
		log.Println("Sent Get")
	}

	//Sends the http/1.1 response Packet to Client
	conn.Write(bytesArray)
	//Resets the Byte array to clean it out
	bytesArray = make([]byte, 0)
	//Closes connection
	conn.Close()
}

//ex: of Response
// HTTP/1.1 200 OK
// Server: cihttpd
// Last-Modified: Tue, 01 Dec 2009 20:18:22 GMT
// Content-Length: 29769

//Else

// HTTP/1.1 404 /notfound.html
// Server: cihttpd
// Last-Modified: Tue, 01 Dec 2009 20:18:22 GMT
// Content-Length: 29769

func handleHEAD(httpData string, conn net.Conn) (response string) {
	log.Println("HEAD Received")

	//Split First Line from Client
	s := strings.Fields(httpData)
	fmt.Println(s)
	// Output: ['GET', '/', 'HTTP/1.1']

	//If client asks for base Site it then returns the index.html Page
	if s[1] == "/" {
		s[1] = "/index.html"
	}

	reader, err := os.Open("./www" + s[1])

	//If file not Found Sends 404 Page
	if err != nil {
		log.Println("Page not found Sending 404")
		return Send404(httpData, conn)
	}

	//Else Tries to grab data from file
	fileStat, err := os.Stat(reader.Name())

	if err != nil {
		log.Println("No Stats Found")
	}

	ModTime := fileStat.ModTime()
	log.Println(ModTime)

	//Grabs 200 OK Data puts it into Data
	buffer := bytes.NewBuffer(bytesArray)

	//First Line Header 200
	buffer.WriteString("HTTP/1.1 200 OK\r\n")
	buffer.WriteString("Server: cihttpd\r\n")
	buffer.WriteString("Last-Modified: " + fileStat.ModTime().Format(time.RFC1123) + "\r\n")
	buffer.WriteString("Content-Length: " + strconv.Itoa(int(fileStat.Size())) + "\r\n")
	buffer.WriteString("\r\n")

	buffer.ReadFrom(reader)

	return buffer.String()
}

func handleGET(httpData string, conn net.Conn) (response string) {
	log.Println("GET Received")

	//Split First Line from Client
	s := strings.Fields(httpData)
	fmt.Println(s)
	// Output: ['GET', '/', 'HTTP/1.1']
	log.Println(s[1])

	if s[1] == "/" {
		s[1] = "/index.html"
	}

	reader, err := os.Open("./www" + s[1])
	//If file not Found Sends 404 Page
	if err != nil {
		log.Println("Page not found Sending 404")
		return Send404(httpData, conn)
	}

	//Else Tries to grab data from file
	fileStat, err := os.Stat(reader.Name())

	if err != nil {
		log.Println("No Stats Found")
	}

	ModTime := fileStat.ModTime()
	log.Println(ModTime)

	//Grabs 200 OK Data puts it into Data
	buffer := bytes.NewBuffer(bytesArray)

	//First Line Header 200
	buffer.WriteString("HTTP/1.1 200 OK\r\n")
	buffer.WriteString("Server: cihttpd\r\n")
	buffer.WriteString("Last-Modified: " + fileStat.ModTime().Format(time.RFC1123) + "\r\n")
	buffer.WriteString("Content-Length: " + strconv.Itoa(int(fileStat.Size())) + "\r\n")
	buffer.WriteString("\r\n")

	buffer.ReadFrom(reader)

	return buffer.String()
}

func Send404(httpData string, conn net.Conn) (response string) {
	reader, err := os.Open(error404)

	if err != nil {
		log.Println("404 Not Found LOL")
	}

	//Else Tries to grab data from file
	fileStat, err := os.Stat(reader.Name())

	if err != nil {
		log.Println("No Stats Found")
	}

	//Grabs 404.html Data puts it into Data
	buffer := bytes.NewBuffer(bytesArray)

	buffer.WriteString("HTTP/1.1 404 Not Found\r\n")
	buffer.WriteString("Server: cihttpd\r\n")
	buffer.WriteString("Last-Modified: " + fileStat.ModTime().Format(time.RFC1123) + "\r\n")
	buffer.WriteString("Content-Length: " + strconv.Itoa(int(fileStat.Size())) + "\r\n")
	buffer.WriteString("\r\n")

	buffer.ReadFrom(reader)

	return buffer.String()
}
