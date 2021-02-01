package main

import (
    "fmt"
    "net/http"
)

func slackChat(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Are you ready to chat !!!")
}


func setupRoutes() {
  http.HandleFunc("/", slackChat)
}

func main() {
    fmt.Println("Welcome Master !!! ")
    setupRoutes()
    http.ListenAndServe(":3000", nil)
}
