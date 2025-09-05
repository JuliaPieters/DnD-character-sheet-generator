package commands

import (
    "dnd-character-sheet/models"
    "dnd-character-sheet/storage"
    "fmt"
)

// CreateCharacter maakt een nieuw D&D character aan
func CreateCharacter(
    name, race, class, background string,
    level int,
    str, dex, con, intt, wis, cha int,
    skills []string,
) error {
    // Gebruik de NewCharacter functie uit models
    c := models.NewCharacter(name, race, class, background, level, str, dex, con, intt, wis, cha)

    // Voeg skills toe en bereken modifiers
    if len(skills) > 0 {
        c.AddSkills(skills)
    }

    // Sla het character op in storage
    if err := storage.SaveCharacter(*c); err != nil {
        return err
    }

    fmt.Printf("✅ Character %s created!\n", name)
    return nil
}
