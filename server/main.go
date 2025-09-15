package main

import (
	"dnd-character-sheet/commands"
	"fmt"
	"html/template"
	"net/http"
)

// Struct die precies past bij de template
type TemplateCharacter struct {
	Name      string
	Race      string
	Class     string
	Level     int
	Background string
	ProfBonus int
	Abilities struct {
		Strength     int
		Dexterity    int
		Constitution int
		Intelligence int
		Wisdom       int
		Charisma     int
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", listCharactersHandler)
	http.HandleFunc("/character", characterDetailHandler)

	fmt.Println("🌐 Server gestart op http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func listCharactersHandler(w http.ResponseWriter, r *http.Request) {
	charactersMap, err := commands.ListCharactersMap()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/characterList.html"))
	if err := tmpl.Execute(w, charactersMap); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func characterDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var tmplChar TemplateCharacter

	if id != "" {
		charactersMap, err := commands.ListCharactersMap()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		foundCharacter, exists := charactersMap[id]
		if !exists {
			http.Error(w, "Character not found", http.StatusNotFound)
			return
		}

		// Vul template struct
		tmplChar.Name = foundCharacter.Name
		tmplChar.Race = foundCharacter.Race
		tmplChar.Class = foundCharacter.Class
		tmplChar.Level = foundCharacter.Level
		tmplChar.Background = foundCharacter.Background
		tmplChar.ProfBonus = foundCharacter.ProfBonus
		tmplChar.Abilities.Strength = foundCharacter.Abilities.Strength
		tmplChar.Abilities.Dexterity = foundCharacter.Abilities.Dexterity
		tmplChar.Abilities.Constitution = foundCharacter.Abilities.Constitution
		tmplChar.Abilities.Intelligence = foundCharacter.Abilities.Intelligence
		tmplChar.Abilities.Wisdom = foundCharacter.Abilities.Wisdom
		tmplChar.Abilities.Charisma = foundCharacter.Abilities.Charisma
	}

	tmpl := template.Must(template.ParseFiles("templates/charactersheet.html"))
	if err := tmpl.Execute(w, tmplChar); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
