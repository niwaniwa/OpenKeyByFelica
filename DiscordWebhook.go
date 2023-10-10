package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

var (
	webhookURL              string = ""
	userName                string = "鍵管理bot"
	avatarURL               string = ""
	openMessage             string = ""
	closeMessage            string = ""
	descriptionMessage      string = ""
	openDescriptionMessage         = ""
	closeDescriptionMessage        = ""
)

type DiscordWebhookBody struct {
	UserName  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Content   string  `json:"content"`
	Embeds    []Embed `json:"embeds"`
}

type Embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Color       int     `json:"color,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func postDiscord(content DiscordWebhookBody) {
	payload, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling the webhook body:", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Error sending the webhook request:", err)
		return
	}
	defer resp.Body.Close()
}
