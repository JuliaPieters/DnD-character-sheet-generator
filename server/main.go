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

var templates = template.Must(template.ParseGlob("../templates/*.html"))

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

		charName := r.FormValue("charname")
		playerName := r.FormValue("playername")
		race := r.FormValue("race")
		class := r.FormValue("classlevel")
		background := r.FormValue("background")
		level, _ := strconv.Atoi(r.FormValue("level"))
		expPoints, _ := strconv.Atoi(r.FormValue("experiencepoints"))

		strength, _ := strconv.Atoi(r.FormValue("Strengthscore"))
		dexterity, _ := strconv.Atoi(r.FormValue("Dexterityscore"))
		constitution, _ := strconv.Atoi(r.FormValue("Constitutionscore"))
		intelligence, _ := strconv.Atoi(r.FormValue("Intelligencescore"))
		wisdom, _ := strconv.Atoi(r.FormValue("Wisdomscore"))
		charisma, _ := strconv.Atoi(r.FormValue("Charismascore"))
		speed, _ := strconv.Atoi(r.FormValue("Speed"))

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

		if len(skillProficiencies) == 0 {
			skillProficiencies = append([]string{}, models.ClassSkills[strings.ToLower(class)]...)
		}

		character, err := storage.GetCharacterByName(charName)
		if err != nil {
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
				Speed:              speed,
			}
		} else {
			character.PlayerName = playerName
			character.Race = race
			character.Class = class
			character.Level = level
			character.Background = background
			character.ExperiencePoints = expPoints
			character.ProficiencyBonus = models.CalculateProfBonus(level)
			character.Abilities = abilities
			character.SkillProficiencies = skillProficiencies
			character.Speed = speed
		}

		character.StrengthMod = character.Abilities.Modifier("Strength")
		character.DexterityMod = character.Abilities.Modifier("Dexterity")
		character.ConstitutionMod = character.Abilities.Modifier("Constitution")
		character.IntelligenceMod = character.Abilities.Modifier("Intelligence")
		character.WisdomMod = character.Abilities.Modifier("Wisdom")
		character.CharismaMod = character.Abilities.Modifier("Charisma")

		character.CalculateAllSkills()
		character.CalculateCombatStats()
		character.SetupSpellcasting()

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

		spells, err := api.GetSpellsForClass(class, character.SpellSlots)
		if err != nil {
			log.Println("Error fetching spells:", err)
		} else {
			character.Spells = spells
		}

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

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/character", characterHandler)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
