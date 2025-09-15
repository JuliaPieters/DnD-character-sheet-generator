package main

import (
	"dnd-character-sheet/commands"
	"dnd-character-sheet/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Start de server
func main() {
	// Static files serveren
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", listCharactersHandler)           // lijst van alle characters
	http.HandleFunc("/character", characterDetailHandler) // detailpagina van één character
	http.HandleFunc("/create", createCharacterHandler)    // character aanmaken via POST
	http.HandleFunc("/create-form", showCreateFormHandler) // formulier tonen via GET

	log.Println("🌐 Server gestart op http://localhost:8080")
	log.Println("   Ga naar http://localhost:8080/create-form voor character creation")
	http.ListenAndServe(":8080", nil)
}

// Lijst van alle characters tonen
func listCharactersHandler(responseWriter http.ResponseWriter, request *http.Request) {
	characters, err := commands.ListCharacters()
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	templateFile := template.Must(template.ParseFiles("templates/characterList.html"))
	if err := templateFile.Execute(responseWriter, characters); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Detailpagina van één character tonen
func characterDetailHandler(responseWriter http.ResponseWriter, request *http.Request) {
	characterIDStr := request.URL.Query().Get("id")

	var characterPointer *models.Character

	if characterIDStr != "" {
		characterID, err := strconv.Atoi(characterIDStr)
		if err != nil {
			http.Error(responseWriter, "Invalid character ID", http.StatusBadRequest)
			return
		}

		character, err := commands.GetCharacter(characterID)
		if err != nil {
			http.Error(responseWriter, "Character not found", http.StatusNotFound)
			return
		}

		characterPointer = character
	} else {
		// Nieuw leeg character voor create form
		characterPointer = &models.Character{
			Abilities: models.AbilityScores{},
			Skills:    make(map[string]int),
			Equipment: models.Equipment{},
		}
	}

	templateFile := template.Must(template.ParseFiles("templates/charactersheet.html"))
	if err := templateFile.Execute(responseWriter, characterPointer); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Character creation form tonen (GET)
func showCreateFormHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "GET methode vereist", http.StatusMethodNotAllowed)
		return
	}

	// Data voor het formulier
	formData := struct {
		Races             map[string]map[string]int
		Classes           map[string][]string
		StandardArray     []int
		SpellcastingClasses map[string]string
		AbilityNames      []string
	}{
		Races:             models.RaceModifiers,
		Classes:           models.ClassSkills,
		StandardArray:     models.StandardArray,
		SpellcastingClasses: models.SpellcastingClasses,
		AbilityNames:      []string{"Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma"},
	}

	templateFile := template.Must(template.ParseFiles("templates/createCharacter.html"))
	if err := templateFile.Execute(responseWriter, formData); err != nil {
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

	// Basis velden uit het formulier
	characterName := request.FormValue("charname")
	characterRace := request.FormValue("race")
	characterClass := request.FormValue("class")
	characterBackground := request.FormValue("background")

	// Level uit form, default 1
	level := 1
	if levelStr := request.FormValue("level"); levelStr != "" {
		if parsedLevel, err := strconv.Atoi(levelStr); err == nil && parsedLevel > 0 && parsedLevel <= 20 {
			level = parsedLevel
		}
	}

	// Standard Array assignment uit dropdowns
	abilityAssignment := make([]string, 6)
	abilityAssignment[0] = request.FormValue("ability_15") // 15 gaat naar...
	abilityAssignment[1] = request.FormValue("ability_14") // 14 gaat naar...
	abilityAssignment[2] = request.FormValue("ability_13") // 13 gaat naar...
	abilityAssignment[3] = request.FormValue("ability_12") // 12 gaat naar...
	abilityAssignment[4] = request.FormValue("ability_10") // 10 gaat naar...
	abilityAssignment[5] = request.FormValue("ability_8")  // 8 gaat naar...

	// Fallback naar default assignment als leeg
	hasEmptyAssignment := false
	for _, assignment := range abilityAssignment {
		if assignment == "" {
			hasEmptyAssignment = true
			break
		}
	}

	if hasEmptyAssignment {
		abilityAssignment = []string{"Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma"}
	}

	// Skills uit checkboxes
	skillChoices := []string{}
	availableSkills := models.GetAvailableSkills(characterClass)
	
	for _, skill := range availableSkills {
		// Convert skill name naar form field name (spaces -> underscores, lowercase)
		fieldName := "skill_" + strings.ToLower(strings.ReplaceAll(skill, " ", "_"))
		if request.FormValue(fieldName) == "on" {
			skillChoices = append(skillChoices, skill)
		}
	}

	// Als geen skills gekozen, neem eerste 2 beschikbare
	if len(skillChoices) == 0 && len(availableSkills) >= 2 {
		skillChoices = availableSkills[:2]
	}

	// Character aanmaken met nieuwe functie
	err := commands.CreateCharacter(
		characterName,
		characterRace,
		characterClass,
		characterBackground,
		level,
		abilityAssignment,
		skillChoices,
	)

	if err != nil {
		// Voor nu: toon error op de pagina
		http.Error(responseWriter, "Fout bij aanmaken character: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect naar character lijst
	http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
}

// API endpoint voor getting available skills voor een class (AJAX support)
func getClassSkillsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	className := request.URL.Query().Get("class")
	if className == "" {
		http.Error(responseWriter, "Class parameter required", http.StatusBadRequest)
		return
	}

	skills := models.GetAvailableSkills(className)
	
	responseWriter.Header().Set("Content-Type", "application/json")
	// Simple JSON response (je kunt ook encoding/json package gebruiken)
	responseWriter.Write([]byte(`{"skills": [`))
	for i, skill := range skills {
		if i > 0 {
			responseWriter.Write([]byte(`, `))
		}
		responseWriter.Write([]byte(`"` + skill + `"`))
	}
	responseWriter.Write([]byte(`]}`))
}