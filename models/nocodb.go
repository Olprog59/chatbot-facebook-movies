package models

import (
	"bytes"
	"chatbot-facebook-movies/commons"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func SendRequestNocoDB(movie Movie) error {
	log.Println("sendRequestNocoDB")

	type Body struct {
		Title    string `json:"Title"`
		Year     int    `json:"Year"`
		Type     string `json:"Type"`
		IdUser   string `json:"IdUser"`
		NameUser string `json:"NameUser"`
	}

	var body = new(Body)

	body.Title = movie.Title
	body.Year = movie.Year
	body.Type = movie.Type
	body.IdUser = movie.IdUser
	body.NameUser = movie.NameUser

	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed marshal body: %w", err)
	}
	log.Println(string(b))
	// setup http request
	url := fmt.Sprintf("%s/%s/records", commons.NOCO_URL, commons.NOCO_TABLE_ID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("xc-auth", commons.NOCO_API_KEY)
	if err != nil {
		return fmt.Errorf("failed wrap request: %w", err)
	}

	// send http request
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed send request: %w", err)
	}

	// Lecture du corps de la réponse
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Affichage du statut de la réponse et du corps
	log.Printf("Response status: %s", res.Status)
	log.Printf("Response body: %s", string(responseBody))

	return nil
}
