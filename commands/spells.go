package commands
import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
)

// VoegSpellToe voegt een spell toe aan een caster character
func VoegSpellToe(characterNaam string, nieuweSpell models.Spell) error {
	characters, laadFout := storage.LoadCharacters()
	if laadFout != nil {
		return fmt.Errorf("kon characters niet laden: %w", laadFout)
	}

	character, bestaat := characters[characterNaam]
	if !bestaat {
		return fmt.Errorf("character '%s' niet gevonden", characterNaam)
	}

	if _, isCaster := models.SpellcastingClasses[character.Class]; !isCaster {
		return fmt.Errorf("character '%s' is geen spellcaster", characterNaam)
	}

	character.Spells = append(character.Spells, nieuweSpell)

	if slaFout := storage.SaveCharacter(character); slaFout != nil {
		return fmt.Errorf("kon character niet opslaan: %w", slaFout)
	}

	fmt.Printf("✅ Spell '%s' toegevoegd aan character '%s'\n", nieuweSpell.Name, characterNaam)
	return nil
}

// VerwijderSpell verwijdert een spell van een caster character
func VerwijderSpell(characterNaam string, naamVanSpell string) error {
	characters, laadFout := storage.LoadCharacters()
	if laadFout != nil {
		return fmt.Errorf("kon characters niet laden: %w", laadFout)
	}

	character, bestaat := characters[characterNaam]
	if !bestaat {
		return fmt.Errorf("character '%s' niet gevonden", characterNaam)
	}

	if _, isCaster := models.SpellcastingClasses[character.Class]; !isCaster {
		return fmt.Errorf("character '%s' is geen spellcaster", characterNaam)
	}

	nieuweLijstVanSpells := []models.Spell{}
	for _, bestaandeSpell := range character.Spells {
		if bestaandeSpell.Name != naamVanSpell {
			nieuweLijstVanSpells = append(nieuweLijstVanSpells, bestaandeSpell)
		}
	}
	character.Spells = nieuweLijstVanSpells

	if slaFout := storage.SaveCharacter(character); slaFout != nil {
		return fmt.Errorf("kon character niet opslaan: %w", slaFout)
	}

	fmt.Printf("✅ Spell '%s' verwijderd van character '%s'\n", naamVanSpell, characterNaam)
	return nil
}
