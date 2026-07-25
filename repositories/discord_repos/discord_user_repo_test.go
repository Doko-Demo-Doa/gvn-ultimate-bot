package discordrepos

import (
	"doko/gvn-ultimate-bot/models"
	"testing"
)

func TestDiscordUserRepo_UpsertCreatesNewUser(t *testing.T) {
	repo := NewDiscordUserRepo(newTestDB(t))

	created, err := repo.Upsert(&models.DiscordUser{
		NativeID: "1",
		Username: "alice",
		Nickname: "ali",
	})
	if err != nil {
		t.Fatalf("Upsert returned error: %v", err)
	}
	if created.Username != "alice" || created.Nickname != "ali" {
		t.Fatalf("created user has wrong fields: %+v", created)
	}

	found, err := repo.GetByNativeID("1")
	if err != nil {
		t.Fatalf("GetByNativeID returned error: %v", err)
	}
	if found == nil || found.Username != "alice" {
		t.Fatalf("expected to find created user, got %+v", found)
	}
}

// Regression test: Upsert must actually refresh Discord-sourced fields on an
// existing row, not silently ignore the incoming data.
func TestDiscordUserRepo_UpsertUpdatesExistingUser(t *testing.T) {
	repo := NewDiscordUserRepo(newTestDB(t))

	if _, err := repo.Upsert(&models.DiscordUser{
		NativeID: "1",
		Username: "alice",
		Nickname: "old-nick",
		Avatar:   "old-avatar",
	}); err != nil {
		t.Fatalf("initial Upsert returned error: %v", err)
	}

	updated, err := repo.Upsert(&models.DiscordUser{
		NativeID: "1",
		Username: "alice",
		Nickname: "new-nick",
		Avatar:   "new-avatar",
	})
	if err != nil {
		t.Fatalf("second Upsert returned error: %v", err)
	}
	if updated.Nickname != "new-nick" || updated.Avatar != "new-avatar" {
		t.Fatalf("expected Upsert to overwrite nickname/avatar, got %+v", updated)
	}

	found, err := repo.GetByNativeID("1")
	if err != nil {
		t.Fatalf("GetByNativeID returned error: %v", err)
	}
	if found.Nickname != "new-nick" {
		t.Fatalf("expected persisted nickname to be updated, got %q", found.Nickname)
	}

	all, err := repo.ListAll()
	if err != nil {
		t.Fatalf("ListAll returned error: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected exactly one row after upserting the same NativeID twice, got %d", len(all))
	}
}

func TestDiscordUserRepo_GetByNativeID_NotFoundReturnsNilNil(t *testing.T) {
	repo := NewDiscordUserRepo(newTestDB(t))

	found, err := repo.GetByNativeID("does-not-exist")
	if err != nil {
		t.Fatalf("expected no error for a missing user, got: %v", err)
	}
	if found != nil {
		t.Fatalf("expected nil for a missing user, got: %+v", found)
	}
}

func TestDiscordUserRepo_DeleteNotIn_RemovesStaleUsers(t *testing.T) {
	repo := NewDiscordUserRepo(newTestDB(t))

	for _, id := range []string{"1", "2", "3"} {
		if _, err := repo.Upsert(&models.DiscordUser{NativeID: id, Username: "u" + id}); err != nil {
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

	all, err := repo.ListAll()
	if err != nil {
		t.Fatalf("ListAll returned error: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 remaining users, got %d", len(all))
	}
	for _, u := range all {
		if u.NativeID == "2" {
			t.Fatalf("expected user 2 to have been deleted, but it's still present")
		}
	}
}

func TestDiscordUserRepo_DeleteNotIn_EmptyListRemovesAll(t *testing.T) {
	repo := NewDiscordUserRepo(newTestDB(t))

	for _, id := range []string{"1", "2"} {
		if _, err := repo.Upsert(&models.DiscordUser{NativeID: id}); err != nil {
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
