package config

// Config represents the complete configuration for the generator.
type Config struct {
	Schema               SchemaConfig                 `yaml:"schema"`
	Output               OutputConfig                 `yaml:"output"`
	Java                 JavaConfig                   `yaml:"java"`
	TypeMappings         TypeMappingsConfig           `yaml:"typeMappings"`
	Features             FeaturesConfig               `yaml:"features"`
	JavaVersionOverrides map[int]JavaVersionOverrides `yaml:"javaVersionOverrides"`
}

// SchemaConfig contains schema-related configuration.
type SchemaConfig struct {
	Path     string   `yaml:"path"`
	Includes []string `yaml:"includes"`
}

// OutputConfig contains output-related configuration.
type OutputConfig struct {
	Directory string `yaml:"directory"`
	Package   string `yaml:"package"`
}

// JavaConfig contains Java generation settings.
type JavaConfig struct {
	Version          int          `yaml:"version"`
	FieldVisibility  string       `yaml:"fieldVisibility"`
	CollectionType   string       `yaml:"collectionType"`
	NullableHandling string       `yaml:"nullableHandling"`
	Naming           NamingConfig `yaml:"naming"`
}

// NamingConfig contains naming convention settings.
type NamingConfig struct {
	FieldCase       string `yaml:"fieldCase"`
	ClassSuffix     string `yaml:"classSuffix"`
	InterfacePrefix string `yaml:"interfacePrefix"`
}

// TypeMappingsConfig contains type mapping configuration.
type TypeMappingsConfig struct {
	Scalars map[string]ScalarMapping `yaml:"scalars"`
}

// ScalarMapping defines how a GraphQL scalar maps to Java.
type ScalarMapping struct {
	JavaType string   `yaml:"javaType"`
	Imports  []string `yaml:"imports"`
}

// FeaturesConfig contains feature toggle configuration.
type FeaturesConfig struct {
	Lombok     LombokConfig     `yaml:"lombok"`
	Validation ValidationConfig `yaml:"validation"`
	Jackson    JacksonConfig    `yaml:"jackson"`
}

// LombokConfig contains Lombok-related settings.
type LombokConfig struct {
	Enabled           bool `yaml:"enabled"`
	Data              bool `yaml:"data"`
	Builder           bool `yaml:"builder"`
	NoArgsConstructor bool `yaml:"noArgsConstructor"`
	AllArgsConstructor bool `yaml:"allArgsConstructor"`
	Getter            bool `yaml:"getter"`
	Setter            bool `yaml:"setter"`
}

// ValidationConfig contains JSR-303 validation settings.
type ValidationConfig struct {
	Enabled        bool   `yaml:"enabled"`
	Package        string `yaml:"package"`
	NotNullOnNonNull bool `yaml:"notNullOnNonNull"`
}

// JacksonConfig contains Jackson serialization settings.
type JacksonConfig struct {
	Enabled bool `yaml:"enabled"`
}

// JavaVersionOverrides contains overrides for specific Java versions.
type JavaVersionOverrides struct {
	Features FeaturesConfig `yaml:"features"`
}

// FieldVisibility constants.
const (
	VisibilityPrivate   = "private"
	VisibilityProtected = "protected"
	VisibilityPackage   = "package"
	VisibilityPublic    = "public"
)

// CollectionType constants.
const (
	CollectionList       = "List"
	CollectionSet        = "Set"
	CollectionCollection = "Collection"
)

// NullableHandling constants.
const (
	NullableWrapper    = "wrapper"
	NullableOptional   = "optional"
	NullableAnnotation = "annotation"
)

// FieldCase constants.
const (
	FieldCaseCamel = "camelCase"
	FieldCaseSnake = "snake_case"
)

// ValidationPackage constants.
const (
	ValidationJakarta = "jakarta"
	ValidationJavax   = "javax"
)
