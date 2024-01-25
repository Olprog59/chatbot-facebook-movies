package main

import (
	"chatbot-facebook-movies/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	PAYLOAD_HELLO  = "bonjour"
	PAYLOAD_FILMS  = "films"
	PAYLOAD_SERIES = "series"
	PAYLOAD_YEAR   = "year"
)

func parcours(message string, senderId string) string {
	name, ok := checkUser(senderId)
	if !ok {
		sendError(senderId)
		return ""
	}

	movie := models.GetMovies(senderId)
	movie.NameUser = name
	movie.UpdateMovies(senderId)

	var result string
	// TODO: Dire Bonjour
	switch strings.ToLower(message) {
	case PAYLOAD_HELLO:
		sendHello(senderId)
	case PAYLOAD_FILMS:
		getType(senderId, PAYLOAD_FILMS)
		movie.State = PAYLOAD_FILMS
	case PAYLOAD_SERIES:
		getType(senderId, PAYLOAD_SERIES)
		movie.State = PAYLOAD_SERIES
	default:
		switch movie.State {
		case PAYLOAD_FILMS:
			movie.Type = "Films"
			movie.Title = message
			movie.State = PAYLOAD_YEAR
			getYear(senderId)
		case PAYLOAD_SERIES:
			movie.Type = "Series"
			movie.Title = message
			movie.State = PAYLOAD_YEAR
			getYear(senderId)
		case PAYLOAD_YEAR:
			var err error
			if regexp.MustCompile(`\d{4}`).MatchString(message) {
				movie.Year, err = strconv.Atoi(message)
				if err != nil {
					log.Println(err)
					sendError(senderId)
					return ""
				}
			}
			movie.State = PAYLOAD_HELLO
			// TODO: envoyer à l'API nocodb
			err = models.SendRequestNocoDB(movie)
			if err != nil {
				sendError(senderId)
				log.Println(err)
				return ""
			}
			log.Println("The end !")
			sendMessageEnding(senderId)
		}
	}
	movie.UpdateMovies(senderId)
	log.Printf("%+v\n", movie)
	return result
}

func getYear(senderId string) {
	msg := new(models.FacebookSendMessage)
	msg.Recipient.Id = senderId
	msg.Message.Text = "Entres l'année"
	err := msg.SendMessage()
	if err != nil {
		log.Println(err)
		return
	}
}
func sendMessageEnding(senderId string) {
	msg := new(models.FacebookSendMessage)
	msg.Recipient.Id = senderId
	msg.Message.Text = "C'est bon pour moi ! A bientôt !"
	err := msg.SendMessage()
	if err != nil {
		log.Println(err)
		return
	}
}
func sendError(senderId string) {
	msg := new(models.FacebookSendMessage)
	msg.Recipient.Id = senderId
	msg.Message.Text = "⚠️ 🥵 Je n'ai pas compris ! Peux-tu recommencer ?"
	err := msg.SendMessage()
	if err != nil {
		log.Println(err)
		return
	}
}

func getType(senderId, typeMovie string) {
	msg := new(models.FacebookSendMessage)
	msg.Recipient.Id = senderId
	if typeMovie == PAYLOAD_FILMS {
		msg.Message.Text = "🎬 Entres le nom du film"
	} else if typeMovie == PAYLOAD_SERIES {
		msg.Message.Text = "📺 Entres le nom de la série"
	}
	err := msg.SendMessage()
	if err != nil {
		log.Println(err)
		return
	}
}

func sendHello(senderId string) {
	msg := new(models.FacebookSendMessage)
	msg.Recipient.Id = senderId
	msg.Message.Text = "Bonjour, dis-moi ce que tu veux voir ?"
	msg.Message.QuickReplies = &[]models.FacebookSendMessageQuickReplies{
		{
			ContentType: "text",
			Title:       cases.Title(language.French).String(PAYLOAD_FILMS),
		},
		{
			ContentType: "text",
			Title:       cases.Title(language.French).String(PAYLOAD_SERIES),
		},
	}
	err := msg.SendMessage()
	if err != nil {
		log.Println(err)
		sendError(senderId)
		return
	}
}

func checkUser(senderId string) (string, bool) {
	var fbUser = new(models.FaceBookUser)
	if user, err := fbUser.GetFacebookUser(senderId); err == nil {
		if name, ok := fbUser.CheckIfUserAuthorized(senderId); ok {
			log.Printf("%+v est autorisé\n", user)
			return name, true
		} else {
			log.Printf("%+v n'est pas autorisé\n", user)
			return "", false
		}
	}
	return "", false
}
