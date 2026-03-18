package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"text/template"
)

func TestParseTemplate(t *testing.T) {
	if _, err := template.New("").Parse(pkgTmpl); err != nil {
		t.Errorf("failed to parse template, got %v", err)
	}
}

func TestLoadNamedTypes(t *testing.T) {
	for _, i := range allowedInterfaces {
		if i == nil {
			t.Fatalf("failed to load named type, got nil")
		}
	}
}

func TestGeneratorWithQueryInterface(t *testing.T) {
	inputPath, err := filepath.Abs("../examples/query.go")
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	goldenPath, err := filepath.Abs("../examples/output/query.go")
	if err != nil {
		t.Fatalf("failed to get absolute output path: %v", err)
	}

	outputDir := filepath.Join(t.TempDir(), "output")

	g := &Generator{Files: map[string]*File{}, OutPath: outputDir}

	if err := g.Process(inputPath); err != nil {
		t.Fatalf("Process error: %v", err)
	}
	if err := g.Gen(); err != nil {
		t.Fatalf("Gen error: %v", err)
	}

	files, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("failed to read output dir: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("no files were generated in %s", outputDir)
	}

	goldenBytes, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file %s: %v", goldenPath, err)
	}
	goldenStr := string(goldenBytes)

	generatedFile := filepath.Join(outputDir, files[0].Name())
	genBytes, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("failed to read generated file %s: %v", generatedFile, err)
	}
	generatedStr := string(genBytes)

	if _, err := parser.ParseFile(token.NewFileSet(), generatedFile, genBytes, parser.AllErrors); err != nil {
		t.Errorf("generated code %s has invalid Go syntax: %v", generatedFile, err)
	}

	if goldenStr != generatedStr {
		t.Errorf("generated file differs from golden file\nGOLDEN: %s\nGENERATED: %s\n%s",
			goldenPath, generatedFile, generatedStr)
	}
}

func TestExcludeInterfacesSkipsInvalidInterfaces(t *testing.T) {
	writeSample := func(dir string, withExclude bool) string {
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module temp.test\n\ngo 1.21\n"), 0o644); err != nil {
			t.Fatalf("write go.mod: %v", err)
		}

		cfg := "ExcludeInterfaces: []any{Entity(nil)},"
		if !withExclude {
			cfg = ""
		}

		src := fmt.Sprintf(`package sample

import "gorm.io/cli/gorm/genconfig"

var _ = genconfig.Config{
	%s
}

type Entity interface {
	TableName() string
}
`, cfg)

		path := filepath.Join(dir, "sample.go")
		if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
			t.Fatalf("write sample.go: %v", err)
		}
		return path
	}

	runGen := func(file string) error {
		g := &Generator{Files: map[string]*File{}, OutPath: filepath.Join(filepath.Dir(file), "out")}
		if err := g.Process(file); err != nil {
			return err
		}
		return g.Gen()
	}

	withExcludeDir := t.TempDir()
	withExcludeFile := writeSample(withExcludeDir, true)
	if err := runGen(withExcludeFile); err != nil {
		t.Fatalf("generator should succeed when interface is excluded: %v", err)
	}

	withoutExcludeDir := t.TempDir()
	withoutExcludeFile := writeSample(withoutExcludeDir, false)
	if err := runGen(withoutExcludeFile); err == nil {
		t.Fatalf("expected generator failure when interface is not excluded")
	}
}

