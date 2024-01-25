package models

type Movie struct {
	IdUser   string `json:"id_user"`
	NameUser string `json:"name_user"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	State    string `json:"state"`
}

var MOVIE_USER = make(map[string]Movie)

func (m *Movie) UpdateMovies(senderId string) {
	// Récupère l'objet Movies de la map, si existant.
	movie, exists := MOVIE_USER[senderId]
	if !exists {
		// Si non existant, initialise un nouvel objet Movies.
		movie = Movie{IdUser: senderId}
	}

	// Met à jour les champs si de nouvelles valeurs sont fournies.
	if m.NameUser != "" {
		movie.NameUser = m.NameUser
	}
	if m.Type != "" {
		movie.Type = m.Type
	}
	if m.Title != "" {
		movie.Title = m.Title
	}
	if m.Year != 0 {
		movie.Year = m.Year
	}

	if m.State != "" || m.State != movie.State {
		movie.State = m.State
	}

	// Remet l'objet modifié dans la map.
	MOVIE_USER[senderId] = movie
}

func GetMovies(senderId string) Movie {
	if movie, exists := MOVIE_USER[senderId]; exists {
		return movie
	}
	return Movie{State: "hello"}
}
