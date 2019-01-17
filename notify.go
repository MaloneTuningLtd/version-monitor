package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SlackMessage struct {
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func slackFromVersions(n, o Version) SlackMessage {
	attachMsg := fmt.Sprintf("Updated to %s", n.Version)
	if o.Version != "" && o.Version != n.Version {
		attachMsg = attachMsg + fmt.Sprintf(" from %s", o.Version)
	}

	attachment := SlackAttachment{
		Title: n.Name,
		Text:  attachMsg,
	}

	message := SlackMessage{
		Text:        fmt.Sprintf("%s has been recently updated!", n.Name),
		Attachments: []SlackAttachment{attachment},
	}

	return message
}

func slackMessage(recent *Version) {
	old, err := vers.Get(recent.Name)
	if err != nil {
		old = &Version{}
	}

	slackMsg := slackFromVersions(*recent, *old)
	buf := bytes.NewBuffer(nil)

	err = json.NewEncoder(buf).Encode(&slackMsg)
	if err != nil {
		log.Println(err)
		return
	}

	http.Post(config.SlackHook, "application/json", buf)
}

func notify(ver *Version) {
	fmt.Printf("REPOSITORY: %s updated to %s\n", ver.Name, ver.Version)
}
