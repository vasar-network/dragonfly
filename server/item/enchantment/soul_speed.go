package enchantment

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// SoulSpeed is an enchantment that can be applied on boots and allows the player to walk more quickly on soul sand or
// soul soil.
type SoulSpeed struct{}

// Name ...
func (SoulSpeed) Name() string {
	return "Soul Speed"
}

// MaxLevel ...
func (SoulSpeed) MaxLevel() int {
	return 3
}

// Cost ...
func (SoulSpeed) Cost(level int) (int, int) {
	min := level * 10
	return min, min + 15
}

// Rarity ...
func (SoulSpeed) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

// Treasure ...
func (SoulSpeed) Treasure() bool {
	return true
}

// CompatibleWithEnchantment ...
func (SoulSpeed) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

// CompatibleWithItem ...
func (SoulSpeed) CompatibleWithItem(i world.Item) bool {
	b, ok := i.(item.BootsType)
	return ok && b.Boots()
}
