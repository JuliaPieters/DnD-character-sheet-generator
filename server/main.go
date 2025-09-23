package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"dnd-character-sheet/api"
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
)

// Templates
var templates = template.Must(template.ParseGlob("../templates/*.html"))

// ------------------------
// Handlers
// ------------------------
func listHandler(w http.ResponseWriter, r *http.Request) {
	characters, err := storage.LoadCharacters()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Loaded characters:")
	for characterName, characterData := range characters {
		log.Printf("%s -> %+v\n", characterName, characterData)
	}

	err = templates.ExecuteTemplate(w, "characterList.html", characters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func characterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		characterID := r.URL.Query().Get("id")
		var character models.Character
		if characterID != "" {
			foundCharacter, err := storage.GetCharacterByName(characterID)
			if err == nil {
				character = foundCharacter
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
		experiencePoints, _ := strconv.Atoi(r.FormValue("experiencepoints"))
		strengthScore, _ := strconv.Atoi(r.FormValue("Strengthscore"))
		dexterityScore, _ := strconv.Atoi(r.FormValue("Dexterityscore"))
		constitutionScore, _ := strconv.Atoi(r.FormValue("Constitutionscore"))
		intelligenceScore, _ := strconv.Atoi(r.FormValue("Intelligencescore"))
		wisdomScore, _ := strconv.Atoi(r.FormValue("Wisdomscore"))
		charismaScore, _ := strconv.Atoi(r.FormValue("Charismascore"))

		// Verzamel skill proficiencies
		var skillProficiencies []string
		for _, skill := range []string{
			"Acrobatics", "Animal Handling", "Arcana", "Athletics",
			"Deception", "History", "Insight", "Intimidation",
			"Investigation", "Medicine", "Nature", "Perception",
			"Performance", "Persuasion", "Religion", "Sleight of Hand",
			"Stealth", "Survival",
		} {
			if r.FormValue(skill+"-prof") == "on" {
				skillProficiencies = append(skillProficiencies, skill)
			}
		}

		// Maak nieuw Character object
		character := models.Character{
			Name:               r.FormValue("charname"),
			PlayerName:         r.FormValue("playername"),
			Race:               r.FormValue("race"),
			Class:              r.FormValue("classlevel"),
			Level:              level,
			Background:         r.FormValue("background"),
			Alignment:          r.FormValue("alignment"),
			ExperiencePoints:   experiencePoints,
			Abilities: models.AbilityScores{
				Strength:     strengthScore,
				Dexterity:    dexterityScore,
				Constitution: constitutionScore,
				Intelligence: intelligenceScore,
				Wisdom:       wisdomScore,
				Charisma:     charismaScore,
			},
			SkillProficiencies: skillProficiencies,
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

		// ---------- Voeg spells toe ----------
		if character.SpellSlots == nil {
			character.SetupSpellcasting() // Zorg dat spell slots aanwezig zijn
		}

		spells, err := api.GetSpellsForClass(character.Class, character.SpellSlots)
		if err != nil {
			log.Println("Error fetching spells:", err)
		} else {
			character.Spells = spells
		}

		// ---------- Voeg equipment toe ----------
		weapons, armor, shield, err := api.GetEquipment()
		if err != nil {
			log.Println("Error fetching equipment:", err)
		} else {
			character.Equipment = models.Equipment{
				Weapons: weapons,
				Armor:   armor,
				Shield:  shield,
			}
		}

		// Sla character op
		err = storage.SaveCharacter(character)
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
