package data

type Proxy struct {
	Name        string   `json:"name";yaml:"name"`
	Port        string   `json:"port";yaml:"port"`
	To          []string `json:"to";yaml:"to"`
	Loadbalacer string   `json:"loadbalacer";yaml:"loadbalacer"`
}

type Router struct {
	Port string `json:"port";yaml:"port"`
	Case []struct {
		Service     string   `json:"service,omitempty";yaml:"service"`
		Loadbalacer string   `json:"loadbalacer";yaml:"loadbalacer"`
		Upstream    []string `json:"upstream";yaml:"upstream"`
		Servie      string   `json:"servie,omitempty";yaml:"servie"`
	} `json:"case";yaml:"case"`
}
