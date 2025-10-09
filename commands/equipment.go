package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// ------------------------
// Maps
// ------------------------
var Armors = map[string]models.Armor{}
var Shields = map[string]models.Shield{}
var Weapons = map[string]models.Weapon{}

var DefaultArmorStats = map[string]models.Armor{
	// Light
	"padded":          {ArmorClass: 11, DexBonus: true, MaxDexBonus: 0},
	"leather":         {ArmorClass: 11, DexBonus: true, MaxDexBonus: 0},
	"studded leather": {ArmorClass: 12, DexBonus: true, MaxDexBonus: 0},

	// Medium
	"hide":        {ArmorClass: 12, DexBonus: true, MaxDexBonus: 2},
	"chain shirt": {ArmorClass: 13, DexBonus: true, MaxDexBonus: 2},
	"scale mail":  {ArmorClass: 14, DexBonus: true, MaxDexBonus: 2},
	"breastplate": {ArmorClass: 14, DexBonus: true, MaxDexBonus: 2},
	"half plate":  {ArmorClass: 15, DexBonus: true, MaxDexBonus: 2},

	// Heavy
	"ring mail":  {ArmorClass: 14, DexBonus: false},
	"chain mail": {ArmorClass: 16, DexBonus: false},
	"splint":     {ArmorClass: 17, DexBonus: false},
	"plate":      {ArmorClass: 18, DexBonus: false},

	// Shields
	"shield": {ArmorClass: 2, DexBonus: false},
}

// ------------------------
// Helpers
// ------------------------
func normalizeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.TrimSuffix(name, " armor")
	name = strings.TrimSpace(name)
	return name
}

// ------------------------
// CSV Loader
// ------------------------
func LoadEquipmentCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open equipment CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read equipment CSV: %w", err)
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		originalName := strings.TrimSpace(record[0])
		eqType := strings.ToLower(strings.TrimSpace(record[1]))
		key := normalizeName(originalName)

		switch eqType {
		case "armor":
			if key == "shield" {
				stats := DefaultArmorStats["shield"]
				shield := models.Shield{
					Name:       strings.ToLower(originalName),
					ArmorClass: stats.ArmorClass,
				}
				Shields[key] = shield
				Shields[strings.ToLower(originalName)] = shield
			} else {
				stats, ok := DefaultArmorStats[key]
				if !ok {
					stats = models.Armor{ArmorClass: 10, DexBonus: true}
				}
				armor := models.Armor{
					Name:        key, 
					ArmorClass:  stats.ArmorClass,
					DexBonus:    stats.DexBonus,
					MaxDexBonus: stats.MaxDexBonus,
				}
				Armors[key] = armor
				Armors[strings.ToLower(originalName)] = armor
			}
		case "weapon":
			weapon := models.Weapon{Name: strings.ToLower(originalName)}
			Weapons[key] = weapon
			Weapons[strings.ToLower(originalName)] = weapon
		}
	}

	return nil
}

// ------------------------
// Weapon functions
// ------------------------
func AddWeapon(characterName string, newWeapon models.Weapon) (string, error) {
	return AddWeaponToSlot(characterName, newWeapon, "")
}

func AddWeaponToSlot(characterName string, newWeapon models.Weapon, slot string) (string, error) {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return "", fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return "", fmt.Errorf("character '%s' not found", characterName)
	}

	newWeapon.Name = strings.ToLower(strings.TrimSpace(newWeapon.Name)) // lowercase
	var hand string

	switch slot {
	case "":
		if character.Equipment.MainHand == nil {
			character.Equipment.MainHand = &newWeapon
			hand = "main hand"
		} else if character.Equipment.OffHand == nil {
			character.Equipment.OffHand = &newWeapon
			hand = "off hand"
		} else {
			return "", fmt.Errorf("both hands already occupied")
		}
	case "main hand":
		if character.Equipment.MainHand != nil {
			return "", fmt.Errorf("main hand already occupied")
		}
		character.Equipment.MainHand = &newWeapon
		hand = "main hand"
	case "off hand":
		if character.Equipment.OffHand != nil {
			return "", fmt.Errorf("off hand already occupied")
		}
		character.Equipment.OffHand = &newWeapon
		hand = "off hand"
	default:
		return "", fmt.Errorf("invalid slot: must be 'main hand' or 'off hand'")
	}

	if err := storage.SaveCharacter(character); err != nil {
		return "", fmt.Errorf("could not save character: %w", err)
	}

	return hand, nil
}

func RemoveWeapon(characterName, weaponName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	weaponName = normalizeName(weaponName)
	removed := false
	if character.Equipment.MainHand != nil && normalizeName(character.Equipment.MainHand.Name) == weaponName {
		character.Equipment.MainHand = nil
		removed = true
	}
	if character.Equipment.OffHand != nil && normalizeName(character.Equipment.OffHand.Name) == weaponName {
		character.Equipment.OffHand = nil
		removed = true
	}

	if !removed {
		return fmt.Errorf("weapon '%s' not found on character '%s'", weaponName, characterName)
	}

	if err := storage.SaveCharacter(character); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	return nil
}

// ------------------------
// Armor & Shield functions
// ------------------------
func AddArmor(characterName, armorName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	key := strings.ToLower(strings.TrimSpace(armorName))
	armor, ok := Armors[key]
	if !ok {
		return fmt.Errorf("armor '%s' not found", armorName)
	}

	displayName := key
	if key == "padded" || key == "leather" || key == "studded leather" || key == "plate" {
		displayName += " armor"
	}

	displayArmor := armor
	displayArmor.Name = displayName
	character.Equipment.Armor = &displayArmor
	character.CalculateCombatStats()

	if err := storage.SaveCharacter(character); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	fmt.Printf("Equipped armor %s\n", displayName)
	return nil
}

func RemoveArmor(characterName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	character.Equipment.Armor = nil
	character.CalculateCombatStats()

	if err := storage.SaveCharacter(character); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	return nil
}

func AddShield(characterName, shieldName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	key := normalizeName(shieldName)
	shield, ok := Shields[key]
	if !ok {
		shield, ok = Shields[strings.ToLower(shieldName)]
		if !ok {
			return fmt.Errorf("shield '%s' not found", shieldName)
		}
	}

	displayShield := shield
	displayShield.Name = key
	character.Equipment.Shield = &displayShield
	character.CalculateCombatStats()

	if err := storage.SaveCharacter(character); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	return nil
}

func RemoveShield(characterName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	character.Equipment.Shield = nil
	character.CalculateCombatStats()

	if err := storage.SaveCharacter(character); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	return nil
}
