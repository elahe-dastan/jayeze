package config

type Cluster struct {
	Path string `koanf:"path"`
	Size int    `koanf:"size"`
}

type Config struct {
	IndexPath string    `koanf:"indexPath"`
	DocsSize  int       `koanf:"docsNum"`
	clusters  []Cluster `koanf:"clusters"`
}

func Read() *Config {
	
}
