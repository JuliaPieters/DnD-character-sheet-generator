package commands

import (
	"dnd-character-sheet/storage"
	"fmt"
)

// DeleteCharacter verwijdert een character en geeft een nette bevestiging
func DeleteCharacter(name string) error {
	chars, err := storage.LoadCharacters()
	if err != nil {
		return err
	}

	if len(chars) == 0 {
		fmt.Println("📜 No characters found to delete.")
		return nil
	}

	c, ok := chars[name]
	if !ok {
		return fmt.Errorf("character %s not found", name)
	}

	delete(chars, name)

	if err := storage.SaveAllCharacters(chars); err != nil {
		return err
	}

	fmt.Printf("🗑️ Character %s (%s %s) deleted.\n", name, c.Race, c.Class)
	return nil
}
