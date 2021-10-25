package snlapi

// SnlFile is the config for a shannel file
type SnlFile struct {
	APIVersion string             `yaml:"apiversion,omitempty"`
	Target     map[string]Target  `yaml:"targets"`
	Vars       map[string]*string `yaml:"vars"`
	Options    map[string]string  `yaml:"options"`
	Args       []Arg              `yaml:"args"`
}

// Arg is the config for an argument of a shannel file
type Arg struct {
	Name       string  `yaml:"name"`
	Type       string  `yaml:"type"`
	Usage      string  `yaml:"usage"`
	Alias      string  `yaml:"alias"`
	Default    *string `yaml:"default"`
	FromEnvVar *string `yaml:"fromEnvVar"`
}

// Target is the config of a target to be executed
type Target struct {
	Sh      string `yaml:"sh"`
	Usage   string `yaml:"desc,omitempty"`
	Timeout string `yaml:"timeout,omitempty"`
}
