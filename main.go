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
	mid := logMiddleware(mux)
	mux.HandleFunc("GET /chatbot/{$}", webhookHandler)
	mux.HandleFunc("POST /chatbot/{$}", webhookPostHandler)
	mux.HandleFunc("GET /privacy-policy", privacyPolicyHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mid,
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

func logMiddleware(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		mux.ServeHTTP(w, r)
	})
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
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            padding: 20px;
            max-width: 700px;
            margin: auto;
        }
        h1 {
            color: #333;
        }
        h2 {
            font-size: 1.2em;
        }
    </style>
</head>
<body>
    <h1>Politique de Confidentialité de Chatbot Facebook Movies</h1>
    
    <h2>1. Collecte de Données</h2>
    <p>Notre application Chatbot Facebook Movies collecte les informations suivantes :</p>
    <ul>
        <li>Nom et prénom pour vérifier si l'utilisateur est autorisé à utiliser le chatbot.</li>
        <li>Préférences de l'utilisateur concernant les films et les séries (type, titre, année) lors de l'utilisation du chatbot.</li>
    </ul>
    
    <h2>2. Utilisation des Données</h2>
    <p>Les données collectées sont utilisées aux fins suivantes :</p>
    <ul>
        <li>Pour permettre une utilisation personnalisée et autorisée du chatbot.</li>
        <li>Pour des analyses générales, dans le but d'améliorer l'expérience utilisateur.</li>
    </ul>
    
    <h2>3. Partage des Données</h2>
    <p>Nous ne partageons pas les informations collectées avec des tiers. Toutes les données restent privées et sont exclusivement utilisées pour les fonctionnalités internes de l'application.</p>
    
    <h2>4. Sécurité des Données</h2>
    <p>La sécurité de vos données est importante pour nous. Bien que nous n'utilisions pas de mesures de sécurité spécifiques au-delà de la connexion HTTPS pour accéder à NocoDB, nous nous engageons à protéger la confidentialité de vos informations.</p>
    
    <h2>5. Droits des Utilisateurs</h2>
    <p>Les utilisateurs peuvent demander l'accès à leurs données personnelles ou leur suppression en nous contactant directement. Comme cette application est destinée à un usage privé et limité à un cercle restreint, toute requête de ce type sera traitée manuellement sur demande.</p>
    
    <h2>6. Modifications de la Politique de Confidentialité</h2>
    <p>Cette politique de confidentialité peut être mise à jour occasionnellement. Comme l'application est destinée à un usage restreint et privé, les notifications de mise à jour ne sont pas prévues. Il incombe aux utilisateurs de consulter régulièrement cette politique pour se tenir informés de nos pratiques.</p>

</body>
</html>
`
	_, err := fmt.Fprintf(w, privacyPolicyHTML)
	if err != nil {
		return
	}
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

	log.Println("Webhook verified")

	return
}
