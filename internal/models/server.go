package models

type Server struct {
	Name    string                 `yaml:"name" json:"name"`
	Host    string                 `yaml:"host" json:"host"`
	User    string                 `yaml:"user,omitempty" json:"user,omitempty"`
	Port    int                    `yaml:"port,omitempty" json:"port,omitempty"`
	Tags    []string               `yaml:"tags,omitempty" json:"tags,omitempty"`
	Key     string                 `yaml:"key,omitempty" json:"key,omitempty"`
	Options map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`
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
