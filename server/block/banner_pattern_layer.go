package block

import (
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
)

// BannerPatternLayer is a wrapper over BannerPatternType with a colour property.
type BannerPatternLayer struct {
	BannerPatternType
	// Colour is the colour the pattern should be rendered in.
	Colour item.Colour
}

// EncodeNBT encodes the given BannerPatternLayer into an NBT map.
func (b BannerPatternLayer) EncodeNBT() map[string]any {
	return map[string]any{
		"Pattern": BannerPatternID(b.BannerPatternType),
		"Color":   int32(invertColour(b.Colour)),
	}
}

// DecodeNBT decodes the given NBT map into a BannerPatternLayer and returns it.
func (b BannerPatternLayer) DecodeNBT(data map[string]any) any {
	b.BannerPatternType, _ = BannerPatternByID(nbtconv.Map[string](data, "Pattern"))
	b.Colour = invertColourID(int16(nbtconv.Map[int32](data, "Color")))
	return b
}