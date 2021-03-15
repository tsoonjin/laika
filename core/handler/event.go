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

// approvalRequest mocks the simple "Approval" template located on block kit builder website
func exampleOne() slack.MsgOption {

	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "You have a new request:\n*<fakeLink.toEmployeeProfile.com|Fred Enriquez - New device request>*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Fields
	typeField := slack.NewTextBlockObject("mrkdwn", "*Type:*\nComputer (laptop)", false, false)
	whenField := slack.NewTextBlockObject("mrkdwn", "*When:*\nSubmitted Aut 10", false, false)
	lastUpdateField := slack.NewTextBlockObject("mrkdwn", "*Last Update:*\nMar 10, 2015 (3 years, 5 months)", false, false)
	reasonField := slack.NewTextBlockObject("mrkdwn", "*Reason:*\nAll vowel keys aren't working.", false, false)
	specsField := slack.NewTextBlockObject("mrkdwn", "*Specs:*\n\"Cheetah Pro 15\" - Fast, really fast\"", false, false)

	fieldSlice := make([]*slack.TextBlockObject, 0)
	fieldSlice = append(fieldSlice, typeField)
	fieldSlice = append(fieldSlice, whenField)
	fieldSlice = append(fieldSlice, lastUpdateField)
	fieldSlice = append(fieldSlice, reasonField)
	fieldSlice = append(fieldSlice, specsField)
    inputBlock := slack.NewInputBlock(
        "url_block",
        slack.NewTextBlockObject("plain_text", "URL", false, false),
        slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "https://www.google.com", false, false), "input_url"),
    )

	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)

	// Approve and Deny Buttons
	approveBtnTxt := slack.NewTextBlockObject("plain_text", "Approve", false, false)
	approveBtn := slack.NewButtonBlockElement("", "click_me_123", approveBtnTxt)

	denyBtnTxt := slack.NewTextBlockObject("plain_text", "Deny", false, false)
	denyBtn := slack.NewButtonBlockElement("", "click_me_123", denyBtnTxt)

	actionBlock := slack.NewActionBlock("", approveBtn, denyBtn)


	// b, err := json.MarshalIndent(msg, "", "    ")
	// if err != nil {
	// 	return nil, err
	// }
    return slack.MsgOptionBlocks(headerSection, fieldsSection, inputBlock, actionBlock)
}

func EventHandler(w http.ResponseWriter, r *http.Request) {

    config := config.LoadConfig()
    var matchingKeyword string = "clubhouse"
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
