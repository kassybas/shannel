package snlapi

type SnlFile struct {
	ApiVersion string             `yaml:"apiversion,omitempty"`
	Target     map[string]Target  `yaml:"targets"`
	Vars       map[string]*string `yaml:"vars"`
	Options    map[string]string  `yaml:"options"`
	Args       []Arg              `yaml:"args"`
}

type Arg struct {
	Name       string  `yaml:"name"`
	Type       string  `yaml:"type"`
	Usage      string  `yaml:"usage"`
	Alias      string  `yaml:"alias"`
	Default    *string `yaml:"default"`
	FromEnvVar *string `yaml:"fromEnvVar"`
}

type Target struct {
	Sh      string `yaml:"sh"`
	Usage   string `yaml:"desc,omitempty"`
	Timeout string `yaml:"timeout,omitempty"`
}
