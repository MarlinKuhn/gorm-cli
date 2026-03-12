package generator

import (
	"fmt"

	"gorm.io/cli/gorm/internal/gen"
)

func Generate(input, output string, typed bool) error {
	g := gen.Generator{
		Typed:   typed,
		Files:   map[string]*gen.File{},
		OutPath: output,
	}

	err := g.Process(input)
	if err != nil {
		return fmt.Errorf("error processing %s: %v", input, err)
	}

	err = g.Gen()
	if err != nil {
		return fmt.Errorf("error render template got error: %v", err)
	}

	return nil
}
