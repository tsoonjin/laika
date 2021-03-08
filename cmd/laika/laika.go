package main

import (
    "os"
    "fmt"
    "net/http"
    "io/ioutil"
    "io"
    "encoding/json"
    "github.com/slack-go/slack"
    "github.com/tsoonjin/laika/core"
    "github.com/tsoonjin/laika/core/scripts/google"
    "github.com/tsoonjin/laika/core/handler"
)


func slackChat(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Are you ready to chat !!!")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Are you ready to rumble ?")
}


func setupRoutes(config laika.Config, api *slack.Client) {
    googleClient := google.Client{
        Config: google.Config{
            ApiKey: os.Getenv("PAGESPEED_INSIGHT_API_KEY"),
        },
    }
    http.HandleFunc("/pagespeed", googleClient.Handler)
    http.HandleFunc("/slack/events", slackBot.EventHandler)
    http.HandleFunc("/slack/interactive", slackBot.InteractionHandler)
    http.HandleFunc("/slack/command", func(w http.ResponseWriter, r *http.Request) {

		verifier, err := slack.NewSecretsVerifier(r.Header, config.Secret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = verifier.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch s.Command {
		case "/lai":
			params := &slack.Msg{
                Text: s.Text,
                ResponseType: "in_channel",
            }
			b, err := json.Marshal(params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func main() {
    config := laika.LoadConfig()
    var api = slack.New(config.Token)
    fmt.Println(fmt.Sprintf("Running %s server at port %s", config.Env, config.Port))
    setupRoutes(config, api)
    http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)
}
