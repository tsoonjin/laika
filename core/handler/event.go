package slackBot

import (
	"encoding/json"
	"io/ioutil"
    "strings"
    "log"
	"net/http"
	"github.com/slack-go/slack"
    "github.com/slack-go/slack/slackevents"
    config "github.com/tsoonjin/laika/core"
)

func EventHandler(w http.ResponseWriter, r *http.Request) {

    config := config.LoadConfig()
    var matchingKeyword string = "lighthouse"
    var api = slack.New(config.Token)
    var signingSecret = config.Secret
    // Read request body
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    if _, err := sv.Write(body); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if err := sv.Ensure(); err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
    log.Println(eventsAPIEvent.Type)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if eventsAPIEvent.Type == slackevents.URLVerification {
        var r *slackevents.ChallengeResponse
        err := json.Unmarshal([]byte(body), &r)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "text")
        w.Write([]byte(r.Challenge))
    }
    if eventsAPIEvent.Type == slackevents.CallbackEvent {
        innerEvent := eventsAPIEvent.InnerEvent
        switch ev := innerEvent.Data.(type) {
        case *slackevents.AppMentionEvent:
            api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
        case *slackevents.MessageEvent:
            if strings.Contains(strings.ToLower(ev.Text), matchingKeyword) {
                log.Println("Lighthouse script activated")
                attachment := slack.Attachment{
                    Text:       "Please list the urls to be scanned",
                    Color:      "#f9a41b",
                    CallbackID: "lighthouse-scan",
                    Actions: []slack.AttachmentAction{
                        slack.AttachmentAction{
                            Name:  "submit",
                            Text:  "Submit",
                            Type:  "button",
                            Value: "submit",
                        },
                        slack.AttachmentAction{
                            Name:  "cancel",
                            Text:  "Cancel",
                            Type:  "button",
                            Value: "cancel",
                        },
                    },
                }
                api.PostMessage(ev.Channel, slack.MsgOptionText("Lighthouse Scan", false), slack.MsgOptionAttachments(attachment))
            }
        }
    }
}
