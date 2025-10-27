package api

import (
	"dnd-character-sheet/domain"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

type APISpell struct {
	Name   string `json:"name"`
	Level  int    `json:"level"`
	School struct {
		Name string `json:"name"`
	} `json:"school"`
	Range   string `json:"range"`
	Classes []struct {
		Name string `json:"name"`
	} `json:"classes"`
}

func GetSpellsForClass(className string, slots map[int]int) ([]domain.Spell, error) {
	var list APIListResponse
	if err := getJSON("https://www.dnd5eapi.co/api/spells", &list); err != nil {
		return nil, err
	}

	type SpellResult struct {
		Spell domain.Spell
		Err   error
	}

	results := make(chan SpellResult)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	classNameLower := strings.ToLower(className)

	for _, res := range list.Results {
		<-ticker.C
		go func(res APIResource) {
			url := "https://www.dnd5eapi.co" + strings.ToLower(strings.ReplaceAll(res.URL, " ", "-"))
			var spell APISpell
			err := getJSON(url, &spell)
			if err != nil {
				results <- SpellResult{Err: fmt.Errorf("error fetching %s: %w", res.Name, err)}
				return
			}

			isForClass := false
			for _, c := range spell.Classes {
				if strings.ToLower(c.Name) == classNameLower {
					isForClass = true
					break
				}
			}
			if !isForClass {
				results <- SpellResult{Err: nil}
				return
			}

			results <- SpellResult{
				Spell: domain.Spell{
					Name:   spell.Name,
					Level:  spell.Level,
					School: spell.School.Name,
					Range:  spell.Range,
				},
			}
		}(res)
	}

	selected := []domain.Spell{}
	for i := 0; i < len(list.Results); i++ {
		result := <-results
		if result.Err != nil {
			log.Println(result.Err)
			continue
		}
		if result.Spell.Name == "" {
			continue
		}
		if _, ok := slots[result.Spell.Level]; ok {
			selected = append(selected, result.Spell)
		}
	}

	rand.Seed(time.Now().UnixNano())
	final := []domain.Spell{}
	for lvl, count := range slots {
		lvlSpells := []domain.Spell{}
		for _, s := range selected {
			if s.Level == lvl {
				lvlSpells = append(lvlSpells, s)
			}
		}
		if len(lvlSpells) > count {
			rand.Shuffle(len(lvlSpells), func(i, j int) { lvlSpells[i], lvlSpells[j] = lvlSpells[j], lvlSpells[i] })
			final = append(final, lvlSpells[:count]...)
		} else {
			final = append(final, lvlSpells...)
		}
	}

	return final, nil
}
