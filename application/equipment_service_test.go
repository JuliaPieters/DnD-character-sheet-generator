package application

import (
	"dnd-character-sheet/domain"
	"testing"
)

func TestCalculateWeaponDamage(t *testing.T) {
	service := &EquipmentService{}

	tests := []struct {
		name      string
		character *domain.Character
		weapon    *domain.Weapon
		expected  string
	}{
		{
			name: "Level 1 Half-Orc Barbarian with Greataxe",
			character: &domain.Character{
				Name:        "Barb1",
				StrengthMod: 3,
				DexterityMod: 2,
			},
			weapon: &domain.Weapon{
				Name:      "Greataxe",
				DamageDie: "1d12",
				IsFinesse: false,
			},
			expected: "1d12 + 3",
		},
		{
			name: "Level 1 Tiefling Ranger with Shortsword",
			character: &domain.Character{
				Name:        "Ranger1",
				StrengthMod: 1,
				DexterityMod: 2,
			},
			weapon: &domain.Weapon{
				Name:      "Shortsword",
				DamageDie: "1d6",
				IsFinesse: true,
			},
			expected: "1d6 + 2",
		},
		{
			name: "Level 1 Dwarf Rogue with Rapier",
			character: &domain.Character{
				Name:        "Rogue1",
				StrengthMod: -1,
				DexterityMod: 2,
			},
			weapon: &domain.Weapon{
				Name:      "Rapier",
				DamageDie: "1d8",
				IsFinesse: true,
			},
			expected: "1d8 + 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			damage := service.CalculateWeaponDamage(tt.character, tt.weapon)
			if damage != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, damage)
			}
		})
	}
}
func TestCalculateWeaponDamage_EdgeCases(t *testing.T) {
	service := &EquipmentService{}

	tests := []struct {
		name      string
		character *domain.Character
		weapon    *domain.Weapon
		expected  string
	}{
		{
			name: "Negative STR and DEX",
			character: &domain.Character{
				Name:        "Weakling",
				StrengthMod: -2,
				DexterityMod: -1,
			},
			weapon: &domain.Weapon{
				Name:      "Dagger",
				DamageDie: "1d4",
				IsFinesse: true,
			},
			expected: "1d4 - 1",
		},
		{
			name: "Zero modifiers",
			character: &domain.Character{
				Name:        "Neutral",
				StrengthMod: 0,
				DexterityMod: 0,
			},
			weapon: &domain.Weapon{
				Name:      "Club",
				DamageDie: "1d6",
				IsFinesse: false,
			},
			expected: "1d6 + 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			damage := service.CalculateWeaponDamage(tt.character, tt.weapon)
			if damage != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, damage)
			}
		})
	}
}
