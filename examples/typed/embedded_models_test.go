package typed

import (
	"testing"

	generated "gorm.io/cli/gorm/examples/typed/models"
)

func TestEmbeddedStructsExposeGeneratedFieldsAndRelations(t *testing.T) {
	if got := generated.EmbeddedUser.Code.Column().Name; got != "code" {
		t.Fatalf("expected embedded field column name %q, got %q", "code", got)
	}
	if got := generated.EmbeddedUser.Company.Name(); got != "Company" {
		t.Fatalf("expected embedded relation name %q, got %q", "Company", got)
	}
	if got := generated.EmbeddedUser.Pets.Name(); got != "Pets" {
		t.Fatalf("expected embedded relation name %q, got %q", "Pets", got)
	}
	if got := generated.EmbeddedUserRelations.Company.Name(); got != "Company" {
		t.Fatalf("expected embedded relation path %q, got %q", "Company", got)
	}
	if got := generated.EmbeddedUserRelations.Pets.Toy.Name(); got != "Pets.Toy" {
		t.Fatalf("expected embedded nested relation path %q, got %q", "Pets.Toy", got)
	}
}
