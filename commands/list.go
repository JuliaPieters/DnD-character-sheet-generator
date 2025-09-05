package commands

import (
    "dnd-character-sheet/storage"
    "fmt"
)

// ListCharacters toont alle characters met level en class
func ListCharacters() error {
    chars, err := storage.LoadCharacters()
    if err != nil {
        return err
    }

    if len(chars) == 0 {
        fmt.Println("📜 No characters found.")
        return nil
    }

    fmt.Println("📜 Characters:")
    for _, c := range chars {
        fmt.Printf(" - %s | Level %d %s %s\n", c.Name, c.Level, c.Race, c.Class)
    }
    return nil
}