func TestProcessStructType(t *testing.T) {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "../examples/models/user.go", nil, parser.AllErrors)
	if err != nil {
		t.Fatalf("failed to parse file: %v", err)
	}

	var structType *ast.StructType

	ast.Inspect(file, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if ok && typeSpec.Name.Name == "User" {
			structType = typeSpec.Type.(*ast.StructType)
			return false
		}
		return true
	})

	if structType == nil {
		t.Fatalf("failed to find User struct")
	}

	expected := Struct{
		Name: "User",
		Fields: []Field{
			{Name: "ID", DBName: "id", GoType: "uint"},
			{Name: "CreatedAt", DBName: "created_at", GoType: "time.Time"},
			{Name: "UpdatedAt", DBName: "updated_at", GoType: "time.Time"},
			{Name: "DeletedAt", DBName: "deleted_at", GoType: "gorm.io/gorm.DeletedAt"},
			{Name: "Name", DBName: "name", GoType: "string"},
			{Name: "Age", DBName: "age", GoType: "int"},
			{Name: "Birthday", DBName: "birthday", GoType: "*time.Time"},
			{Name: "Score", DBName: "score", GoType: "sql.NullInt64"},
			{Name: "LastLogin", DBName: "last_login", GoType: "sql.NullTime"},
			{Name: "Account", DBName: "account", GoType: "Account"},
			{Name: "Pets", DBName: "pets", GoType: "[]*Pet"},
			{Name: "Toys", DBName: "toys", GoType: "[]Toy"},
			{Name: "CompanyID", DBName: "company_id", GoType: "*int"},
			{Name: "Company", DBName: "company", GoType: "Company"},
			{Name: "ManagerID", DBName: "manager_id", GoType: "*uint"},
			{Name: "Manager", DBName: "manager", GoType: "*User"},
			{Name: "Team", DBName: "team", GoType: "[]User"},
			{Name: "Languages", DBName: "languages", GoType: "[]Language"},
			{Name: "Friends", DBName: "friends", GoType: "[]*User"},
			{Name: "Role", DBName: "role", GoType: "string"},
			{Name: "IsAdult", DBName: "is_adult", GoType: "bool"},
			{Name: "Profile", DBName: "profile", GoType: "string", NamedGoType: "json"},
			{Name: "AwardTypes", DBName: "award_types", GoType: "datatypes.JSONSlice[int]"},
			{Name: "TagTypes", DBName: "tag_types", GoType: "datatypes.JSONSlice[UserTagType]"},
			{Name: "Tag", DBName: "tag", GoType: "UserTagType"},
			{Name: "Enum", DBName: "enum", GoType: "enum.Enum"}, // 添加
			{Name: "Enum2", DBName: "enum2", GoType: "enum2.Enum"},
		},
	}

	p := File{
		Imports: []Import{
			{Name: "gorm", Path: "gorm.io/gorm"},
		},
	}

	result := p.processStructType(&ast.TypeSpec{Name: &ast.Ident{Name: "User"}}, structType, "")
	// Only compare stable fields (Name, DBName, GoType); ignore tags/alias and internal pointers.
	trimmed := Struct{Name: result.Name}
	for _, f := range result.Fields {
		trimmed.Fields = append(trimmed.Fields, Field{Name: f.Name, DBName: f.DBName, GoType: f.GoType, NamedGoType: f.NamedGoType})
	}
	if !reflect.DeepEqual(trimmed, expected) {
		t.Errorf("Expected %+v, \n got %+v", expected, trimmed)
	}
}

func TestGenericRelationsIncludeTypeParametersInHelpers(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module temp.test\n\ngo 1.21\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	src := `package sample

type Node[T any] struct {
	Next     *Node[T]
	Children []Node[T]
	Value    T
}
`

	inputPath := filepath.Join(dir, "sample.go")
	if err := os.WriteFile(inputPath, []byte(src), 0o644); err != nil {
		t.Fatalf("write sample.go: %v", err)
	}

	outputDir := filepath.Join(dir, "out")
	g := &Generator{Files: map[string]*File{}, OutPath: outputDir}

	if err := g.Process(inputPath); err != nil {
		t.Fatalf("Process error: %v", err)
	}
	if err := g.Gen(); err != nil {
		t.Fatalf("Gen error: %v", err)
	}

	files, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("read output dir: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected one generated file, got %d", len(files))
	}

	generatedFile := filepath.Join(outputDir, files[0].Name())
	genBytes, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	generated := string(genBytes)

	for _, want := range []string{
		"type nodeRelationsFields[T any] struct {",
		"type nodeStructRelation[T any] struct {",
		"type nodeSliceRelation[T any] struct {",
		"func newNodeStructRelation[T any](prefix string, depth int) *nodeStructRelation[T] {",
		"func newNodeSliceRelation[T any](prefix string, depth int) *nodeSliceRelation[T] {",
		"field.Struct[sample.Node[T]]{}.WithName(prefix)",
		"field.Slice[sample.Node[T]]{}.WithName(prefix)",
		"Next:     newNodeStructRelation[T](strings.TrimPrefix(prefix+\".Next\", \".\"), depth-1),",
		"Children: newNodeSliceRelation[T](strings.TrimPrefix(prefix+\".Children\", \".\"), depth-1),",
	} {
		if !strings.Contains(generated, want) {
			t.Fatalf("generated output missing %q\n%s", want, generated)
		}
	}
}

