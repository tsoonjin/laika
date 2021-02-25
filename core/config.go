package laika

import (
    "github.com/joho/godotenv"
    "os"
    "log"
)

type Config struct {
    Env string
    Port string
    Token string
    Secret string
}

func LoadConfig() Config {
    if os.Getenv("ENV") != "production" {
        if err := godotenv.Load(); err != nil {
            log.Println("Failed to load .env file")
        }
    }
    config := Config {
        Env: "develop",
        Port: "4000",
        Token: "",
    }
    env := os.Getenv("ENV")
    if env != "" {
        config.Env = env
    }
    port := os.Getenv("PORT")
    if port != "" {
        config.Port = port
    }
    token := os.Getenv("BOT_TOKEN")
    if token != "" {
        config.Token = token
    }
    secret := os.Getenv("SECRET")
    if secret != "" {
        config.Secret = secret
    }
    return config
}
