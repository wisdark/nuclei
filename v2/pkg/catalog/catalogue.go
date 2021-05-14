package catalog

// Catalog is a template catalog helper implementation
type Catalog struct {
	ignoreFiles        []string
	templatesDirectory string
}

// New creates a new Catalog structure using provided input items
func New(directory string) *Catalog {
	catalog := &Catalog{templatesDirectory: directory}
	return catalog
}

// AppendIgnore appends to the catalog store ignore list.
func (c *Catalog) AppendIgnore(list []string) {
	c.ignoreFiles = append(c.ignoreFiles, list...)
}
