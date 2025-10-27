package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"dnd-character-sheet/api"
	"dnd-character-sheet/application"
	"dnd-character-sheet/commands"
	"dnd-character-sheet/domain"
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
	for name, c := range characters {
		log.Printf("%s -> %+v\n", name, c)
	}

	if err := templates.ExecuteTemplate(w, "characterList.html", characters); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func characterHandler(w http.ResponseWriter, r *http.Request) {
	characterService := application.CharacterService{}
	spellService := application.SpellService{}

	switch r.Method {
	case http.MethodGet:
		characterID := r.URL.Query().Get("id")
		var character *domain.Character
		if characterID != "" {
			foundCharacter, err := storage.GetCharacterByName(characterID)
			if err == nil {
				character = foundCharacter
			}
		}

		if err := templates.ExecuteTemplate(w, "charactersheet.html", character); err != nil {
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

		abilityScores := []int{strength, dexterity, constitution, intelligence, wisdom, charisma}

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
			skillProficiencies = characterService.GetAvailableSkills(class)
		}

		existingChar, err := storage.GetCharacterByName(charName)
		var character *domain.Character
		if err != nil || existingChar == nil {
			allChars, _ := storage.LoadCharacters()
			newID := len(allChars) + 1

			character = characterService.NewCharacter(
				newID,
				charName,
				race,
				class,
				background,
				level,
				abilityScores,
				skillProficiencies,
				&spellService, 
			)
			character.PlayerName = playerName
			character.Speed = speed
			character.ExperiencePoints = expPoints
		} else {
			character = existingChar
			character.PlayerName = playerName
			character.Race = race
			character.Class = class
			character.Level = level
			character.Background = background
			character.ExperiencePoints = expPoints
			character.Abilities = domain.AbilityScores{
				Strength:     strength,
				Dexterity:    dexterity,
				Constitution: constitution,
				Intelligence: intelligence,
				Wisdom:       wisdom,
				Charisma:     charisma,
			}
			character.SkillProficiencies = skillProficiencies
			character.Speed = speed

			characterService.UpdateModifiers(character)
			characterService.CalculateAllSkills(character)
			characterService.CalculateCombatStats(character)
		}

		spellService.SetupSpellcasting(character)
		if err := commands.GiveStartingSpells(character); err != nil {
			log.Println("Failed to give starting spells:", err)
		}

		mainHand, offHand, armor, shield, err := api.GetEquipment()
		if err != nil {
			log.Println("Error fetching equipment:", err)
		} else {
			character.Equipment = domain.Equipment{
				MainHand: mainHand,
				OffHand:  offHand,
				Armor:    armor,
				Shield:   shield,
			}
			characterService.CalculateCombatStats(character)
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
