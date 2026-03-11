package typed

import (
	"context"
	"testing"

	"gorm.io/cli/gorm/examples/models"
	generated "gorm.io/cli/gorm/examples/typed/models"
	"gorm.io/cli/gorm/typed"
)

func TestAssociation_Create_SingleParent(t *testing.T) {
	db := setupTestDB(t)
	users := seedUsers(t, db)
	u := users[0]

	ctx := context.Background()

	// Create one pet for the single parent
	_, err := typed.G[models.User](db).
		Where(generated.User.ID.Eq(u.ID)).
		Set(
			generated.User.Pets.Create(
				generated.Pet.Name.Set("test-pet"),
			),
		).
		Update(ctx)
	if err != nil {
		t.Fatalf("assoc create single failed: %v", err)
	}

	// Verify pet created and associated
	pets, err := typed.G[models.Pet](db).
		Where(generated.Pet.Name.Eq("test-pet")).
		Find(ctx)
	if err != nil {
		t.Fatalf("query pets failed: %v", err)
	}
	if len(pets) != 1 {
		t.Fatalf("expected 1 pet, got %d", len(pets))
	}
	if pets[0].UserID == nil || *pets[0].UserID != u.ID {
		t.Fatalf("pet not associated to user %d: %#v", u.ID, pets[0])
	}
}

func TestAssociation_Create_MultipleParents(t *testing.T) {
	db := setupTestDB(t)
	users := seedUsers(t, db)
	u1, u2 := users[0], users[1]

	ctx := context.Background()

	// Create one pet for each matched parent (two users)
	_, err := typed.G[models.User](db).
		Where(generated.User.Name.In(u1.Name, u2.Name)).
		Set(
			generated.User.Pets.Create(
				generated.Pet.Name.Set("multi-pet"),
			),
		).
		Update(ctx)
	if err != nil {
		t.Fatalf("assoc create multi failed: %v", err)
	}

	// Verify two pets created with correct names
	pets, err := typed.G[models.Pet](db).
		Where(generated.Pet.Name.Eq("multi-pet")).
		Find(ctx)
	if err != nil {
		t.Fatalf("query multi-pets failed: %v", err)
	}
	if len(pets) != 2 {
		t.Fatalf("expected 2 pets, got %d", len(pets))
	}
}

func TestAssociation_Preload_ChainedRelations(t *testing.T) {
	db := setupTestDB(t)
	users := seedUsers(t, db)
	u := users[0]

	pet := models.Pet{Name: "buddy", UserID: &u.ID}
	if err := db.Create(&pet).Error; err != nil {
		t.Fatalf("seed pet failed: %v", err)
	}
	if err := db.Create(&models.Toy{Name: "rope", OwnerID: pet.ID, OwnerType: "pets"}).Error; err != nil {
		t.Fatalf("seed toy failed: %v", err)
	}

	got, err := typed.G[models.User](db).
		Where(generated.User.ID.Eq(u.ID)).
		Preload(generated.UserRelations.Pets.Toy, func(db typed.PreloadBuilder) error { return nil }).
		First(context.Background())
	if err != nil {
		t.Fatalf("preload chained relation failed: %v", err)
	}
	if len(got.Pets) != 1 {
		t.Fatalf("expected 1 preloaded pet, got %d", len(got.Pets))
	}
	if got.Pets[0].Toy.ID == 0 || got.Pets[0].Toy.Name != "rope" {
		t.Fatalf("expected pet toy to be preloaded, got %#v", got.Pets[0].Toy)
	}
}

func TestAssociation_Update_WithConditions(t *testing.T) {
	db := setupTestDB(t)
	users := seedUsers(t, db)
	u := users[0]

	// Seed one pet for the user
	if err := db.Create(&models.Pet{Name: "old", UserID: &u.ID}).Error; err != nil {
		t.Fatalf("seed pet failed: %v", err)
	}

	ctx := context.Background()

	// Update the associated pet name where name='old'
	_, err := typed.G[models.User](db).
		Where(generated.User.ID.Eq(u.ID)).
		Set(
			generated.User.Pets.Where(generated.Pet.Name.Eq("old")).Update(
				generated.Pet.Name.Set("new"),
			),
		).
		Update(ctx)
	if err != nil {
		t.Fatalf("assoc update failed: %v", err)
	}

	// Verify pet name was updated
	pets, err := typed.G[models.Pet](db).
		Where(generated.Pet.Name.Eq("new")).
		Find(ctx)
	if err != nil {
		t.Fatalf("query updated pet failed: %v", err)
	}
	if len(pets) != 1 {
		t.Fatalf("expected 1 updated pet, got %d", len(pets))
	}
}
