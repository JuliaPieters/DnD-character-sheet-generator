package commands

import (
	"dnd-character-sheet/storage"
	"fmt"
)

func ListCharacters() error {
	allCharacters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}

	if len(allCharacters) == 0 {
		fmt.Println(" No characters found.")
		return nil
	}

	fmt.Println("Characters:")
	for _, character := range allCharacters {
		fmt.Printf("- Name: %s | Level: %d | Race: %s | Class: %s\n",
			character.Name, character.Level, character.Race, character.Class)
	}
	return nil
}
