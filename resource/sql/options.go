package sql

type GroupConfig struct {
	Name      string   `yaml:"name"`
	Master    string   `yaml:"master"`
	Slaves    []string `yaml:"slaves"`
	StatLevel string   `yaml:"stat_level"`
	LogFormat string   `yaml:"log_format"`
	LogLevel  string   `yaml:"log_level"`
}
