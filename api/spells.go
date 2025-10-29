package api

import (
	"dnd-character-sheet/domain"
	"fmt"
	"log"
	"math/rand"
	"strings"
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

type SpellResult struct {
	Spell domain.Spell
	Err   error
}

func GetSpellsForClass(className string, slots map[int]int) ([]domain.Spell, error) {
	list, err := fetchSpellList()
	if err != nil {
		return nil, err
	}

	spells := fetchClassSpells(list.Results, className)
	selected := filterSpellsBySlots(spells, slots)
	final := selectRandomSpells(selected, slots)

	return final, nil
}

func fetchSpellList() (APIListResponse, error) {
	var list APIListResponse
	if err := getJSON("https://www.dnd5eapi.co/api/spells", &list); err != nil {
		return APIListResponse{}, err
	}
	return list, nil
}

func fetchClassSpells(resources []APIResource, className string) []domain.Spell {
	classNameLower := strings.ToLower(className)

	results := make(chan SpellResult, len(resources))
	jobs := make(chan APIResource, len(resources))

	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		go func() {
			for res := range jobs {
				results <- fetchSpell(res, classNameLower)
			}
		}()
	}

	for _, res := range resources {
		jobs <- res
	}
	close(jobs)

	var spells []domain.Spell
	for i := 0; i < len(resources); i++ {
		r := <-results
		if r.Err != nil {
			log.Println(r.Err)
			continue
		}
		if r.Spell.Name != "" {
			spells = append(spells, r.Spell)
		}
	}
	return spells
}

func fetchSpell(res APIResource, classNameLower string) SpellResult {
	url := "https://www.dnd5eapi.co" + strings.ToLower(strings.ReplaceAll(res.URL, " ", "-"))
	var spell APISpell
	if err := getJSON(url, &spell); err != nil {
		return SpellResult{Err: fmt.Errorf("error fetching %s: %w", res.Name, err)}
	}

	for _, c := range spell.Classes {
		if strings.ToLower(c.Name) == classNameLower {
			return SpellResult{
				Spell: domain.Spell{
					Name:   spell.Name,
					Level:  spell.Level,
					School: spell.School.Name,
					Range:  spell.Range,
				},
			}
		}
	}

	return SpellResult{}
}

func filterSpellsBySlots(spells []domain.Spell, slots map[int]int) []domain.Spell {
	var selected []domain.Spell
	for _, s := range spells {
		if _, ok := slots[s.Level]; ok {
			selected = append(selected, s)
		}
	}
	return selected
}

func selectRandomSpells(spells []domain.Spell, slots map[int]int) []domain.Spell {
	rand.Seed(rand.Int63())
	var final []domain.Spell
	for lvl, count := range slots {
		var lvlSpells []domain.Spell
		for _, s := range spells {
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
	return final
}
