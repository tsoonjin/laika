package main

import (
    "fmt"
    "net/http"
    "github.com/joho/godotenv"
)

type Config struct {
    Env string
    Port string
}

func loadConfig() Config {
    if os.Getenv("ENV") != "production" {
        if err := godotenv.Load(); err != nil {
            log.Println("Failed to load .env file")
        }
    }
    config := Config {
        Env: "develop",
        Port: "3000",
    }
    env := os.Getenv("ENV")
    if env != "" {
        config.Env = env
    }
    port := os.Getenv("PORT")
    if port != "" {
        config.Port = port
    }
    return config
}
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
    fmt.Println("Welcome Master !!! ")
    setupRoutes()
    http.ListenAndServe(":3000", nil)
}
