package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	// Parse command-line arguments
	schemaPath := flag.String("schema", "", "Path to GraphQL schema file")
	outputDir := flag.String("output", "output", "Directory for generated Java classes")
	packageName := flag.String("package", "com.example.model", "Java package name for generated classes")
	flag.Parse()

	if *schemaPath == "" {
		log.Fatal("Schema file path is required")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Read and parse schema
	schema, err := loadSchema(*schemaPath)
	if err != nil {
		log.Fatalf("Failed to load schema: %v", err)
	}

	// Generate Java classes
	generator := NewJavaGenerator(*packageName)
	for _, typeDef := range schema.Types {
		// Skip built-in GraphQL types
		if isBuiltinType(typeDef.Name) {
			continue
		}

		// Only process object types, interfaces, input objects, and enums
		switch typeDef.Kind {
		case ast.Object, ast.Interface, ast.InputObject, ast.Enum:
			javaClass, err := generator.GenerateClass(typeDef)
			if err != nil {
				log.Printf("Failed to generate Java class for %s: %v", typeDef.Name, err)
				continue
			}

			// Write the Java class to a file
			outputPath := filepath.Join(*outputDir, typeDef.Name+".java")
			if err := os.WriteFile(outputPath, []byte(javaClass), 0644); err != nil {
				log.Printf("Failed to write Java class file for %s: %v", typeDef.Name, err)
			} else {
				fmt.Printf("Generated %s\n", outputPath)
			}
		}
	}

	fmt.Println("Java classes generation completed")
}

