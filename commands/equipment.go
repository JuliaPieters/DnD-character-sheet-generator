package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
)

// AddWeapon equips a weapon to the first available hand (main or off) and returns which hand it was equipped to
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

// RemoveWeapon removes a weapon from main hand or off hand
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

// AddArmor equips armor
func AddArmor(characterName string, newArmor models.Armor) error {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	character.Equipment.Armor = &newArmor
	character.CalculateCombatStats()

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	return nil
}

// RemoveArmor removes armor
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

// AddShield equips shield
func AddShield(characterName string, newShield models.Shield) error {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	character.Equipment.Shield = &newShield
	character.CalculateCombatStats()

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	return nil
}

// RemoveShield removes shield
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
