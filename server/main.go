package main

import (
	"dnd-character-sheet/commands"
	"fmt"
	"html/template"
	"net/http"
)

// Struct voor character data in de template
type Character struct {
	Name  string
	Level int
	Race  string
	Class string
}

func main() {
	// Static files serveren (CSS, JS, afbeeldingen)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))

	// Routes
	http.HandleFunc("/", listCharactersHandler)
	http.HandleFunc("/create", createCharacterHandler)

	fmt.Println("🌐 Server gestart op http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Handler om alle characters te tonen
func listCharactersHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Haal characters op (zorg dat ListCharactersMap een slice of map teruggeeft)
	characters, err := commands.ListCharactersMap()
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// Laad template
	tmpl := template.Must(template.ParseFiles("../templates/charactersheet.html"))

	// Render template met characters
	err = tmpl.Execute(responseWriter, characters)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler om een character aan te maken
func createCharacterHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "POST methode vereist", http.StatusMethodNotAllowed)
		return
	}

	// Haal form values op
	characterName := request.FormValue("charname")
	characterRace := request.FormValue("race")
	characterClass := request.FormValue("classlevel")
	characterBackground := request.FormValue("background")

	// Standaard ability scores
	strength := 10
	dexterity := 10
	constitution := 10
	intelligence := 10
	wisdom := 10
	charisma := 10
	skills := []string{}

	// Maak character aan via bestaande CLI functie
	if err := commands.CreateCharacter(characterName, characterRace, characterClass, characterBackground, 1,
		strength, dexterity, constitution, intelligence, wisdom, charisma, skills); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect terug naar homepage na aanmaken
	http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
}