func TestEmbeddedStructFieldsAndRelationsAcrossFiles(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module temp.test\n\ngo 1.21\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	baseSrc := `package sample

type Meta struct {
	Code string
	Pets []Pet
}

type Audit struct {
	Meta
	Company Company
}

type Company struct {
	ID uint
}

type Pet struct {
	ID uint
}
`
	if err := os.WriteFile(filepath.Join(dir, "base.go"), []byte(baseSrc), 0o644); err != nil {
		t.Fatalf("write base.go: %v", err)
	}

	userSrc := `package sample

type User struct {
	Audit
	Company Company
	Name string
}
`
	if err := os.WriteFile(filepath.Join(dir, "user.go"), []byte(userSrc), 0o644); err != nil {
		t.Fatalf("write user.go: %v", err)
	}

	g := &Generator{Files: map[string]*File{}, OutPath: filepath.Join(dir, "out")}
	if err := g.Process(dir); err != nil {
		t.Fatalf("Process error: %v", err)
	}

	var userFile *File
	var userStruct *Struct
	for _, file := range g.Files {
		for i := range file.Structs {
			if file.Structs[i].Name == "User" {
				userFile = file
				userStruct = &file.Structs[i]
				break
			}
		}
		if userStruct != nil {
			break
		}
	}

	if userStruct == nil || userFile == nil {
		t.Fatalf("failed to find parsed User struct")
	}

	var gotFields []Field
	for _, f := range userStruct.Fields {
		gotFields = append(gotFields, Field{Name: f.Name, DBName: f.DBName, GoType: f.GoType})
	}

	wantFields := []Field{
		{Name: "Code", DBName: "code", GoType: "string"},
		{Name: "Pets", DBName: "pets", GoType: "[]temp.test.Pet"},
		{Name: "Company", DBName: "company", GoType: "temp.test.Company"},
		{Name: "Name", DBName: "name", GoType: "string"},
	}
	if !reflect.DeepEqual(gotFields, wantFields) {
		t.Fatalf("expected embedded fields %+v, got %+v", wantFields, gotFields)
	}

	companyCount := 0
	for _, f := range userStruct.Fields {
		if f.Name == "Company" {
			companyCount++
		}
	}
	if companyCount != 1 {
		t.Fatalf("expected embedded shadowed field to appear once, got %d entries", companyCount)
	}

	userStruct.RelationFields = userFile.buildRelationFields(*userStruct)

	gotRelations := make([]string, 0, len(userStruct.RelationFields))
	for _, rel := range userStruct.RelationFields {
		gotRelations = append(gotRelations, rel.Name)
	}

	wantRelations := []string{"Pets", "Company"}
	if !reflect.DeepEqual(gotRelations, wantRelations) {
		t.Fatalf("expected embedded relations %v, got %v", wantRelations, gotRelations)
	}
}

func TestEmbeddedStructShadowedFieldsStayUnique(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module temp.test\n\ngo 1.21\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	src := `package sample

type companyDefaults struct {
	Company Company
}

type financeDefaults struct {
	companyDefaults
	RecipientAddress string ` + "`gorm:\"column:recipient_address\"`" + `
	HeadNote         string ` + "`gorm:\"column:head_note\"`" + `
}

type Estimate struct {
	financeDefaults
	RecipientAddress *string ` + "`gorm:\"column:recipient_address\"`" + `
	HeadNote         string ` + "`gorm:\"column:head_note\"`" + `
	Company          Company
}

type Company struct {
	ID uint
}
`
	inputPath := filepath.Join(dir, "sample.go")
	if err := os.WriteFile(inputPath, []byte(src), 0o644); err != nil {
		t.Fatalf("write sample.go: %v", err)
	}

	g := &Generator{Files: map[string]*File{}, OutPath: filepath.Join(dir, "out")}
	if err := g.Process(inputPath); err != nil {
		t.Fatalf("Process error: %v", err)
	}

	var estimateFile *File
	var estimateStruct *Struct
	for _, file := range g.Files {
		for i := range file.Structs {
			if file.Structs[i].Name == "Estimate" {
				estimateFile = file
				estimateStruct = &file.Structs[i]
				break
			}
		}
	}

	if estimateStruct == nil || estimateFile == nil {
		t.Fatalf("failed to find parsed Estimate struct")
	}

	nameCounts := map[string]int{}
	dbNameCounts := map[string]int{}
	for _, f := range estimateStruct.Fields {
		nameCounts[f.Name]++
		dbNameCounts[f.DBName]++
	}

	for key, count := range map[string]int{
		"RecipientAddress": nameCounts["RecipientAddress"],
		"HeadNote":         nameCounts["HeadNote"],
		"Company":          nameCounts["Company"],
	} {
		if count != 1 {
			t.Fatalf("expected field %s exactly once, got %d", key, count)
		}
	}

	for key, count := range map[string]int{
		"recipient_address": dbNameCounts["recipient_address"],
		"head_note":         dbNameCounts["head_note"],
		"company":           dbNameCounts["company"],
	} {
		if count != 1 {
			t.Fatalf("expected column %s exactly once, got %d", key, count)
		}
	}

	estimateStruct.RelationFields = estimateFile.buildRelationFields(*estimateStruct)
	companyRelations := 0
	for _, rel := range estimateStruct.RelationFields {
		if rel.Name == "Company" {
			companyRelations++
		}
	}
	if companyRelations != 1 {
		t.Fatalf("expected Company relation exactly once, got %d", companyRelations)
	}
}

