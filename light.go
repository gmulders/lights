package lights

// Light represents a hue light but simpler
type Light struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
	On   bool   `json:"on"`
	Bri  uint8  `json:"bri"`
	Hue  uint16 `json:"hue"`
	Sat  uint8  `json:"sat"`
}
