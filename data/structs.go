package data

// Proxy is data format to parse the configuration file
type Proxy struct {
	Name        string   `json:"name";yaml:"name"`
	Port        string   `json:"port";yaml:"port"`
	To          []string `json:"to";yaml:"to"`
	Loadbalacer string   `json:"loadbalacer";yaml:"loadbalacer"`
	Static      string   `json:"static,omitempty";yaml:"static,omitempty"`
}

// Router is data format to parse the configuration file
type Router struct {
	Port string `json:"port";yaml:"port"`
	Case []struct {
		Service     string   `json:"service,omitempty";yaml:"service"`
		Loadbalacer string   `json:"loadbalacer";yaml:"loadbalacer"`
		Upstream    []string `json:"upstream";yaml:"upstream"`
		Static      []string `json:"static,omitempty";yaml:"static,omitempty"`
	} `json:"case";yaml:"case"`
}
