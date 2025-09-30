package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	switch r.Method {
	case http.MethodGet:
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

	case http.MethodPost:
		r.ParseForm()

		// Parse basic fields
		charName := r.FormValue("charname")
		playerName := r.FormValue("playername")
		race := r.FormValue("race")
		class := r.FormValue("classlevel")
		background := r.FormValue("background")
		level, _ := strconv.Atoi(r.FormValue("level"))
		expPoints, _ := strconv.Atoi(r.FormValue("experiencepoints"))

		// Parse ability scores
		strength, _ := strconv.Atoi(r.FormValue("Strengthscore"))
		dexterity, _ := strconv.Atoi(r.FormValue("Dexterityscore"))
		constitution, _ := strconv.Atoi(r.FormValue("Constitutionscore"))
		intelligence, _ := strconv.Atoi(r.FormValue("Intelligencescore"))
		wisdom, _ := strconv.Atoi(r.FormValue("Wisdomscore"))
		charisma, _ := strconv.Atoi(r.FormValue("Charismascore"))

		// Apply racial modifiers
		raceKey := strings.ToLower(race)
		modifiers := models.RaceModifiers[raceKey]
		abilities := models.AbilityScores{
			Strength:     strength + modifiers["Strength"],
			Dexterity:    dexterity + modifiers["Dexterity"],
			Constitution: constitution + modifiers["Constitution"],
			Intelligence: intelligence + modifiers["Intelligence"],
			Wisdom:       wisdom + modifiers["Wisdom"],
			Charisma:     charisma + modifiers["Charisma"],
		}

		// Skill proficiencies (from form)
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

		// Als er geen skills in het formulier zijn ingevuld, gebruik class skills inclusief duplicaten
		if len(skillProficiencies) == 0 {
			skillProficiencies = append([]string{}, models.ClassSkills[strings.ToLower(class)]...)
		}

		// Check if character exists
		character, err := storage.GetCharacterByName(charName)
		if err != nil {
			// Create new character
			character = models.Character{
				Name:               charName,
				PlayerName:         playerName,
				Race:               race,
				Class:              class,
				Level:              level,
				Background:         background,
				ExperiencePoints:   expPoints,
				ProficiencyBonus:   models.CalculateProfBonus(level),
				Abilities:          abilities,
				SkillProficiencies: skillProficiencies,
			}
		} else {
			// Update existing character
			character.PlayerName = playerName
			character.Race = race
			character.Class = class
			character.Level = level
			character.Background = background
			character.ExperiencePoints = expPoints
			character.ProficiencyBonus = models.CalculateProfBonus(level)
			character.Abilities = abilities
			character.SkillProficiencies = skillProficiencies
		}

		// Calculate derived stats
		character.CalculateAllSkills()
		character.CalculateCombatStats()
		character.SetupSpellcasting()

		// Fetch equipment from API (fixed assignment)
		mainHand, offHand, armor, shield, err := api.GetEquipment()
		if err != nil {
			log.Println("Error fetching equipment:", err)
		} else {
			character.Equipment = models.Equipment{
				MainHand: mainHand,
				OffHand:  offHand,
				Armor:    armor,
				Shield:   shield,
			}
		}

		// Fetch spells for class
		spells, err := api.GetSpellsForClass(class, character.SpellSlots)
		if err != nil {
			log.Println("Error fetching spells:", err)
		} else {
			character.Spells = spells
		}

		// Save character
		if err := storage.SaveCharacter(character); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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
