package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"
)

var fromHost = flag.String("from", "", "From host:port")
var toHost = flag.String("to", "", "To host:port")
var web = flag.String("web", "", "Web server host:port")

func main() {
	flag.Parse()
	if *fromHost == "" || *toHost == "" || *web == "" {
		println("secret-proxy is a small web server that starts connection forwarding on-demand")
		println("\nUsage:\n")
		flag.PrintDefaults()
		return
	}
	http.Handle("/", http.FileServer(http.Dir("html")))
	http.HandleFunc("/action", action)

	log.Printf("Web server listening on %s", *web)
	log.Fatal(http.ListenAndServe(*web, nil))
}

func action(w http.ResponseWriter, req * http.Request) {
	actionId := req.URL.Query().Get("id")
	if actionId == "startForwarding" {
		go runListener()
		w.Write([]byte("OK"))
	} else {
		w.Write([]byte("Error"))
	}
}

func runListener() {
	// listen
	listen, err := net.Listen("tcp", *fromHost)
	if checkErr(err) {
		return
	}
	defer listen.Close()
	if l,ok := listen.(*net.TCPListener); ok {
		l.SetDeadline(time.Now().Add(5 * time.Second))
	}
	log.Printf("Listening on %s", *fromHost)

	// accept and close listener
	accept, err := listen.Accept()
	if checkErr(err) {
		return
	}
	// allow only one connection in
	listen.Close()
	log.Printf("Accepted %s, Listener stopped", accept.RemoteAddr())
	defer accept.Close()

	// client conn
	dial, err := net.Dial("tcp", *toHost)
	defer dial.Close()
	if checkErr(err) {
		return
	}

	// forward
	done := make(chan bool)
	log.Printf("Forwarding to %s", *toHost)
	go copyData(accept, dial, done)
	go copyData(dial, accept, done)
	<- done
	<- done
	log.Print("Connection closed")
}

func copyData(from net.Conn, to net.Conn, done chan<- bool) {
	defer func(){
		done <- true
	}()
	buf := make([]byte, 65536)
	for {
		read, err := from.Read(buf)
		if checkErr(err) || read <=0 {
			from.Close()
			to.Close()
			return
		}
		write, err := to.Write(buf[:read])
		if checkErr(err) || write <=0 {
			from.Close()
			to.Close()
			return
		}
	}
}

func checkErr(err error) bool {
	if err != nil {
		log.Printf("%T: %s", err, err)
		return true
	}
	return false
}