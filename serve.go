package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var additionalOutput = false

var helpText = `Usage: serve [OPTION]... [PATH]...

Serve path and all its files via HTTP.
Provides directory listing when the served path is accessed with browser.

If PATH parameter is omitted, the current directory is served by default.

OPTION(s):
`

func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r)

		uri := r.URL.String()
		remoteAddr := r.RemoteAddr
		userAgent := r.Header.Get("User-Agent")
		method := r.Method

		log.Println(fmt.Sprintf("%s \"%s %s\" \"%s\"", remoteAddr, method, uri, userAgent))
	}

	return http.HandlerFunc(fn)
}

func determineIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine IP address: %s", err.Error())
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func determinePublicIP() string {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine IP address, connection error: %s", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine IP address, request error: %s", err.Error())
	}

	return string(body)
}

func main() {
	var port int
	var logRequests bool

	flag.IntVar(&port, "p", 80, "Port number to use")
	timeout := flag.Duration("t", 0, "Timeout in seconds after which stop execution (for example 2h30m)")
	flag.BoolVar(&additionalOutput, "v", false, "Enable additional output")
	flag.BoolVar(&logRequests, "l", false, "Log requests to the server")
	var urlOutput string
	flag.StringVar(&urlOutput, "o", "", "Output URL using method: [hostname, dns, public]")

	flag.Usage = func() {
		fmt.Println(helpText)
		flag.PrintDefaults()
	}

	flag.Parse()

	path := flag.Arg(0)
	if len(path) == 0 {
		path, _ = os.Getwd()
	}
	var fileSystem http.Dir
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		fileSystem = http.Dir(path)
	} else {
		fileSystem = http.Dir(filepath.Dir(path))
	}
	address := fmt.Sprintf(":%d", port)
	if port == 80 {
		address = ""
	}

	if additionalOutput {
		log.Println(fmt.Sprintf("Using port: %d", port))
		if timeout.Seconds() == 0 {
			log.Println("No timeout, server will run until stopped")
		} else {
			log.Println(fmt.Sprintf("Timeout set to: %s", timeout))
		}
		log.Println("Serving path:", fileSystem)
		if len(flag.Args()) > 1 {
			log.Println("Additional arguments were given but unused", flag.Args())
		}
		if logRequests {
			log.Println("Request logging is enabled.")
		} else {
			log.Println("Request logging is disabled.")
		}
	}

	switch urlOutput {
	case "":
		break
	case "hostname":
		hostname, _ := os.Hostname()
		fmt.Fprintln(os.Stdout, "Your files are now reachable at:")
		fmt.Fprintln(os.Stdout, fmt.Sprintf("http://%s%s/ ", hostname, address))
		break
	case "dns":
		fmt.Fprintln(os.Stdout, "Your files are now reachable at:")
		fmt.Fprintln(os.Stdout, fmt.Sprintf("http://%s%s/ ", determineIP(), address))
		break
	case "public":
		fmt.Fprintln(os.Stdout, "Your files are now reachable at:")
		fmt.Fprintln(os.Stdout, fmt.Sprintf("http://%s%s/ ", determinePublicIP(), address))
		break
	default:
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Unknown URL lookup method %s, choose one of the following: [hostname, dns, public]", urlOutput))
		break
	}

	handler := http.FileServer(fileSystem)
	if logRequests {
		// Wrap handler to the logger
		handler = logRequestHandler(handler)
	}

	if timeout.Seconds() == 0 {
		// Serve until break signal
		err := http.ListenAndServe(address, handler)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %s", err.Error())
		}
	} else {
		// Serve with timeout
		server := http.Server{Addr: address, Handler: handler}
		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Fprintf(os.Stderr, "Server error: %s", err.Error())
			}
		}()
		select {
		case <-ctx.Done():
			// Shutdown the server when the context is canceled
			server.Shutdown(ctx)
		}
		if additionalOutput {
			log.Println(fmt.Sprintf("Server shutdown at %s", time.Now().Local()))
		}
	}
}
