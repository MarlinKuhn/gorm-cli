package generator

import (
	"bytes"
	_ "database/sql"
	_ "database/sql/driver"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
	"gorm.io/cli/gorm/genconfig"
	_ "gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var allowedInterfaces = []types.Type{
	loadNamedType("", "database/sql", "Scanner"),
	loadNamedType("", "database/sql/driver", "Valuer"),
	loadNamedType("", "gorm.io/gorm", "Valuer"),
	loadNamedType("", "gorm.io/gorm/schema", "SerializerInterface"),
}

type ExtractedSQL struct {
	Raw    string
	Where  string
	Select string
}

func extractSQL(comment string, methodName string) ExtractedSQL {
	comment = strings.TrimSpace(comment)

	if index := strings.Index(comment, "\n\n"); index != -1 {
		if strings.Contains(comment[index+2:], methodName) {
			comment = comment[:index]
		} else {
			comment = comment[index+2:]
		}
	}

	sql := strings.TrimPrefix(comment, methodName)
	if strings.HasPrefix(sql, "where(") && strings.HasSuffix(sql, ")") {
		content := strings.TrimSuffix(strings.TrimPrefix(sql, "where("), ")")
		content = strings.Trim(content, "\"")
		content = strings.TrimSpace(content)
		return ExtractedSQL{Where: content}
	} else if strings.HasPrefix(sql, "select(") && strings.HasSuffix(sql, ")") {
		content := strings.TrimSuffix(strings.TrimPrefix(sql, "select("), ")")
		content = strings.Trim(content, "\"")
		content = strings.TrimSpace(content)
		return ExtractedSQL{Select: content}
	}
	return ExtractedSQL{Raw: sql}
}

// ImplementsAllowedInterfaces reports whether typ or *typ implements any allowed interface.
func ImplementsAllowedInterfaces(typ types.Type) bool {
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}
	for _, t := range allowedInterfaces {
		iface, _ := t.Underlying().(*types.Interface)
		if types.Implements(typ, iface) || types.Implements(types.NewPointer(typ), iface) {
			return true
		}
	}
	return false
}

func IsUnderlyingComparable(typ types.Type) bool {
	underlying := typ.Underlying()
	if _, ok := underlying.(*types.Struct); ok {
		return false
	}
	return types.Comparable(underlying)
}

func findGoModDir(filename string) string {
	cmd := exec.Command("go", "env", "GOMOD")
	cmd.Dir = filepath.Dir(filename)
	out, _ := cmd.Output()
	return filepath.Dir(string(out))
}

// getCurrentPackagePath gets the full import path of the current file's package
func getCurrentPackagePath(filename string) string {
	cfg := &packages.Config{
		Mode: packages.NeedName,
		Dir:  findGoModDir(filename),
	}

	pkgs, err := packages.Load(cfg, filepath.Dir(filename))
	if err == nil && len(pkgs) > 0 && pkgs[0].PkgPath != "" {
		return pkgs[0].PkgPath
	}
	return ""
}

// loadNamedType returns a named type from a package with basic caching.
func loadNamedType(modRoot, pkgPath, name string) types.Type {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedName | packages.NeedDeps,
		Dir:  modRoot,
	}

	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil || len(pkgs) == 0 || pkgs[0].Types == nil {
		return nil
	}
	if obj := pkgs[0].Types.Scope().Lookup(name); obj != nil {
		return obj.Type()
	}
	return nil
}

