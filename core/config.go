package core

type Config struct {
	PrivateKeyFile string   `yaml:"privateKeyFile"`
	WorkspacePath  string   `yaml:"workspacePath"`
	GitUsername    string   `yaml:"gitUsername"`
	Category       Category `yaml:"repos"`
}

type Category []struct {
	Group string   `yaml:"group"`
	Repos []string `yaml:"repos"`
}
