package main

import (
  "webmachine"
  "http"
  "io"
  "log"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
  wm := webmachine.NewWebMachine()
  wm.AddRouteHandler(webmachine.NewFileResource("/Users/pomack/", "/", false))
	err := http.ListenAndServe(":12345", wm)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.String())
	}
}


