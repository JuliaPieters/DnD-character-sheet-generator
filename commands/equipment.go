package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
)

func AddWeapon(characterName string, newWeapon models.Weapon) error {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	character.Equipment.Weapons = append(character.Equipment.Weapons, newWeapon)

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	fmt.Printf("✅ Weapon '%s' added to character '%s'\n", newWeapon.Name, characterName)
	return nil
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

	newWeaponList := []models.Weapon{}
	for _, existingWeapon := range character.Equipment.Weapons {
		if existingWeapon.Name != weaponName {
			newWeaponList = append(newWeaponList, existingWeapon)
		}
	}
	character.Equipment.Weapons = newWeaponList

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	fmt.Printf(" Weapon '%s' removed from character '%s'\n", weaponName, characterName)
	return nil
}

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

	fmt.Printf("✅ Armor '%s' added to character '%s'\n", newArmor.Name, characterName)
	return nil
}

// RemoveArmor removes armor from a character and recalculates AC
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

	fmt.Printf("✅ Armor removed from character '%s'\n", characterName)
	return nil
}

// AddShield adds a shield to a character and recalculates AC
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

	fmt.Printf("✅ Shield '%s' added to character '%s'\n", newShield.Name, characterName)
	return nil
}

// RemoveShield removes a shield from a character and recalculates AC
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

	fmt.Printf("✅ Shield removed from character '%s'\n", characterName)
	return nil
}
