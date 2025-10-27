package commands

import (
	"dnd-character-sheet/domain"
	"dnd-character-sheet/storage"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var SpellcastingClasses = map[string]bool{
	"bard":     true,
	"cleric":   true,
	"druid":    true,
	"paladin":  true,
	"sorcerer": true,
	"warlock":  true,
	"wizard":   true,
}

var PreparedCasters = map[string]bool{
	"cleric":  true,
	"druid":   true,
	"paladin": true,
	"wizard":  true,
}

var SpellList []domain.Spell
var SpellClasses = map[string][]string{}

func LoadSpellsFromCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open spells file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read spells file: %v", err)
	}

	SpellList = []domain.Spell{}
	SpellClasses = map[string][]string{}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(row[0]))
		level, _ := strconv.Atoi(row[1])
		classes := strings.Split(strings.ToLower(row[2]), ",")
		for j := range classes {
			classes[j] = strings.TrimSpace(classes[j])
		}

		SpellList = append(SpellList, domain.Spell{
			Name:  name,
			Level: level,
		})
		SpellClasses[name] = classes
	}
	return nil
}

func FindSpellByName(name string) *domain.Spell {
	name = strings.ToLower(name)
	for i := range SpellList {
		if SpellList[i].Name == name {
			return &SpellList[i]
		}
	}
	return nil
}

func FindSpellsForClass(class string) []domain.Spell {
	class = strings.ToLower(class)
	var spells []domain.Spell
	for _, s := range SpellList {
		for _, c := range SpellClasses[s.Name] {
			if c == class {
				spells = append(spells, s)
				break
			}
		}
	}
	return spells
}

func GiveStartingSpells(character *domain.Character) error {
	if character.CanPrepareSpells {
		spells := FindSpellsForClass(character.Class)
		for _, s := range spells {
			if s.Level == 0 {
				character.Spells = append(character.Spells, domain.Spell{
					Name:     s.Name,
					Level:    s.Level,
					Prepared: false,
				})
			}
		}
	}
	SetupSpellcasting(character)
	return storage.SaveCharacter(character)
}

func LearnSpell(characterName, spellName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}
	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character \"%s\" not found", characterName)
	}
	characterPtr := &character

	if !SpellcastingClasses[characterPtr.Class] {
		return fmt.Errorf("this class can't cast spells")
	}
	if characterPtr.CanPrepareSpells {
		return fmt.Errorf("this class prepares spells and can't learn them")
	}

	spell := FindSpellByName(spellName)
	if spell == nil {
		return fmt.Errorf("spell '%s' not found in spell list", spellName)
	}

	valid := false
	for _, c := range SpellClasses[spell.Name] {
		if c == strings.ToLower(characterPtr.Class) {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("%s cannot learn %s", characterPtr.Class, spellName)
	}

	for _, s := range characterPtr.Spells {
		if s.Name == spell.Name {
			return fmt.Errorf("character '%s' already knows spell '%s'", characterName, spell.Name)
		}
	}

	characterPtr.Spells = append(characterPtr.Spells, domain.Spell{
		Name:     spell.Name,
		Level:    spell.Level,
		Prepared: false,
	})
	if err := storage.SaveCharacter(characterPtr); err != nil {
		return err
	}
	fmt.Printf("Learned spell %s\n", spell.Name)
	return nil
}

func PrepareSpell(characterName, spellName string, spellLevel int) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}
	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf(`character "%s" not found`, characterName)
	}
	characterPtr := &character

	if !SpellcastingClasses[characterPtr.Class] {
		return fmt.Errorf("this class can't cast spells")
	}
	if !characterPtr.CanPrepareSpells {
		return fmt.Errorf("this class learns spells and can't prepare them")
	}

	spellIndex := -1
	for i, s := range characterPtr.Spells {
		if s.Name == spellName {
			spellIndex = i
			break
		}
	}

	if spellIndex == -1 && characterPtr.CanPrepareSpells {
		classSpells := FindSpellsForClass(characterPtr.Class)
		for _, s := range classSpells {
			if s.Name == spellName {
				spellIndex = -2
				break
			}
		}
		if spellIndex == -1 {
			return fmt.Errorf("spell '%s' not available for class '%s'", spellName, characterPtr.Class)
		}
	}

	spell := FindSpellByName(spellName)
	if spell == nil {
		return fmt.Errorf("spell '%s' not found in spell list", spellName)
	}

	if spellLevel < spell.Level {
		return fmt.Errorf("the spell has higher level than the available spell slots")
	}
	if slots, ok := characterPtr.SpellSlots[spellLevel]; !ok || slots == 0 {
		return fmt.Errorf("no available spell slots of level %d", spellLevel)
	}

	if spellIndex == -2 {
		characterPtr.Spells = append(characterPtr.Spells, domain.Spell{
			Name:     spell.Name,
			Level:    spellLevel,
			Prepared: true,
		})
	} else {
		characterPtr.Spells[spellIndex].Prepared = true
		characterPtr.Spells[spellIndex].Level = spellLevel
	}

	if err := storage.SaveCharacter(characterPtr); err != nil {
		return err
	}
	fmt.Printf("Prepared spell %s\n", spellName)
	return nil
}

