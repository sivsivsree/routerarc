package data

type Proxy struct {
	Name        string   `json:"name"`
	Port        string   `json:"port"`
	To          []string `json:"to"`
	Loadbalacer string   `json:"loadbalacer"`
}

type Router struct {
	Port string `json:"port"`
	Case []struct {
		Service     string   `json:"service,omitempty"`
		Loadbalacer string   `json:"loadbalacer"`
		Upstream    []string `json:"upstream"`
		Servie      string   `json:"servie,omitempty"`
	} `json:"case"`
}
