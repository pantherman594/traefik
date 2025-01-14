package plugins

// Descriptor The static part of a plugin configuration (prod).
type Descriptor struct {
	// ModuleName (required)
	ModuleName string `description:"plugin's module name." json:"moduleName,omitempty" toml:"moduleName,omitempty" yaml:"moduleName,omitempty" export:"true"`

	// Version (required)
	Version string `description:"plugin's version." json:"version,omitempty" toml:"version,omitempty" yaml:"version,omitempty" export:"true"`

	// AllowUnsafe (optional)
	AllowUnsafe bool `description:"allow plugin to use unsafe."  json:"allowUnsafe" toml:"allowUnsafe" yaml:"allowUnsafe" export:"true"`
}

// DevPlugin The static part of a plugin configuration (only for dev).
type DevPlugin struct {
	// GoPath plugin's GOPATH. (required)
	GoPath string `description:"plugin's GOPATH." json:"goPath,omitempty" toml:"goPath,omitempty" yaml:"goPath,omitempty" export:"true"`

	// ModuleName (required)
	ModuleName string `description:"plugin's module name."  json:"moduleName,omitempty" toml:"moduleName,omitempty" yaml:"moduleName,omitempty" export:"true"`

	// AllowUnsafe (optional)
	AllowUnsafe bool `description:"allow plugin to use unsafe."  json:"allowUnsafe" toml:"allowUnsafe" yaml:"allowUnsafe" export:"true"`
}

// Manifest The plugin manifest.
type Manifest struct {
	DisplayName   string                 `yaml:"displayName"`
	Type          string                 `yaml:"type"`
	Import        string                 `yaml:"import"`
	BasePkg       string                 `yaml:"basePkg"`
	Compatibility string                 `yaml:"compatibility"`
	Summary       string                 `yaml:"summary"`
	TestData      map[string]interface{} `yaml:"testData"`
}
