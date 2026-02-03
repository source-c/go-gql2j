package typemap

// CollectionInfo contains information about a Java collection type.
type CollectionInfo struct {
	Interface      string   // e.g., "List"
	Implementation string   // e.g., "ArrayList"
	Imports        []string // Required imports
}

// CollectionTypes returns information about supported collection types.
func CollectionTypes() map[string]CollectionInfo {
	return map[string]CollectionInfo{
		"List": {
			Interface:      "List",
			Implementation: "ArrayList",
			Imports:        []string{"java.util.List", "java.util.ArrayList"},
		},
		"Set": {
			Interface:      "Set",
			Implementation: "HashSet",
			Imports:        []string{"java.util.Set", "java.util.HashSet"},
		},
		"Collection": {
			Interface:      "Collection",
			Implementation: "ArrayList",
			Imports:        []string{"java.util.Collection", "java.util.ArrayList"},
		},
		"SortedSet": {
			Interface:      "SortedSet",
			Implementation: "TreeSet",
			Imports:        []string{"java.util.SortedSet", "java.util.TreeSet"},
		},
		"LinkedList": {
			Interface:      "List",
			Implementation: "LinkedList",
			Imports:        []string{"java.util.List", "java.util.LinkedList"},
		},
	}
}

// GetCollectionInfo returns info for a collection type, defaulting to List.
func GetCollectionInfo(collectionType string) CollectionInfo {
	collections := CollectionTypes()
	if info, ok := collections[collectionType]; ok {
		return info
	}
	return collections["List"]
}

// FormatCollectionType formats a collection type with element type.
func FormatCollectionType(collectionType, elementType string) string {
	info := GetCollectionInfo(collectionType)
	return info.Interface + "<" + elementType + ">"
}

// GetCollectionImports returns the imports needed for a collection type.
func GetCollectionImports(collectionType string) []string {
	info := GetCollectionInfo(collectionType)
	return info.Imports
}

// OptionalInfo contains information about Optional types.
type OptionalInfo struct {
	Type    string
	Imports []string
}

// GetOptionalInfo returns information for Optional types.
func GetOptionalInfo() OptionalInfo {
	return OptionalInfo{
		Type:    "Optional",
		Imports: []string{"java.util.Optional"},
	}
}

// FormatOptionalType formats an Optional type with element type.
func FormatOptionalType(elementType string) string {
	return "Optional<" + BoxType(elementType) + ">"
}
