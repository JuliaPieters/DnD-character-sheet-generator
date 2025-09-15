package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
)

// VoegWapenToe voegt een wapen toe aan een character
func VoegWapenToe(characterNaam string, nieuwWapen models.Weapon) error {
	characters, laadFout := storage.LoadCharacters()
	if laadFout != nil {
		return fmt.Errorf("kon characters niet laden: %w", laadFout)
	}

	character, bestaat := characters[characterNaam]
	if !bestaat {
		return fmt.Errorf("character '%s' niet gevonden", characterNaam)
	}

	character.Equipment.Weapons = append(character.Equipment.Weapons, nieuwWapen)

	if slaFout := storage.SaveCharacter(character); slaFout != nil {
		return fmt.Errorf("kon character niet opslaan: %w", slaFout)
	}

	fmt.Printf("✅ Wapen '%s' toegevoegd aan character '%s'\n", nieuwWapen.Name, characterNaam)
	return nil
}

// VerwijderWapen verwijdert een wapen van een character
func VerwijderWapen(characterNaam string, naamVanWapen string) error {
	characters, laadFout := storage.LoadCharacters()
	if laadFout != nil {
		return fmt.Errorf("kon characters niet laden: %w", laadFout)
	}

	character, bestaat := characters[characterNaam]
	if !bestaat {
		return fmt.Errorf("character '%s' niet gevonden", characterNaam)
	}

	nieuweLijstVanWapens := []models.Weapon{}
	for _, bestaandWapen := range character.Equipment.Weapons {
		if bestaandWapen.Name != naamVanWapen {
			nieuweLijstVanWapens = append(nieuweLijstVanWapens, bestaandWapen)
		}
	}
	character.Equipment.Weapons = nieuweLijstVanWapens

	if slaFout := storage.SaveCharacter(character); slaFout != nil {
		return fmt.Errorf("kon character niet opslaan: %w", slaFout)
	}

	fmt.Printf("✅ Wapen '%s' verwijderd van character '%s'\n", naamVanWapen, characterNaam)
	return nil
}

// VoegArmorToe voegt armor toe en berekent Armor Class opnieuw
func VoegArmorToe(characterNaam string, nieuwArmor models.Armor) error {
	characters, laadFout := storage.LoadCharacters()
	if laadFout != nil {
		return fmt.Errorf("kon characters niet laden: %w", laadFout)
	}

	character, bestaat := characters[characterNaam]
	if !bestaat {
		return fmt.Errorf("character '%s' niet gevonden", characterNaam)
	}

	character.Equipment.Armor = &nieuwArmor
	character.CalculateCombatStats()

	if slaFout := storage.SaveCharacter(character); slaFout != nil {
		return fmt.Errorf("kon character niet opslaan: %w", slaFout)
	}

	fmt.Printf("✅ Armor '%s' toegevoegd aan character '%s'\n", nieuwArmor.Name, characterNaam)
	return nil
}

// VerwijderArmor verwijdert armor en herberekent Armor Class
func VerwijderArmor(characterNaam string) error {
	characters, laadFout := storage.LoadCharacters()
	if laadFout != nil {
		return fmt.Errorf("kon characters niet laden: %w", laadFout)
	}

	character, bestaat := characters[characterNaam]
	if !bestaat {
		return fmt.Errorf("character '%s' niet gevonden", characterNaam)
	}

	character.Equipment.Armor = nil
	character.CalculateCombatStats()

	if slaFout := storage.SaveCharacter(character); slaFout != nil {
		return fmt.Errorf("kon character niet opslaan: %w", slaFout)
	}

	fmt.Printf("✅ Armor verwijderd van character '%s'\n", characterNaam)
	return nil
}

// VoegShieldToe voegt een shield toe en herberekent Armor Class
func VoegShieldToe(characterNaam string, nieuwShield models.Shield) error {
	characters, laadFout := storage.LoadCharacters()
	if laadFout != nil {
		return fmt.Errorf("kon characters niet laden: %w", laadFout)
	}

	character, bestaat := characters[characterNaam]
	if !bestaat {
		return fmt.Errorf("character '%s' niet gevonden", characterNaam)
	}

	character.Equipment.Shield = &nieuwShield
	character.CalculateCombatStats()

	if slaFout := storage.SaveCharacter(character); slaFout != nil {
		return fmt.Errorf("kon character niet opslaan: %w", slaFout)
	}

	fmt.Printf("✅ Shield '%s' toegevoegd aan character '%s'\n", nieuwShield.Name, characterNaam)
	return nil
}

// VerwijderShield verwijdert het shield en herberekent Armor Class
func VerwijderShield(characterNaam string) error {
	characters, laadFout := storage.LoadCharacters()
	if laadFout != nil {
		return fmt.Errorf("kon characters niet laden: %w", laadFout)
	}

	character, bestaat := characters[characterNaam]
	if !bestaat {
		return fmt.Errorf("character '%s' niet gevonden", characterNaam)
	}

	character.Equipment.Shield = nil
	character.CalculateCombatStats()

	if slaFout := storage.SaveCharacter(character); slaFout != nil {
		return fmt.Errorf("kon character niet opslaan: %w", slaFout)
	}

	fmt.Printf("✅ Shield verwijderd van character '%s'\n", characterNaam)
	return nil
}
