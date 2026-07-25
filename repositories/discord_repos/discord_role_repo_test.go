package discordrepos

import (
	"doko/gvn-ultimate-bot/models"
	"testing"
)

// Regression test: EditRole used to load the existing row and save it back
// unchanged, so edits silently never took effect.
func TestDiscordRoleRepo_EditRoleActuallyUpdatesFields(t *testing.T) {
	repo := NewDiscordRoleRepo(newTestDB(t))

	created, err := repo.CreateRole(&models.DiscordRole{
		NativeID:    "1",
		Name:        "Old Name",
		Mentionable: 0,
		Hoist:       0,
		Color:       111,
	})
	if err != nil {
		t.Fatalf("CreateRole returned error: %v", err)
	}
	if created.Name != "Old Name" {
		t.Fatalf("expected role to be created with 'Old Name', got %+v", created)
	}

	edited, err := repo.EditRole(&models.DiscordRole{
		NativeID:    "1",
		Name:        "New Name",
		Mentionable: 1,
		Hoist:       1,
		Color:       222,
	})
	if err != nil {
		t.Fatalf("EditRole returned error: %v", err)
	}
	if edited.Name != "New Name" || edited.Color != 222 {
		t.Fatalf("expected EditRole to persist new values, got %+v", edited)
	}

	found, err := repo.GetByNativeID("1")
	if err != nil {
		t.Fatalf("GetByNativeID returned error: %v", err)
	}
	if found.Name != "New Name" {
		t.Fatalf("expected persisted role name to be 'New Name', got %q", found.Name)
	}
}

func TestDiscordRoleRepo_UpsertCreatesNewRole(t *testing.T) {
	repo := NewDiscordRoleRepo(newTestDB(t))

	created, err := repo.Upsert(&models.DiscordRole{
		NativeID: "1",
		Name:     "Admin",
		Color:    123,
	})
	if err != nil {
		t.Fatalf("Upsert returned error: %v", err)
	}
	if created.Name != "Admin" {
		t.Fatalf("expected created role name 'Admin', got %+v", created)
	}
}

// Regression test: Upsert must refresh Discord-sourced fields (Name, Color,
// Mentionable, Hoist) but must NOT clobber locally managed fields
// (ImplicitType) that Discord has no concept of.
func TestDiscordRoleRepo_UpsertUpdatesDiscordFieldsButPreservesImplicitType(t *testing.T) {
	repo := NewDiscordRoleRepo(newTestDB(t))

	if _, err := repo.Upsert(&models.DiscordRole{
		NativeID:     "1",
		Name:         "Admin",
		Color:        111,
		ImplicitType: 7,
	}); err != nil {
		t.Fatalf("initial Upsert returned error: %v", err)
	}

	updated, err := repo.Upsert(&models.DiscordRole{
		NativeID: "1",
		Name:     "Administrator",
		Color:    222,
		// ImplicitType intentionally left at zero value, mirroring what a
		// Discord-sourced sync payload looks like (it has no such field).
	})
	if err != nil {
		t.Fatalf("second Upsert returned error: %v", err)
	}
	if updated.Name != "Administrator" || updated.Color != 222 {
		t.Fatalf("expected Discord-sourced fields to be refreshed, got %+v", updated)
	}
	if updated.ImplicitType != 7 {
		t.Fatalf("expected ImplicitType to be preserved as 7, got %d", updated.ImplicitType)
	}
}

func TestDiscordRoleRepo_DeleteNotIn_RemovesStaleRoles(t *testing.T) {
	repo := NewDiscordRoleRepo(newTestDB(t))

	for _, id := range []string{"1", "2", "3"} {
		if _, err := repo.Upsert(&models.DiscordRole{NativeID: id, Name: "role-" + id}); err != nil {
			t.Fatalf("Upsert(%s) returned error: %v", id, err)
		}
	}

	removed, err := repo.DeleteNotIn([]string{"1", "3"})
	if err != nil {
		t.Fatalf("DeleteNotIn returned error: %v", err)
	}
	if removed != 1 {
		t.Fatalf("expected 1 row removed, got %d", removed)
	}

	all, err := repo.ListRoles()
	if err != nil {
		t.Fatalf("ListRoles returned error: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 remaining roles, got %d", len(all))
	}
}

func TestDiscordRoleRepo_DeleteNotIn_EmptyListRemovesAll(t *testing.T) {
	repo := NewDiscordRoleRepo(newTestDB(t))

	for _, id := range []string{"1", "2"} {
		if _, err := repo.Upsert(&models.DiscordRole{NativeID: id}); err != nil {
			t.Fatalf("Upsert(%s) returned error: %v", id, err)
		}
	}

	removed, err := repo.DeleteNotIn(nil)
	if err != nil {
		t.Fatalf("DeleteNotIn returned error: %v", err)
	}
	if removed != 2 {
		t.Fatalf("expected all 2 rows removed, got %d", removed)
	}
}