func loadSchema(schemaPath string) (*ast.Schema, error) {
	// Read schema file
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Parse schema
	sources := []*ast.Source{
		{
			Name:  schemaPath,
			Input: string(schemaBytes),
		},
	}

	schema, err := gqlparser.LoadSchema(sources...)
	if err != nil {
		// Try with schema definition
		source := fmt.Sprintf("schema { query: Query }\n%s", string(schemaBytes))
		schema, err = gqlparser.LoadSchema(&ast.Source{
			Name:  schemaPath,
			Input: source,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to parse schema: %w", err)
		}
	}

	return schema, nil
}

func isBuiltinType(typeName string) bool {
	builtinTypes := map[string]bool{
		"String":              true,
		"Int":                 true,
		"Float":               true,
		"Boolean":             true,
		"ID":                  true,
		"__Schema":            true,
		"__Type":              true,
		"__Field":             true,
		"__Directive":         true,
		"__EnumValue":         true,
		"__InputValue":        true,
		"__TypeKind":          true,
		"__DirectiveLocation": true,
	}
	return builtinTypes[typeName] || strings.HasPrefix(typeName, "__")
}

type JavaGenerator struct {
	packageName string
}

func NewJavaGenerator(packageName string) *JavaGenerator {
	return &JavaGenerator{
		packageName: packageName,
	}
}

func (g *JavaGenerator) GenerateClass(typeDef *ast.Definition) (string, error) {
	var builder strings.Builder

	// Write package declaration
	builder.WriteString(fmt.Sprintf("package %s;\n\n", g.packageName))

	// Add imports
	g.writeImports(&builder, typeDef)

	// Add class javadoc
	if typeDef.Description != "" {
		builder.WriteString("/**\n")
		builder.WriteString(fmt.Sprintf(" * %s\n", typeDef.Description))
		builder.WriteString(" */\n")
	}

	// Generate class/enum/interface declaration
	switch typeDef.Kind {
	case ast.Object, ast.InputObject:
		builder.WriteString(fmt.Sprintf("public class %s {\n", typeDef.Name))
		g.generateFields(&builder, typeDef)
		g.generateGettersAndSetters(&builder, typeDef)
		builder.WriteString("}\n")
	case ast.Interface:
		builder.WriteString(fmt.Sprintf("public interface %s {\n", typeDef.Name))
		g.generateInterfaceMethods(&builder, typeDef)
		builder.WriteString("}\n")
	case ast.Enum:
		builder.WriteString(fmt.Sprintf("public enum %s {\n", typeDef.Name))
		g.generateEnumValues(&builder, typeDef)
		builder.WriteString("}\n")
	default:
		return "", fmt.Errorf("unsupported type kind: %s", typeDef.Kind)
	}

	return builder.String(), nil
}

func (g *JavaGenerator) writeImports(builder *strings.Builder, typeDef *ast.Definition) {
	// Add needed imports for collections, etc.
	importList := map[string]bool{}

	for _, field := range typeDef.Fields {
		// Check if field type requires imports
		if isListType(field.Type) {
			importList["import java.util.List;"] = true
			importList["import java.util.ArrayList;"] = true
		}
	}

	if len(importList) > 0 {
		for imp := range importList {
			builder.WriteString(imp + "\n")
		}
		builder.WriteString("\n")
	}
}

func (g *JavaGenerator) generateFields(builder *strings.Builder, typeDef *ast.Definition) {
	for _, field := range typeDef.Fields {
		// Add field javadoc if description exists
		if field.Description != "" {
			builder.WriteString(fmt.Sprintf("    /**\n     * %s\n     */\n", field.Description))
		}

		// Add field declaration
		javaType := mapGraphQLTypeToJava(field.Type)
		builder.WriteString(fmt.Sprintf("    private %s %s;\n\n", javaType, field.Name))
	}
}

func (g *JavaGenerator) generateGettersAndSetters(builder *strings.Builder, typeDef *ast.Definition) {
	for _, field := range typeDef.Fields {
		javaType := mapGraphQLTypeToJava(field.Type)
		fieldName := field.Name
		upperFieldName := strings.ToUpper(fieldName[:1]) + fieldName[1:]

		// Getter
		builder.WriteString(fmt.Sprintf("    public %s get%s() {\n", javaType, upperFieldName))
		builder.WriteString(fmt.Sprintf("        return this.%s;\n", fieldName))
		builder.WriteString("    }\n\n")

		// Setter
		builder.WriteString(fmt.Sprintf("    public void set%s(%s %s) {\n", upperFieldName, javaType, fieldName))
		builder.WriteString(fmt.Sprintf("        this.%s = %s;\n", fieldName, fieldName))
		builder.WriteString("    }\n\n")
	}
}

func (g *JavaGenerator) generateInterfaceMethods(builder *strings.Builder, typeDef *ast.Definition) {
	for _, field := range typeDef.Fields {
		// Add method javadoc if description exists
		if field.Description != "" {
			builder.WriteString(fmt.Sprintf("    /**\n     * %s\n     */\n", field.Description))
		}

		javaType := mapGraphQLTypeToJava(field.Type)
		fieldName := field.Name
		upperFieldName := strings.ToUpper(fieldName[:1]) + fieldName[1:]

		// Getter
		builder.WriteString(fmt.Sprintf("    %s get%s();\n\n", javaType, upperFieldName))
	}
}

func (g *JavaGenerator) generateEnumValues(builder *strings.Builder, typeDef *ast.Definition) {
	for i, enumValue := range typeDef.EnumValues {
		if enumValue.Description != "" {
			builder.WriteString(fmt.Sprintf("    /**\n     * %s\n     */\n", enumValue.Description))
		}

		builder.WriteString(fmt.Sprintf("    %s", enumValue.Name))
		if i < len(typeDef.EnumValues)-1 {
			builder.WriteString(",\n")
		} else {
			builder.WriteString(";\n")
		}
	}
}

func mapGraphQLTypeToJava(fieldType *ast.Type) string {
	if fieldType == nil {
		return "Object"
	}

	// Handle lists
	if isListType(fieldType) {
		elementType := mapGraphQLTypeToJava(fieldType.Elem)
		return fmt.Sprintf("List<%s>", elementType)
	}

	// Map scalar types
	typeName := fieldType.NamedType
	switch typeName {
	case "String":
		return "String"
	case "Int":
		if fieldType.NonNull {
			return "int"
		}
		return "Integer"
	case "Float":
		if fieldType.NonNull {
			return "double"
		}
		return "Double"
	case "Boolean":
		if fieldType.NonNull {
			return "boolean"
		}
		return "Boolean"
	case "ID":
		return "String"
	default:
		// Custom type
		return typeName
	}
}

func isListType(fieldType *ast.Type) bool {
	return fieldType.Elem != nil && fieldType.NamedType == ""
}
