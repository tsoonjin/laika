package lai

import (
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
