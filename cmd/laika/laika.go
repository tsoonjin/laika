package main

import (
    "fmt"
    "net/http"
    "github.com/tsoonjin/laika/core"
)

func slackChat(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Are you ready to chat !!!")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Are you ready to rumble ?")
}


func setupRoutes() {
  http.HandleFunc("/slack", slackChat)
  http.HandleFunc("/", rootHandler)
}

func main() {
    config := laika.LoadConfig()
    fmt.Println(fmt.Sprintf("Running %s server at port %s", config.Env, config.Port))
    setupRoutes()
    http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)
}
