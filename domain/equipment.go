package domain

type Weapon struct {
	Name      string
	Category  string
	Range     string
	TwoHanded bool
}

type Armor struct {
	Name        string
	ArmorClass  int
	DexBonus    bool
	MaxDexBonus int
}

type Shield struct {
	Name       string
	ArmorClass int
}

type Equipment struct {
	MainHand *Weapon
	OffHand  *Weapon
	Armor    *Armor
	Shield   *Shield
}