// loadStructFromPackage loads a struct type definition from an external package by name
func loadNamedStructType(modRoot, pkgPath, name string) (*ast.StructType, error) {
	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedName,
		Dir:  modRoot,
	}

	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load package %q from %v: %w", pkgPath, modRoot, err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found for path %q from %v", pkgPath, modRoot)
	}

	for _, pkg := range pkgs {
		for _, syntax := range pkg.Syntax {
			for _, decl := range syntax.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				for _, spec := range gen.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if ok && ts.Name.Name == name {
						if st, ok := ts.Type.(*ast.StructType); ok {
							return st, nil
						}
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("struct %s not found in package %s", name, pkgPath)
}

// generateDBName generates database column name using GORM's NamingStrategy and COLUMN tag.
func generateDBName(fieldName, gormTag string) string {
	tagSettings := schema.ParseTagSetting(reflect.StructTag(gormTag).Get("gorm"), ";")
	if tagSettings["COLUMN"] != "" {
		return tagSettings["COLUMN"]
	}

	// Use GORM's NamingStrategy with IdentifierMaxLength: 64
	ns := schema.NamingStrategy{IdentifierMaxLength: 64}
	return ns.ColumnName("", fieldName)
}

// mergeImports appends imports from src into dst if not already present (by Path)
func mergeImports(dst *[]Import, src []Import) {
	existing := map[string]bool{}
	for _, i := range *dst {
		existing[i.Path] = true
	}
	for _, i := range src {
		if !existing[i.Path] {
			*dst = append(*dst, i)
			existing[i.Path] = true
		}
	}
}

// shouldSkipFile checks if a file contains the generated code header and should be skipped
func shouldSkipFile(filePath string) bool {
	if !strings.HasSuffix(filePath, ".go") {
		return true
	}

	content, err := os.ReadFile(filePath)
	return err == nil && bytes.Contains(content, []byte(codeGenHint))
}

func checksumMatchesGeneratedFile(filePath, checksum string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	found, ok := extractGeneratedChecksum(content)
	return ok && found == checksum
}

func extractGeneratedChecksum(content []byte) (string, bool) {
	const prefix = "// Generation checksum:"

	for _, line := range strings.Split(string(content), "\n") {
		if after, ok := strings.CutPrefix(line, prefix); ok {
			checksum := strings.TrimSpace(after)
			return checksum, checksum != ""
		}
	}

	return "", false
}

func generationChecksum(file *File, typed bool) string {
	payload := struct {
		Typed            bool             `json:"typed"`
		Source           string           `json:"source"`
		ApplicableConfig []configSnapshot `json:"applicable_config"`
	}{
		Typed:            typed,
		Source:           checksumBytes(mustReadFile(file.inputPath)),
		ApplicableConfig: snapshotConfigs(file.applicableConfigs),
	}

	b, _ := json.Marshal(payload)
	return checksumBytes(b)
}

type configSnapshot struct {
	OutPath           string            `json:"out_path"`
	FieldTypeMap      map[string]string `json:"field_type_map"`
	FieldNameMap      map[string]string `json:"field_name_map"`
	FileLevel         bool              `json:"file_level"`
	IncludeInterfaces []string          `json:"include_interfaces"`
	ExcludeInterfaces []string          `json:"exclude_interfaces"`
	IncludeStructs    []string          `json:"include_structs"`
	ExcludeStructs    []string          `json:"exclude_structs"`
}

func snapshotConfigs(configs []*genconfig.Config) []configSnapshot {
	snapshots := make([]configSnapshot, 0, len(configs))
	for _, cfg := range configs {
		snapshots = append(snapshots, configSnapshot{
			OutPath:           cfg.OutPath,
			FieldTypeMap:      stringifyMap(cfg.FieldTypeMap),
			FieldNameMap:      stringifyStringMap(cfg.FieldNameMap),
			FileLevel:         cfg.FileLevel,
			IncludeInterfaces: stringifySlice(cfg.IncludeInterfaces),
			ExcludeInterfaces: stringifySlice(cfg.ExcludeInterfaces),
			IncludeStructs:    stringifySlice(cfg.IncludeStructs),
			ExcludeStructs:    stringifySlice(cfg.ExcludeStructs),
		})
	}
	return snapshots
}

func stringifyMap(input map[any]any) map[string]string {
	out := make(map[string]string, len(input))
	for k, v := range input {
		out[fmt.Sprint(k)] = fmt.Sprint(v)
	}
	return out
}

func stringifyStringMap(input map[string]any) map[string]string {
	out := make(map[string]string, len(input))
	for k, v := range input {
		out[k] = fmt.Sprint(v)
	}
	return out
}

func stringifySlice(input []any) []string {
	out := make([]string, 0, len(input))
	for _, item := range input {
		out = append(out, fmt.Sprint(item))
	}
	return out
}

func mustReadFile(filePath string) []byte {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return content
}

// strLit returns the unquoted string if expr is a string literal; otherwise "".
func strLit(expr ast.Expr) string {
	if bl, ok := expr.(*ast.BasicLit); ok && bl.Kind == token.STRING {
		if s, err := strconv.Unquote(bl.Value); err == nil {
			return s
		}
	}

	return ""
}

func stripGeneric(s string) string {
	if i := strings.Index(s, "["); i >= 0 {
		return s[:i]
	}
	return s
}

// splitGenericArgs splits a generic type argument string into individual arguments.
func splitGenericArgs(s string) []string {
	var args []string
	depth := 0
	start := 0
	for i, char := range s {
		switch char {
		case '[':
			depth++
		case ']':
			depth--
		case ',':
			if depth == 0 {
				args = append(args, s[start:i])
				start = i + 1
			}
		}
	}
	if start < len(s) {
		args = append(args, s[start:])
	}
	return args
}
