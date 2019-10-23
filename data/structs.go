package data

type Proxy struct {
	Name        string   `json:"name";yaml:"name"`
	Port        string   `json:"port";yaml:"port"`
	To          []string `json:"to";yaml:"to"`
	Loadbalacer string   `json:"loadbalacer";yaml:"loadbalacer"`
	Static      string   `json:"static,omitempty";yaml:"static,omitempty"`
}

type Router struct {
	Port string `json:"port";yaml:"port"`
	Case []struct {
		Service     string   `json:"service,omitempty";yaml:"service"`
		Loadbalacer string   `json:"loadbalacer";yaml:"loadbalacer"`
		Upstream    []string `json:"upstream";yaml:"upstream"`
		Static      []string `json:"static,omitempty";yaml:"static,omitempty"`
	} `json:"case";yaml:"case"`
}
