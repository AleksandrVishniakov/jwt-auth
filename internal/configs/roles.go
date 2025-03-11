package configs

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Role struct {
	Permissions []string `yaml:"permissions"`
	Default     bool
	Super bool
}

type yamlStructure struct {
	Roles map[string]Role `yaml:"roles"`
}

func MustParseRoles(path string) map[string]Role {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open %s: %s", path, err.Error())
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read %s: %s", path, err.Error())
	}

	var yamlData yamlStructure
	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		log.Fatalf("Failed to parse yaml %s: %s", path, err.Error())
	}

	return yamlData.Roles
}