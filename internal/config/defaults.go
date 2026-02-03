package config

// DefaultConfig returns a new Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Schema: SchemaConfig{
			Path:     "",
			Includes: []string{},
		},
		Output: OutputConfig{
			Directory: "./generated",
			Package:   "com.example.model",
		},
		Java: JavaConfig{
			Version:          17,
			FieldVisibility:  VisibilityPrivate,
			CollectionType:   CollectionList,
			NullableHandling: NullableWrapper,
			Naming: NamingConfig{
				FieldCase:       FieldCaseCamel,
				ClassSuffix:     "",
				InterfacePrefix: "",
			},
		},
		TypeMappings: TypeMappingsConfig{
			Scalars: map[string]ScalarMapping{},
		},
		Features: FeaturesConfig{
			Lombok: LombokConfig{
				Enabled:           false,
				Data:              true,
				Builder:           false,
				NoArgsConstructor: true,
				AllArgsConstructor: false,
				Getter:            false,
				Setter:            false,
			},
			Validation: ValidationConfig{
				Enabled:        false,
				Package:        ValidationJakarta,
				NotNullOnNonNull: true,
			},
			Jackson: JacksonConfig{
				Enabled: false,
			},
		},
		JavaVersionOverrides: map[int]JavaVersionOverrides{
			8: {
				Features: FeaturesConfig{
					Validation: ValidationConfig{
						Package: ValidationJavax,
					},
				},
			},
		},
	}
}

// DefaultScalarMappings returns commonly used scalar type mappings.
func DefaultScalarMappings() map[string]ScalarMapping {
	return map[string]ScalarMapping{
		"DateTime": {
			JavaType: "java.time.LocalDateTime",
			Imports:  []string{"java.time.LocalDateTime"},
		},
		"Date": {
			JavaType: "java.time.LocalDate",
			Imports:  []string{"java.time.LocalDate"},
		},
		"Time": {
			JavaType: "java.time.LocalTime",
			Imports:  []string{"java.time.LocalTime"},
		},
		"UUID": {
			JavaType: "java.util.UUID",
			Imports:  []string{"java.util.UUID"},
		},
		"BigDecimal": {
			JavaType: "java.math.BigDecimal",
			Imports:  []string{"java.math.BigDecimal"},
		},
		"BigInteger": {
			JavaType: "java.math.BigInteger",
			Imports:  []string{"java.math.BigInteger"},
		},
		"Long": {
			JavaType: "Long",
			Imports:  []string{},
		},
		"JSON": {
			JavaType: "com.fasterxml.jackson.databind.JsonNode",
			Imports:  []string{"com.fasterxml.jackson.databind.JsonNode"},
		},
	}
}

// SupportedJavaVersions returns the list of supported Java versions.
func SupportedJavaVersions() []int {
	return []int{8, 11, 17, 21}
}
