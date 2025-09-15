package main

import (
	"dnd-character-sheet/commands"
	"dnd-character-sheet/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Start de server
func main() {
	// Static files serveren
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", listCharactersHandler)           // lijst van alle characters
	http.HandleFunc("/character", characterDetailHandler) // detailpagina van één character
	http.HandleFunc("/create", createCharacterHandler)    // character aanmaken via POST

	log.Println("🌐 Server gestart op http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Lijst van alle characters tonen
func listCharactersHandler(responseWriter http.ResponseWriter, request *http.Request) {
	charactersMap, err := commands.ListCharactersMap()
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	templateFile := template.Must(template.ParseFiles("templates/characterList.html"))
	if err := templateFile.Execute(responseWriter, charactersMap); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Detailpagina van één character tonen
func characterDetailHandler(responseWriter http.ResponseWriter, request *http.Request) {
	characterID := request.URL.Query().Get("id")

	var characterPointer *models.Character

	if characterID != "" {
		charactersMap, err := commands.ListCharactersMap()
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		foundCharacter, exists := charactersMap[characterID]
		if !exists {
			http.Error(responseWriter, "Character not found", http.StatusNotFound)
			return
		}

		characterPointer = &foundCharacter
	} else {
		// Nieuw leeg character voor create
		characterPointer = &models.Character{
			Abilities: models.AbilityScores{},
			Skills:    make(map[string]int),
			Attacks:   []models.Attack{{"", "", ""}, {"", "", ""}, {"", "", ""}},
		}
	}

	templateFile := template.Must(template.ParseFiles("templates/charactersheet.html"))
	if err := templateFile.Execute(responseWriter, characterPointer); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Character aanmaken via POST
func createCharacterHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "POST methode vereist", http.StatusMethodNotAllowed)
		return
	}

	// Velden uit het formulier
	characterName := request.FormValue("charname")
	playerName := request.FormValue("playername")
	characterRace := request.FormValue("race")
	characterClass := request.FormValue("classlevel")
	characterBackground := request.FormValue("background")
	alignment := request.FormValue("alignment")

	// Ability scores uit form, fallback op 10
	strength, _ := strconv.Atoi(request.FormValue("Strength"))
	dexterity, _ := strconv.Atoi(request.FormValue("Dexterity"))
	constitution, _ := strconv.Atoi(request.FormValue("Constitution"))
	intelligence, _ := strconv.Atoi(request.FormValue("Intelligence"))
	wisdom, _ := strconv.Atoi(request.FormValue("Wisdom"))
	charisma, _ := strconv.Atoi(request.FormValue("Charisma"))

	// Standaard skills, kan later uitgebreid met checkboxes
	skills := []string{}

	// Level standaard op 1
	level := 1

	err := commands.CreateCharacter(
		characterName,
		playerName,
		characterRace,
		characterClass,
		characterBackground,
		alignment,
		level,
		strength,
		dexterity,
		constitution,
		intelligence,
		wisdom,
		charisma,
		skills,
	)

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
}