func TestNamedSliceRelationUsesUnderlyingModel(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module temp.test\n\ngo 1.21\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	src := `package sample

type EmploymentContract struct {
	ID uint
}

type EmploymentContracts []*EmploymentContract

type User struct {
	Contracts EmploymentContracts
}
`

	inputPath := filepath.Join(dir, "sample.go")
	if err := os.WriteFile(inputPath, []byte(src), 0o644); err != nil {
		t.Fatalf("write sample.go: %v", err)
	}

	g := &Generator{Files: map[string]*File{}, OutPath: filepath.Join(dir, "out")}
	if err := g.Process(inputPath); err != nil {
		t.Fatalf("Process error: %v", err)
	}

	var userFile *File
	var userStruct *Struct
	for _, file := range g.Files {
		for i := range file.Structs {
			if file.Structs[i].Name == "User" {
				userFile = file
				userStruct = &file.Structs[i]
				break
			}
		}
		if userStruct != nil {
			break
		}
	}

	if userStruct == nil || userFile == nil {
		t.Fatalf("failed to find parsed User struct")
	}

	userStruct.RelationFields = userFile.buildRelationFields(*userStruct)
	if len(userStruct.RelationFields) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(userStruct.RelationFields))
	}

	rel := userStruct.RelationFields[0]
	if rel.Name != "Contracts" {
		t.Fatalf("expected relation name %q, got %q", "Contracts", rel.Name)
	}
	if !rel.IsSlice {
		t.Fatalf("expected named slice relation to be treated as slice")
	}
	if rel.Target == nil || rel.Target.Name != "EmploymentContract" {
		t.Fatalf("expected relation target %q, got %#v", "EmploymentContract", rel.Target)
	}
	if rel.BaseType != "field.Slice[sample.EmploymentContract]" {
		t.Fatalf("expected slice relation base type, got %q", rel.BaseType)
	}
}

func TestGenericCustomValueTypeUsesFieldWrapper(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module temp.test\n\ngo 1.21\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	src := `package sample

import (
	"database/sql/driver"
	"encoding/json"
)

type JSON[T any] []T

func (j *JSON[T]) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return nil
	}
}

func (j JSON[T]) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type User struct {
	Tags JSON[string]
}
`

	inputPath := filepath.Join(dir, "sample.go")
	if err := os.WriteFile(inputPath, []byte(src), 0o644); err != nil {
		t.Fatalf("write sample.go: %v", err)
	}

	outputDir := filepath.Join(dir, "out")
	g := &Generator{Files: map[string]*File{}, OutPath: outputDir}

	if err := g.Process(inputPath); err != nil {
		t.Fatalf("Process error: %v", err)
	}
	if err := g.Gen(); err != nil {
		t.Fatalf("Gen error: %v", err)
	}

	files, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("read output dir: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected one generated file, got %d", len(files))
	}

	generatedFile := filepath.Join(outputDir, files[0].Name())
	genBytes, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	generated := string(genBytes)

	if !strings.Contains(generated, "Tags field.Field[sample.JSON[string]]") {
		t.Fatalf("generated output should use field.Field for custom value type\n%s", generated)
	}

	if strings.Contains(generated, "Tags field.Struct[sample.JSON[string]]") {
		t.Fatalf("generated output should not use field.Struct for custom value type\n%s", generated)
	}
}
