package generator

import (
	"sort"
	"strings"
)

// ImportManager manages Java imports.
type ImportManager struct {
	imports   map[string]bool
	package_  string
}

// NewImportManager creates a new import manager.
func NewImportManager(pkg string) *ImportManager {
	return &ImportManager{
		imports:  make(map[string]bool),
		package_: pkg,
	}
}

// Add adds an import to the manager.
func (m *ImportManager) Add(imp string) {
	if imp == "" {
		return
	}
	// Don't import java.lang types
	if strings.HasPrefix(imp, "java.lang.") && !strings.Contains(imp[10:], ".") {
		return
	}
	// Don't import types from the same package
	if m.isSamePackage(imp) {
		return
	}
	m.imports[imp] = true
}

// AddAll adds multiple imports.
func (m *ImportManager) AddAll(imports []string) {
	for _, imp := range imports {
		m.Add(imp)
	}
}

// Has checks if an import exists.
func (m *ImportManager) Has(imp string) bool {
	return m.imports[imp]
}

// GetSorted returns all imports sorted.
func (m *ImportManager) GetSorted() []string {
	result := make([]string, 0, len(m.imports))
	for imp := range m.imports {
		result = append(result, imp)
	}
	sort.Strings(result)
	return result
}

// GetGrouped returns imports grouped by package prefix.
func (m *ImportManager) GetGrouped() [][]string {
	sorted := m.GetSorted()
	if len(sorted) == 0 {
		return nil
	}

	var groups [][]string
	var currentGroup []string
	currentPrefix := ""

	for _, imp := range sorted {
		prefix := getImportPrefix(imp)
		if currentPrefix != "" && prefix != currentPrefix {
			if len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
			}
			currentGroup = nil
		}
		currentGroup = append(currentGroup, imp)
		currentPrefix = prefix
	}

	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

// GenerateImportBlock generates the import statements.
func (m *ImportManager) GenerateImportBlock() string {
	groups := m.GetGrouped()
	if len(groups) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, group := range groups {
		for _, imp := range group {
			sb.WriteString("import ")
			sb.WriteString(imp)
			sb.WriteString(";\n")
		}
		if i < len(groups)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m *ImportManager) isSamePackage(imp string) bool {
	if m.package_ == "" {
		return false
	}
	lastDot := strings.LastIndex(imp, ".")
	if lastDot < 0 {
		return false
	}
	return imp[:lastDot] == m.package_
}

func getImportPrefix(imp string) string {
	// Group by first two package segments
	parts := strings.SplitN(imp, ".", 3)
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	return parts[0]
}

// StandardImportOrder defines the order for import grouping.
type StandardImportOrder struct {
	groups []string
}

// NewStandardImportOrder creates a new standard import order.
func NewStandardImportOrder() *StandardImportOrder {
	return &StandardImportOrder{
		groups: []string{
			"java.",
			"javax.",
			"jakarta.",
			"org.",
			"com.",
		},
	}
}

// GetGroupIndex returns the group index for an import.
func (o *StandardImportOrder) GetGroupIndex(imp string) int {
	for i, prefix := range o.groups {
		if strings.HasPrefix(imp, prefix) {
			return i
		}
	}
	return len(o.groups) // Other imports go last
}
