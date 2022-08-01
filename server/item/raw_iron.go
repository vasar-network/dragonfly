package item

// RawIron is a raw metal resource obtained from mining iron ore.
type RawIron struct{}

// SmeltInfo ...
func (RawIron) SmeltInfo() SmeltInfo {
	return newOreSmeltInfo(NewStack(IronIngot{}, 1), 0.7)
}

// EncodeItem ...
func (RawIron) EncodeItem() (name string, meta int16) {
	return "minecraft:raw_iron", 0
}
