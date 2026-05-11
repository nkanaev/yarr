package storage

import (
	"testing"
)

func TestUpdateFolder(t *testing.T) {
	db := testDB()
	folder := db.CreateFolder("old title")
	if folder.IsExpanded != true {
		t.Fatal("expected folder to be expanded by default")
	}

	t.Run("rename only", func(t *testing.T) {
		newTitle := "new title"
		ok, err := db.UpdateFolder(folder.Id, UpdateFolderParams{
			Title: &newTitle,
		})
		if !ok || err != nil {
			t.Fatalf("UpdateFolder failed: %v", err)
		}

		folders := db.ListFolders()
		if len(folders) != 1 || folders[0].Title != "new title" {
			t.Errorf("expected title to be updated, got %s", folders[0].Title)
		}
		if folders[0].IsExpanded != true {
			t.Error("expected expansion state to remain unchanged")
		}
	})

	t.Run("toggle expanded only", func(t *testing.T) {
		isExpanded := false
		ok, err := db.UpdateFolder(folder.Id, UpdateFolderParams{
			IsExpanded: &isExpanded,
		})
		if !ok || err != nil {
			t.Fatalf("UpdateFolder failed: %v", err)
		}

		folders := db.ListFolders()
		if len(folders) != 1 || folders[0].IsExpanded != false {
			t.Errorf("expected is_expanded to be false, got %v", folders[0].IsExpanded)
		}
		if folders[0].Title != "new title" {
			t.Error("expected title to remain unchanged")
		}
	})

	t.Run("update both", func(t *testing.T) {
		bothTitle := "both"
		isExpanded := true
		ok, err := db.UpdateFolder(folder.Id, UpdateFolderParams{
			Title:      &bothTitle,
			IsExpanded: &isExpanded,
		})
		if !ok || err != nil {
			t.Fatalf("UpdateFolder failed: %v", err)
		}

		folders := db.ListFolders()
		if len(folders) != 1 || folders[0].Title != "both" || folders[0].IsExpanded != true {
			t.Errorf("expected both to be updated, got title=%s expanded=%v", folders[0].Title, folders[0].IsExpanded)
		}
	})

	t.Run("update none", func(t *testing.T) {
		ok, err := db.UpdateFolder(folder.Id, UpdateFolderParams{})
		if !ok || err != nil {
			t.Fatalf("UpdateFolder failed: %v", err)
		}

		folders := db.ListFolders()
		if len(folders) != 1 || folders[0].Title != "both" || folders[0].IsExpanded != true {
			t.Errorf("expected no changes, got title=%s expanded=%v", folders[0].Title, folders[0].IsExpanded)
		}
	})
}
