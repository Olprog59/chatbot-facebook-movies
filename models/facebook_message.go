package models

import (
	"bytes"
	"chatbot-facebook-movies/commons"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type FacebookMessage struct {
	Object string `json:"object"`
	Entry  []struct {
		Id         string `json:"id"`
		Time       int64  `json:"time"`
		TimeString string
		Messaging  []struct {
			Sender struct {
				Id string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				Id string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   struct {
				Mid  string `json:"mid"`
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

type FacebookMessageWelcome struct {
	Greeting []struct {
		Locale string `json:"locale"`
		Text   string `json:"text"`
	} `json:"greeting"`
}

type FacebookSendMessage struct {
	Recipient struct {
		Id string `json:"id"`
	} `json:"recipient"`
	MessagingType string `json:"messaging_type"`
	Message       struct {
		Text         string                             `json:"text"`
		QuickReplies *[]FacebookSendMessageQuickReplies `json:"quick_replies"`
	} `json:"message"`
}

type FacebookSendMessageQuickReplies struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Payload     string `json:"payload"`
	ImageUrl    string `json:"image_url"`
}

func (f *FacebookMessage) GetSenderId() string {
	return f.Entry[0].Messaging[0].Sender.Id
}

func (f *FacebookMessage) GetEntryMessages() (message, senderId string) {
	for _, entry := range f.Entry {
		for _, messaging := range entry.Messaging {
			return messaging.Message.Text, messaging.Sender.Id
		}
	}

	return "", ""
}

func (f *FacebookSendMessage) SendMessage() error {
	// validate empty message
	if len(f.Message.Text) == 0 {
		return errors.New("message can't be empty")
	}

	// marshal request data
	data, err := json.Marshal(f)
	if err != nil {
		return fmt.Errorf("error marshall request: %w", err)
	}

	// setup http request
	url := fmt.Sprintf("%s/%s?access_token=%s", commons.GRAPHQL_URL, "me/messages", commons.FB_TOKEN)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed wrap request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// send http request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed send request: %w", err)
	}
	defer res.Body.Close()

	return nil
}

func SendWelcome() error {
	// validate empty message

	var msg = new(FacebookMessageWelcome)
	msg.Greeting = make([]struct {
		Locale string `json:"locale"`
		Text   string `json:"text"`
	}, 1)
	msg.Greeting[0].Locale = "default"
	msg.Greeting[0].Text = "Bonjour {{user_first_name}} ! Je suis un bot qui peut vous aider à trouver des films et des séries. Pour commencer, tapez le mot clé 'bonjour'."

	// marshal request data
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshall request: %w", err)
	}

	// setup http request
	url := fmt.Sprintf("%s/%s?access_token=%s", commons.GRAPHQL_URL, "me/messager_profile", commons.FB_TOKEN)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed wrap request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// send http request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed send request: %w", err)
	}
	defer res.Body.Close()

	return nil
}
