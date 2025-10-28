package main

import (
	"dnd-character-sheet/api"
	"dnd-character-sheet/application"
	"dnd-character-sheet/commands"
	"dnd-character-sheet/domain"
	"dnd-character-sheet/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
	characterService := &application.CharacterService{}
	spellService := &application.SpellService{}

	switch r.Method {
	case http.MethodGet:
		characterID := r.URL.Query().Get("id")
		var character *domain.Character

		if characterID != "" {
			character, _ = storage.GetCharacterByName(characterID)
			if character == nil {
				http.Error(w, "Character not found", http.StatusNotFound)
				return
			}
		} else {
			character = &domain.Character{
				Abilities: domain.AbilityScores{},
				Equipment: domain.Equipment{},
				Skills:    make(map[string]int),
				Spells:    []domain.Spell{},
			}
		}

		if err := templates.ExecuteTemplate(w, "charactersheet.html", character); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		character, err := setupCharacterFromForm(r, spellService, characterService)
		if err != nil {
			http.Error(w, "Failed to setup character", http.StatusInternalServerError)
			return
		}

		spellService.SetupSpellcasting(character)
		if err := commands.GiveStartingSpells(character); err != nil {
			log.Println("Failed to give starting spells:", err)
		}

		handleEquipmentAndCombat(character, characterService)
		handleSpells(character)

		if err := storage.SaveCharacter(character); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func parseSkillProficiencies(r *http.Request) []string {
	var skillProficiencies []string
	for _, skill := range []string{
		"Acrobatics", "Animal Handling", "Arcana", "Athletics", "Deception",
		"History", "Insight", "Intimidation", "Investigation", "Medicine",
		"Nature", "Perception", "Performance", "Persuasion", "Religion",
		"Sleight of Hand", "Stealth", "Survival",
	} {
		if r.FormValue(skill+"-prof") == "on" {
			skillProficiencies = append(skillProficiencies, skill)
		}
	}
	return skillProficiencies
}

func getFormInt(r *http.Request, field string, defaultValue int) int {
	if val, err := strconv.Atoi(r.FormValue(field)); err == nil {
		return val
	}
	return defaultValue
}

func setupCharacterFromForm(r *http.Request, spellService *application.SpellService, service *application.CharacterService) (*domain.Character, error) {
	charName := r.FormValue("charname")
	race := r.FormValue("race")
	class := r.FormValue("classlevel")
	background := r.FormValue("background")
	level := getFormInt(r, "level", 1)
	expPoints := getFormInt(r, "experiencepoints", 0)
	speed := getFormInt(r, "Speed", 30)

	abilityScores := []int{
		getFormInt(r, "Strengthscore", 10),
		getFormInt(r, "Dexterityscore", 10),
		getFormInt(r, "Constitutionscore", 10),
		getFormInt(r, "Intelligencescore", 10),
		getFormInt(r, "Wisdomscore", 10),
		getFormInt(r, "Charismascore", 10),
	}

	skillProficiencies := parseSkillProficiencies(r)
	if len(skillProficiencies) == 0 {
		skillProficiencies = service.GetAvailableSkills(class)
	}

	existingChar, _ := storage.GetCharacterByName(charName)
	var char *domain.Character
	if existingChar == nil {
		allChars, _ := storage.LoadCharacters()
		newID := len(allChars) + 1

		char = service.NewCharacter(application.NewCharacterParams{
			ID:            newID,
			Name:          charName,
			Race:          race,
			Class:         class,
			Background:    background,
			Level:         level,
			AbilityScores: abilityScores,
			SkillChoices:  skillProficiencies,
			SpellService:  spellService,
		})
	} else {
		char = existingChar
		char.PlayerName = r.FormValue("playername")
		char.Race = race
		char.Class = class
		char.Level = level
		char.Background = background
		char.ExperiencePoints = expPoints
		char.Abilities = domain.AbilityScores{
			Strength:     abilityScores[0],
			Dexterity:    abilityScores[1],
			Constitution: abilityScores[2],
			Intelligence: abilityScores[3],
			Wisdom:       abilityScores[4],
			Charisma:     abilityScores[5],
		}
		char.SkillProficiencies = skillProficiencies
		char.Speed = speed
		service.UpdateModifiers(char)
		service.CalculateAllSkills(char)
		service.CalculateCombatStats(char)
	}

	char.PlayerName = r.FormValue("playername")
	char.Speed = speed
	char.ExperiencePoints = expPoints

	if char.Equipment == (domain.Equipment{}) {
		char.Equipment = domain.Equipment{}
	}

	if char.Skills == nil {
		char.Skills = make(map[string]int)
	}
	if char.Spells == nil {
		char.Spells = []domain.Spell{}
	}

	return char, nil
}

func handleEquipmentAndCombat(character *domain.Character, service *application.CharacterService) {
	allWeapons, armor, shield, err := api.GetAllEquipment()
	if err != nil {
		log.Println("Error fetching equipment:", err)
		return
	}

	if character.Equipment.MainHand == nil && len(allWeapons) > 0 {
		character.Equipment.MainHand = allWeapons[0]
	} else if character.Equipment.MainHand != nil {
		fillWeaponData(character.Equipment.MainHand, allWeapons)
	}

	if character.Equipment.OffHand == nil && len(allWeapons) > 1 {
		character.Equipment.OffHand = allWeapons[1]
	} else if character.Equipment.OffHand != nil {
		fillWeaponData(character.Equipment.OffHand, allWeapons)
	}

	if character.Equipment.Armor == nil && armor != nil {
		character.Equipment.Armor = armor
	}

	if character.Equipment.Shield == nil && shield != nil {
		character.Equipment.Shield = shield
	}

	service.CalculateCombatStats(character)
}

func fillWeaponData(existing *domain.Weapon, allWeapons []*domain.Weapon) {
	for _, w := range allWeapons {
		if w.Name == existing.Name {
			if existing.DamageDie == "" {
				existing.DamageDie = w.DamageDie
			}
			if existing.Range == "" {
				existing.Range = w.Range
			}
			if !existing.TwoHanded && w.TwoHanded {
				existing.TwoHanded = true
			}
			if !existing.IsFinesse && w.IsFinesse {
				existing.IsFinesse = true
			}
			break
		}
	}
}

func handleSpells(character *domain.Character) {
	spells, err := api.GetSpellsForClass(character.Class, character.SpellSlots)
	if err != nil {
		log.Println("Error fetching spells:", err)
		return
	}
	character.Spells = spells
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/character", characterHandler)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
