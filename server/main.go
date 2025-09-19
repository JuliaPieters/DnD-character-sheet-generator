package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
)

// Templates
var templates = template.Must(template.ParseGlob("../templates/*.html"))

// ------------------------
// Handlers
// ------------------------
func listHandler(w http.ResponseWriter, r *http.Request) {
	// Load characters van JSON
	characters, err := storage.LoadCharacters()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// characters is map[string]models.Character
	// log om te checken
	log.Println("Loaded characters:")
	for k, v := range characters {
		log.Printf("%s -> %+v\n", k, v)
	}

	err = templates.ExecuteTemplate(w, "characterList.html", characters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func characterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		name := r.URL.Query().Get("id")
		var character models.Character
		if name != "" {
			c, err := storage.GetCharacterByName(name)
			if err == nil {
				character = c
			}
		}

		err := templates.ExecuteTemplate(w, "charactersheet.html", character)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		level, _ := strconv.Atoi(r.FormValue("level"))
		xp, _ := strconv.Atoi(r.FormValue("experiencepoints"))
		str, _ := strconv.Atoi(r.FormValue("Strengthscore"))
		dex, _ := strconv.Atoi(r.FormValue("Dexterityscore"))
		con, _ := strconv.Atoi(r.FormValue("Constitutionscore"))
		intt, _ := strconv.Atoi(r.FormValue("Intelligencescore"))
		wis, _ := strconv.Atoi(r.FormValue("Wisdomscore"))
		cha, _ := strconv.Atoi(r.FormValue("Charismascore"))

		var skills []string
		for _, skill := range []string{
			"Acrobatics", "Animal Handling", "Arcana", "Athletics",
			"Deception", "History", "Insight", "Intimidation",
			"Investigation", "Medicine", "Nature", "Perception",
			"Performance", "Persuasion", "Religion", "Sleight of Hand",
			"Stealth", "Survival",
		} {
			if r.FormValue(skill+"-prof") == "on" {
				skills = append(skills, skill)
			}
		}

		char := models.Character{
			Name:       r.FormValue("charname"),
			PlayerName: r.FormValue("playername"),
			Race:       r.FormValue("race"),
			Class:      r.FormValue("classlevel"),
			Level:      level,
			Background: r.FormValue("background"),
			Alignment:  r.FormValue("alignment"),
			ExperiencePoints: xp,
			Abilities: models.AbilityScores{
				Strength:     str,
				Dexterity:    dex,
				Constitution: con,
				Intelligence: intt,
				Wisdom:       wis,
				Charisma:     cha,
			},
			SkillProficiencies: skills,
			Equipment: models.Equipment{
				Weapons: []models.Weapon{},
				Armor:   nil,
				Shield:  nil,
			},
			Personality: r.FormValue("personality"),
			Ideals:      r.FormValue("ideals"),
			Bonds:       r.FormValue("bonds"),
			Flaws:       r.FormValue("flaws"),
			Features:    r.FormValue("features"),
		}

		// Opslaan via storage
		err := storage.SaveCharacter(char)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// ------------------------
// Main
// ------------------------
func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/character", characterHandler)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
