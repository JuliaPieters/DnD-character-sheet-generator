package commands

import (
    "dnd-character-sheet/models"
    "dnd-character-sheet/storage"
    "fmt"
)

// CreateCharacter maakt een nieuw D&D character aan
func CreateCharacter(
    name string,
    playerName string,
    race string,
    class string,
    background string,
    alignment string,
    level int,
    str int,
    dex int,
    con int,
    intt int,
    wis int,
    cha int,
    skills []string,
) error {

    // Maak een nieuw character aan met alle velden uit models
    character := models.NewCharacter(
        name,
        playerName,
        race,
        class,
        background,
        alignment,
        level,
        str,
        dex,
        con,
        intt,
        wis,
        cha,
    )

    // Voeg skills toe en bereken modifiers
    if len(skills) > 0 {
        character.AddSkills(skills)
    }

    // Sla het character op in storage
    if err := storage.SaveCharacter(*character); err != nil {
        return err
    }

    fmt.Printf("✅ Character %s created!\n", name)
    return nil
}