func SetupSpellcasting(c *domain.Character) {
	class := strings.ToLower(c.Class)
	if !SpellcastingClasses[class] {
		c.SpellSlots = nil
		c.CanPrepareSpells = false
		return
	}
	c.CanPrepareSpells = PreparedCasters[class]
	c.SpellSlots = GenerateSpellSlots(class, c.Level)
	if class == "warlock" {
		c.SpellSlots[0] = 4
	}
}

func GenerateSpellSlots(class string, level int) map[int]int {
	slots := make(map[int]int)
	class = strings.ToLower(class)

	fullCasters := map[string]bool{
		"wizard": true, "cleric": true, "druid": true, "bard": true, "sorcerer": true,
	}

	if fullCasters[class] {
		switch {
		case level <= 3:
			slots[0] = 3
		case level <= 9:
			slots[0] = 4
		default:
			slots[0] = 5
		}

		fullCasterSlots := [][]int{
			{2, 0, 0, 0, 0, 0, 0, 0, 0},
			{3, 0, 0, 0, 0, 0, 0, 0, 0},
			{4, 2, 0, 0, 0, 0, 0, 0, 0},
			{4, 3, 0, 0, 0, 0, 0, 0, 0},
			{4, 3, 2, 0, 0, 0, 0, 0, 0},
			{4, 3, 3, 0, 0, 0, 0, 0, 0},
			{4, 3, 3, 1, 0, 0, 0, 0, 0},
			{4, 3, 3, 2, 0, 0, 0, 0, 0},
			{4, 3, 3, 3, 0, 0, 0, 0, 0},
			{4, 3, 3, 3, 2, 0, 0, 0, 0},
			{4, 3, 3, 3, 3, 0, 0, 0, 0},
			{4, 3, 3, 3, 3, 1, 0, 0, 0},
			{4, 3, 3, 3, 3, 2, 0, 0, 0},
			{4, 3, 3, 3, 3, 2, 1, 0, 0},
			{4, 3, 3, 3, 3, 3, 1, 0, 0},
			{4, 3, 3, 3, 3, 3, 2, 0, 0},
			{4, 3, 3, 3, 3, 3, 2, 1, 0},
			{4, 3, 3, 3, 3, 3, 2, 2, 1},
			{4, 3, 3, 3, 3, 3, 3, 2, 2},
			{4, 3, 3, 3, 3, 2, 2, 1, 1},
		}

		if class == "cleric" || class == "druid" {
			if level > 10 {
				level = 10
			}
		}

		if level >= 1 && level <= len(fullCasterSlots) {
			for i, count := range fullCasterSlots[level-1] {
				if count > 0 {
					slots[i+1] = count
				}
			}
		}
		return slots
	}

	if class == "paladin" || class == "ranger" {
		halfCaster := []struct {
			slotLevel int
			reqLevel  int
			slots     int
		}{
			{1, 2, 4}, {2, 5, 3}, {3, 9, 3}, {4, 13, 3}, {5, 17, 2},
		}
		for _, h := range halfCaster {
			if level >= h.reqLevel {
				slots[h.slotLevel] = h.slots
			}
		}
		return slots
	}

	if class == "warlock" {
		switch {
		case level == 1:
			slots[1] = 1
		case level == 2:
			slots[1] = 2
		case level >= 3 && level <= 4:
			slots[2] = 2
		case level >= 5 && level <= 6:
			slots[3] = 2
		case level >= 7 && level <= 8:
			slots[4] = 2
		case level >= 9 && level <= 10:
			slots[5] = 2
		case level >= 11 && level <= 16:
			slots[5] = 3
		case level >= 17 && level <= 19:
			slots[5] = 4
		case level == 20:
			slots[5] = 4
			slots[6] = 1
			slots[7] = 1
			slots[8] = 1
			slots[9] = 1
		}
		return slots
	}

	return slots
}