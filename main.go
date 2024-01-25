package main

import (
	"chatbot-facebook-movies/commons"
	"chatbot-facebook-movies/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Ldate | log.Ltime)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", webhookHandler)
	mux.HandleFunc("POST /{$}", webhookPostHandler)
	mux.HandleFunc("GET /privacy-policy", privacyPolicyHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	go func() {
		err := models.SendWelcome()
		if err != nil {
			return
		}
	}()

	log.Println("Starting server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}

func webhookPostHandler(w http.ResponseWriter, r *http.Request) {
	var fbMessage models.FacebookMessage
	err := json.NewDecoder(r.Body).Decode(&fbMessage)
	if err != nil {
		return
	}
	messages, senderId := fbMessage.GetEntryMessages()
	if messages == "" {
		return
	}
	parcours(messages, senderId)
	return
}

func privacyPolicyHandler(w http.ResponseWriter, r *http.Request) {
	const privacyPolicyHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Politique de Confidentialité</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { width: 80%; margin: auto; overflow: hidden; }
        h1, p { padding: 15px 0; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Politique de Confidentialité</h1>
        <p>
            Cette page présente la politique de confidentialité de notre application.
        </p>
        <p>
            <!-- Ajoute ici le contenu détaillé de ta politique de confidentialité -->
        </p>
    </div>
</body>
</html>
`
	fmt.Fprintf(w, privacyPolicyHTML)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	// read token from query parameter
	verifyToken := r.URL.Query().Get("hub.verify_token")

	// verify the token included in the incoming request
	if verifyToken != commons.WEBHOOK_TOKEN {
		log.Printf("invalid verification token: %s", verifyToken)
		return
	}

	// write string from challenge query parameter
	if _, err := w.Write([]byte(r.URL.Query().Get("hub.challenge"))); err != nil {
		log.Printf("failed to write response body: %v", err)
		return
	}

	return
}
