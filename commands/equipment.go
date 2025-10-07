package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
)

// ------------------------
// Standaard Armors & Shields
// ------------------------
var Armors = map[string]models.Armor{
	// Light
	"padded":          {Name: "padded", ArmorClass: 11, DexBonus: true},
	"leather armor":   {Name: "leather armor", ArmorClass: 11, DexBonus: true},
	"studded leather": {Name: "studded leather", ArmorClass: 12, DexBonus: true},

	// Medium
	"hide":        {Name: "hide", ArmorClass: 12, DexBonus: true},
	"chain shirt": {Name: "chain shirt", ArmorClass: 13, DexBonus: true},
	"scale mail":  {Name: "scale mail", ArmorClass: 14, DexBonus: true},
	"breastplate": {Name: "breastplate", ArmorClass: 14, DexBonus: true},
	"half plate":  {Name: "half plate", ArmorClass: 15, DexBonus: true},

	// Heavy
	"ring mail":   {Name: "ring mail", ArmorClass: 14, DexBonus: false},
	"chain mail":  {Name: "chain mail", ArmorClass: 16, DexBonus: false},
	"splint":      {Name: "splint", ArmorClass: 17, DexBonus: false},
	"plate armor": {Name: "plate armor", ArmorClass: 18, DexBonus: false},
}

var Shields = map[string]models.Shield{
	"shield": {Name: "shield", ArmorClass: 2},
}

// ------------------------
// Weapon functions
// ------------------------
func AddWeapon(characterName string, newWeapon models.Weapon) (string, error) {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return "", fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return "", fmt.Errorf("character '%s' not found", characterName)
	}

	var hand string
	if character.Equipment.MainHand == nil {
		character.Equipment.MainHand = &newWeapon
		hand = "main hand"
	} else if character.Equipment.OffHand == nil {
		character.Equipment.OffHand = &newWeapon
		hand = "off hand"
	} else {
		return "", fmt.Errorf("both hands already occupied")
	}

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return "", fmt.Errorf("could not save character: %w", saveErr)
	}

	return hand, nil
}

func AddWeaponToSlot(characterName string, newWeapon models.Weapon, slot string) (string, error) {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return "", fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return "", fmt.Errorf("character '%s' not found", characterName)
	}

	var hand string
	switch slot {
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

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return "", fmt.Errorf("could not save character: %w", saveErr)
	}

	return hand, nil
}

func RemoveWeapon(characterName string, weaponName string) error {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	removed := false
	if character.Equipment.MainHand != nil && character.Equipment.MainHand.Name == weaponName {
		character.Equipment.MainHand = nil
		removed = true
	}
	if character.Equipment.OffHand != nil && character.Equipment.OffHand.Name == weaponName {
		character.Equipment.OffHand = nil
		removed = true
	}

	if !removed {
		return fmt.Errorf("weapon '%s' not found on character '%s'", weaponName, characterName)
	}

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	return nil
}

// ------------------------
// Armor functions
// ------------------------
func AddArmor(characterName string, armorName string) error {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	armor, ok := Armors[armorName]
	if !ok {
		return fmt.Errorf("armor '%s' not found", armorName)
	}

	character.Equipment.Armor = &armor
	character.CalculateCombatStats()

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	return nil
}

func RemoveArmor(characterName string) error {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	character.Equipment.Armor = nil
	character.CalculateCombatStats()

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	return nil
}

// ------------------------
// Shield functions
// ------------------------
func AddShield(characterName string, shieldName string) error {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	shield, ok := Shields[shieldName]
	if !ok {
		return fmt.Errorf("shield '%s' not found", shieldName)
	}

	character.Equipment.Shield = &shield
	character.CalculateCombatStats()

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	return nil
}

func RemoveShield(characterName string) error {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	character.Equipment.Shield = nil
	character.CalculateCombatStats()

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	return nil
}
