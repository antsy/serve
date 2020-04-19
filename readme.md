# (Simple HTTP) Serve

## Background

If you want to use Python's `SimpleHTTPServer` feature but don't want to install Python, then this program might be what you need.

It will serve any locally available folder via HTTP using Go's builtin [FileServer](https://golang.org/pkg/net/http/#FileServer).

The project only uses Go's standard library, there are no additional dependencies.


## Usage

As a prerequisite, you need to have [Go](https://golang.org/doc/install) installed.

You can execute the program without producing a permanent binary by running:

`go run serve.go`

By default, you will now have your current directory served at [https://localhost/](https://localhost/). You can give a path as a parameter to the program to serve any other directory.

Press [Ctrl+C](https://en.wikipedia.org/wiki/Control-C) to end execution or set timeout with `-t` -parameter.


### Parameters

Additional command line parameters for the server.

| Parameter | What it does  |
|-----------|---------------|
| -h        | Display help. Basically the same thing you are reading right now.
| -p        | Sets a custom port for Serve to use. Defaults to 80 (HTTP).
| -t        | Timeout for execution.<br><br>If you want to leave server on only for a specific period of time, this setting will end the execution after the given duration.<br><br>A duration string is a sequence of decimal numbers, each with optional fraction and a unit suffix, such as "120s", "1.5h" or "2h45m".
| -l        | Log requests. Will log the requester IP address and what was requested. <br><br>As of now, the logging is to stdout only.
| -o        | This method can be used to output HTTP URL where your files are being served.<br>Determining this is always bit of a guess work if you want to share the files somewhere outside your own computer.<br><br>Valid options are as follows:<br><br>*hostname* - Use OS Hostname setting to determine the URL.<br>*dns*      - Use Google's DNS server (8.8.8.8) to determine your IP address.<br>*public*   - Use ipify.org to determine your IP address.<br><br>All, some or none of these might work for your use case.
| -v        | Verbose output. Will log additional information about what is happening in the program.
 

## Build

As a prerequisite, you need to have [Go](https://golang.org/doc/install) installed.

If you wish to compile stand alone version of the program, simply run `go build` to compile executable of your current platform.

See Go's own documentation for additional information about compilation in general and cross compilation to different platforms (Windows/Linux/macOS/etc.)
