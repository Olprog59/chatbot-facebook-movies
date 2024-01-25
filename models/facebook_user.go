package models

import (
	"chatbot-facebook-movies/commons"
	"encoding/json"
	"fmt"
	"net/http"
)

type FaceBookUser struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Id         string
	Authorized bool
}

func (f *FaceBookUser) CheckIfUserAuthorized(senderId string) (string, bool) {
	user, err := f.GetFacebookUser(senderId)
	if err != nil {
		return "", false
	}

	users := GetAuthorizedUsers()

	for _, u := range users {
		if u.FirstName == user.FirstName && u.LastName == user.LastName {
			return u.FirstName + " " + u.LastName, true
		}
	}
	return "", false
}

func (f *FaceBookUser) GetFacebookUser(senderId string) (*FaceBookUser, error) {

	// setup http request
	url := fmt.Sprintf("%s/%s?fields=first_name,last_name&access_token=%s", commons.GRAPHQL_URL, senderId, commons.FB_TOKEN)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return f, fmt.Errorf("failed wrap request: %w", err)
	}

	// send http request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return f, fmt.Errorf("failed send request: %w", err)
	}
	defer res.Body.Close()

	// unmarshal response body
	if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
		return f, fmt.Errorf("failed unmarshal response body: %w", err)
	}
	f.Id = senderId
	f.Authorized = true

	return f, nil
}
