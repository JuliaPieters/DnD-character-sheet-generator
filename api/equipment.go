package api

import (
	"dnd-character-sheet/domain"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type EquipmentRange struct {
	Normal int `json:"normal,omitempty"`
	Long   int `json:"long,omitempty"`
}

type APIEquipment struct {
	Name              string `json:"name"`
	EquipmentCategory struct {
		Name string `json:"name"`
	} `json:"equipment_category"`
	WeaponCategory string `json:"weapon_category"` 
	WeaponRange    string `json:"weapon_range"`    
	ArmorClass struct {
		Base     int  `json:"base"`
		DexBonus bool `json:"dex_bonus"`
		MaxDex   int  `json:"max_bonus"`
	} `json:"armor_class,omitempty"`
	TwoHanded  bool            `json:"two_handed,omitempty"`
	Range      json.RawMessage `json:"range,omitempty"`
	Properties []struct {
		Name string `json:"name"`
	} `json:"properties,omitempty"`
	Damage struct {
		DamageDice string `json:"damage_dice"`
	} `json:"damage,omitempty"`
}

func parseRange(raw json.RawMessage) string {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var r EquipmentRange
	if err := json.Unmarshal(raw, &r); err == nil {
		if r.Normal > 0 {
			return strconv.Itoa(r.Normal)
		}
	}
	return ""
}

func normalizeDamageDie(d string) string {
	if d == "" {
		return "1d4"
	}
	if strings.Contains(d, "d") {
		return d
	}
	num, err := strconv.Atoi(d)
	if err != nil || num <= 0 {
		return "1d4"
	}
	return fmt.Sprintf("%dd%d", num, num)
}

func GetAllEquipment() ([]*domain.Weapon, *domain.Armor, *domain.Shield, error) {
	var list APIListResponse
	if err := getJSON("https://www.dnd5eapi.co/api/equipment", &list); err != nil {
		return nil, nil, nil, err
	}

	var weapons []*domain.Weapon
	var armor *domain.Armor
	var shield *domain.Shield

	type result struct {
		eq  APIEquipment
		err error
	}

	jobs := make(chan APIResource, len(list.Results))
	results := make(chan result, len(list.Results))

	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		go func() {
			for res := range jobs {
				var eq APIEquipment
				url := "https://www.dnd5eapi.co" + strings.ToLower(strings.ReplaceAll(res.URL, " ", "-"))
				err := getJSON(url, &eq)
				results <- result{eq: eq, err: err}
			}
		}()
	}

	for _, res := range list.Results {
		jobs <- res
	}
	close(jobs)

	for i := 0; i < len(list.Results); i++ {
		r := <-results
		if r.err != nil {
			log.Println("Error fetching equipment:", r.err)
			continue
		}
		eq := r.eq
		switch eq.EquipmentCategory.Name {
		case "Weapon":
			damage := normalizeDamageDie(eq.Damage.DamageDice)
			isFinesse := false
			for _, p := range eq.Properties {
				if strings.ToLower(p.Name) == "finesse" {
					isFinesse = true
				}
			}
			weapon := &domain.Weapon{
				Name:      eq.Name,
				TwoHanded: eq.TwoHanded,
				Range:     parseRange(eq.Range),
				DamageDie: damage,
				IsFinesse: isFinesse,
				Category:  strings.TrimSpace(eq.WeaponCategory + " " + eq.WeaponRange),
			}
			weapons = append(weapons, weapon)
		case "Armor":
			armor = &domain.Armor{
				Name:        eq.Name,
				ArmorClass:  eq.ArmorClass.Base,
				DexBonus:    eq.ArmorClass.DexBonus,
				MaxDexBonus: eq.ArmorClass.MaxDex,
			}
		case "Shield":
			shield = &domain.Shield{
				Name:       eq.Name,
				ArmorClass: eq.ArmorClass.Base,
			}
		}
	}

	return weapons, armor, shield, nil
}
