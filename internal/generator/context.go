package generator

import (
	"github.com/source-c/go-gql2j/internal/annotations"
	"github.com/source-c/go-gql2j/internal/config"
	"github.com/source-c/go-gql2j/internal/parser"
	"github.com/source-c/go-gql2j/internal/typemap"
)

// Context holds the generation context.
type Context struct {
	Config           *config.Config
	Schema           *parser.Schema
	TypeMapper       *typemap.TypeMapper
	NamingHelper     *NamingHelper
	LombokGen        *annotations.LombokGenerator
	ValidationGen    *annotations.ValidationGenerator
	CustomAnnotation *annotations.CustomAnnotationGenerator
}

// NewContext creates a new generation context.
func NewContext(cfg *config.Config, schema *parser.Schema) *Context {
	typeMapper := typemap.NewTypeMapper(cfg)
	typeMapper.SetSchemaTypes(schema.Types)

	return &Context{
		Config:           cfg,
		Schema:           schema,
		TypeMapper:       typeMapper,
		NamingHelper:     NewNamingHelper(&cfg.Java.Naming),
		LombokGen:        annotations.NewLombokGenerator(&cfg.Features.Lombok),
		ValidationGen:    annotations.NewValidationGenerator(&cfg.Features.Validation),
		CustomAnnotation: annotations.NewCustomAnnotationGenerator(),
	}
}

// TypeContext holds context for generating a specific type.
type TypeContext struct {
	*Context
	TypeDef  *parser.TypeDef
	TypeName string
	Imports  *ImportManager
}

// NewTypeContext creates a type-specific context.
func NewTypeContext(ctx *Context, typeDef *parser.TypeDef) *TypeContext {
	typeName := ctx.NamingHelper.GetTypeName(typeDef)
	return &TypeContext{
		Context:  ctx,
		TypeDef:  typeDef,
		TypeName: typeName,
		Imports:  NewImportManager(ctx.Config.Output.Package),
	}
}

// ShouldSkip returns true if the type should be skipped.
func (tc *TypeContext) ShouldSkip() bool {
	return parser.ExtractSkipDirective(tc.TypeDef.Directives) != nil
}

// GetVisibility returns the field visibility keyword.
func (tc *TypeContext) GetVisibility() string {
	switch tc.Config.Java.FieldVisibility {
	case config.VisibilityPrivate:
		return "private"
	case config.VisibilityProtected:
		return "protected"
	case config.VisibilityPackage:
		return "" // Package-private has no keyword
	case config.VisibilityPublic:
		return "public"
	default:
		return "private"
	}
}

// FieldContext holds context for generating a specific field.
type FieldContext struct {
	*TypeContext
	Field     *parser.FieldDef
	FieldName string
	JavaType  string
	Imports   []string
	IsNonNull bool
}

// NewFieldContext creates a field-specific context.
func NewFieldContext(tc *TypeContext, field *parser.FieldDef) (*FieldContext, error) {
	fieldName := tc.NamingHelper.GetFieldName(field)
	fieldName = EscapeJavaKeyword(fieldName)

	mapResult, err := tc.TypeMapper.MapFieldType(field)
	if err != nil {
		return nil, err
	}

	return &FieldContext{
		TypeContext: tc,
		Field:       field,
		FieldName:   fieldName,
		JavaType:    mapResult.JavaType,
		Imports:     mapResult.Imports,
		IsNonNull:   field.Type != nil && field.Type.NonNull,
	}, nil
}

// ShouldSkip returns true if the field should be skipped.
func (fc *FieldContext) ShouldSkip() bool {
	// Skip if @skip directive is present
	if parser.ExtractSkipDirective(fc.Field.Directives) != nil {
		return true
	}

	// Skip fields that reference introspection types (starting with __)
	if fc.Field.Type != nil {
		namedType := fc.Field.Type.NamedType()
		if len(namedType) >= 2 && namedType[:2] == "__" {
			return true
		}
	}

	return false
}

// IsBooleanType returns true if the field is a boolean type.
func (fc *FieldContext) IsBooleanType() bool {
	return fc.JavaType == "boolean" || fc.JavaType == "Boolean"
}
