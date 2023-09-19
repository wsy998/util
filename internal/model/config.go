package model

type (
	Frontend struct {
		Path      string            `json:"path"`
		Install   string            `json:"install"`
		Dev       string            `json:"dev"`
		Build     string            `json:"build"`
		BuildPath string            `json:"buildPath"`
		Port      int               `json:"port"`
		Alias     map[string]string `json:"alias"`
	}
	Backend struct {
		Path      string            `json:"path,omitempty"`
		Dev       string            `json:"dev,omitempty"`
		Build     string            `json:"build,omitempty"`
		BuildPath string            `json:"buildPath,omitempty"`
		OutName   string            `json:"outName,omitempty"`
		Alias     map[string]string `json:"alias"`
		Mirror    string            `json:"mirror"`
	}
	TypeConfig struct {
		Backend  Backend           `json:"backend"`
		Frontend Frontend          `json:"frontend"`
		Alias    map[string]string `json:"alias"`
	}
)
