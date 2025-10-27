package api

import (
	"dnd-character-sheet/domain"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
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
	ArmorClass struct {
		Base     int  `json:"base"`
		DexBonus bool `json:"dex_bonus"`
		MaxDex   int  `json:"max_bonus"`
	} `json:"armor_class,omitempty"`
	TwoHanded bool            `json:"two_handed,omitempty"`
	Range     json.RawMessage `json:"range,omitempty"`
}

func parseRange(raw json.RawMessage) string {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}

	var r EquipmentRange
	if err := json.Unmarshal(raw, &r); err == nil {
		if r.Normal > 0 {
			return fmt.Sprintf("%d", r.Normal)
		}
	}
	return ""
}

func GetEquipment() (*domain.Weapon, *domain.Weapon, *domain.Armor, *domain.Shield, error) {
	var list APIListResponse
	if err := getJSON("https://www.dnd5eapi.co/api/equipment", &list); err != nil {
		return nil, nil, nil, nil, err
	}

	var mainHand *domain.Weapon
	var offHand *domain.Weapon
	var armor *domain.Armor
	var shield *domain.Shield

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for _, res := range list.Results {
		<-ticker.C
		var eq APIEquipment
		url := "https://www.dnd5eapi.co" + strings.ToLower(strings.ReplaceAll(res.URL, " ", "-"))
		if err := getJSON(url, &eq); err != nil {
			log.Println("Error fetching equipment:", res.Name, err)
			continue
		}

		switch eq.EquipmentCategory.Name {
		case "Weapon":
			weapon := &domain.Weapon{
				Name:      eq.Name,
				TwoHanded: eq.TwoHanded,
				Range:     parseRange(eq.Range),
			}
			if mainHand == nil {
				mainHand = weapon
			} else if offHand == nil {
				offHand = weapon
			}
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

	return mainHand, offHand, armor, shield, nil
}
