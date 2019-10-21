package data

type Proxy struct {
	Name        string   `json:"name"`
	Port        string   `json:"port"`
	To          []string `json:"to"`
	Loadbalacer string   `json:"loadbalacer"`
}
