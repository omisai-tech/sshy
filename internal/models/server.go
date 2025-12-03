package models

type Server struct {
	Name    string                 `yaml:"name"`
	Host    string                 `yaml:"host"`
	User    string                 `yaml:"user,omitempty"`
	Port    int                    `yaml:"port,omitempty"`
	Tags    []string               `yaml:"tags,omitempty"`
	Key     string                 `yaml:"key,omitempty"`
	Options map[string]interface{} `yaml:"options,omitempty"`
}

type Servers []Server

type ServerSource int

const (
	SourceShared ServerSource = iota
	SourceLocal
	SourceOverride
)

type ServerWithSource struct {
	Server Server
	Source ServerSource
}
