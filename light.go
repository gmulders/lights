package lights

// Light represents a hue light but simpler
type Light struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
	On   bool   `json:"on"`
	Bri  uint8  `json:"bri"` // 1 - 254: 1 really dark, 254 really bright
	Hue  uint16 `json:"hue"` // 0 - 65535: 0 -> red, 21845 -> green, 43690 -> blue, 65535 -> red
	Sat  uint8  `json:"sat"` // 0 - 254: 0 -> white, 254 -> coloured
}
