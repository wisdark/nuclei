package catalog

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResolvePath resolves the path to an absolute one in various ways.
//
// It checks if the filename is an absolute path, looks in the current directory
// or checking the nuclei templates directory. If a second path is given,
// it also tries to find paths relative to that second path.
func (c *Catalog) ResolvePath(templateName, second string) (string, error) {
	if filepath.IsAbs(templateName) {
		return templateName, nil
	}

	if second != "" {
		secondBasePath := filepath.Join(filepath.Dir(second), templateName)
		if _, err := os.Stat(secondBasePath); !os.IsNotExist(err) {
			return secondBasePath, nil
		}
	}

	curDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	templatePath := filepath.Join(curDirectory, templateName)
	if _, err := os.Stat(templatePath); !os.IsNotExist(err) {
		return templatePath, nil
	}

	if c.templatesDirectory != "" {
		templatePath := filepath.Join(c.templatesDirectory, templateName)
		if _, err := os.Stat(templatePath); !os.IsNotExist(err) {
			return templatePath, nil
		}
	}
	return "", fmt.Errorf("no such path found: %s", templateName)
}
