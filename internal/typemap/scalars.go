package typemap

// ScalarInfo contains information about a Java scalar type.
type ScalarInfo struct {
	JavaType    string
	PrimitiveType string // Optional primitive version (e.g., "int" for "Integer")
	Imports     []string
	Boxed       string // Boxed version if this is a primitive
}

// BuiltinScalars returns the default mappings for GraphQL built-in scalars.
func BuiltinScalars() map[string]ScalarInfo {
	return map[string]ScalarInfo{
		"String": {
			JavaType: "String",
			Imports:  nil,
		},
		"Int": {
			JavaType:      "Integer",
			PrimitiveType: "int",
			Imports:       nil,
		},
		"Float": {
			JavaType:      "Double",
			PrimitiveType: "double",
			Imports:       nil,
		},
		"Boolean": {
			JavaType:      "Boolean",
			PrimitiveType: "boolean",
			Imports:       nil,
		},
		"ID": {
			JavaType: "String",
			Imports:  nil,
		},
	}
}

// CommonScalars returns mappings for commonly used custom scalars.
func CommonScalars() map[string]ScalarInfo {
	return map[string]ScalarInfo{
		"DateTime": {
			JavaType: "LocalDateTime",
			Imports:  []string{"java.time.LocalDateTime"},
		},
		"Date": {
			JavaType: "LocalDate",
			Imports:  []string{"java.time.LocalDate"},
		},
		"Time": {
			JavaType: "LocalTime",
			Imports:  []string{"java.time.LocalTime"},
		},
		"Instant": {
			JavaType: "Instant",
			Imports:  []string{"java.time.Instant"},
		},
		"UUID": {
			JavaType: "UUID",
			Imports:  []string{"java.util.UUID"},
		},
		"BigDecimal": {
			JavaType: "BigDecimal",
			Imports:  []string{"java.math.BigDecimal"},
		},
		"BigInteger": {
			JavaType: "BigInteger",
			Imports:  []string{"java.math.BigInteger"},
		},
		"Long": {
			JavaType:      "Long",
			PrimitiveType: "long",
			Imports:       nil,
		},
		"Short": {
			JavaType:      "Short",
			PrimitiveType: "short",
			Imports:       nil,
		},
		"Byte": {
			JavaType:      "Byte",
			PrimitiveType: "byte",
			Imports:       nil,
		},
		"Char": {
			JavaType:      "Character",
			PrimitiveType: "char",
			Imports:       nil,
		},
		"JSON": {
			JavaType: "JsonNode",
			Imports:  []string{"com.fasterxml.jackson.databind.JsonNode"},
		},
		"JSONObject": {
			JavaType: "ObjectNode",
			Imports:  []string{"com.fasterxml.jackson.databind.node.ObjectNode"},
		},
		"JSONArray": {
			JavaType: "ArrayNode",
			Imports:  []string{"com.fasterxml.jackson.databind.node.ArrayNode"},
		},
		"URL": {
			JavaType: "URL",
			Imports:  []string{"java.net.URL"},
		},
		"URI": {
			JavaType: "URI",
			Imports:  []string{"java.net.URI"},
		},
	}
}

// PrimitiveToBoxed maps primitive types to their boxed equivalents.
var PrimitiveToBoxed = map[string]string{
	"int":     "Integer",
	"long":    "Long",
	"short":   "Short",
	"byte":    "Byte",
	"float":   "Float",
	"double":  "Double",
	"boolean": "Boolean",
	"char":    "Character",
}

// BoxedToPrimitive maps boxed types to their primitive equivalents.
var BoxedToPrimitive = map[string]string{
	"Integer":   "int",
	"Long":      "long",
	"Short":     "short",
	"Byte":      "byte",
	"Float":     "float",
	"Double":    "double",
	"Boolean":   "boolean",
	"Character": "char",
}

// IsPrimitive returns true if the type is a Java primitive.
func IsPrimitive(javaType string) bool {
	_, ok := PrimitiveToBoxed[javaType]
	return ok
}

// IsBoxed returns true if the type is a boxed primitive.
func IsBoxed(javaType string) bool {
	_, ok := BoxedToPrimitive[javaType]
	return ok
}

// BoxType converts a primitive to its boxed equivalent.
func BoxType(javaType string) string {
	if boxed, ok := PrimitiveToBoxed[javaType]; ok {
		return boxed
	}
	return javaType
}

// UnboxType converts a boxed type to its primitive equivalent.
func UnboxType(javaType string) string {
	if primitive, ok := BoxedToPrimitive[javaType]; ok {
		return primitive
	}
	return javaType
}
