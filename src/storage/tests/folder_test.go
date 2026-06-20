package tests

import (
	"testing"

	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

func TestCreateFolder(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		folder := db.CreateFolder("test-folder")
		if folder == nil || folder.Id == 0 {
			t.Fatal("expected folder with id")
		}
		if folder.Title != "test-folder" {
			t.Errorf("expected title 'test-folder', got %s", folder.Title)
		}
		if !folder.IsExpanded {
			t.Error("expected folder to be expanded by default")
		}

		// upsert: same title returns existing folder
		folder2 := db.CreateFolder("test-folder")
		if folder2 == nil || folder2.Id != folder.Id {
			t.Errorf("expected same folder id on upsert, got %d != %d", folder2.Id, folder.Id)
		}

		folders := db.ListFolders()
		if len(folders) != 1 || folders[0].Id != folder.Id {
			t.Errorf("expected folder in ListFolders")
		}
	})
}

func TestDeleteFolder(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		// delete non-existent returns true (err == nil)
		if !db.DeleteFolder(99999) {
			t.Error("expected true when deleting non-existent folder")
		}

		folder := db.CreateFolder("test")
		if !db.DeleteFolder(folder.Id) {
			t.Fatal("delete failed")
		}

		folders := db.ListFolders()
		if len(folders) != 0 {
			t.Errorf("expected 0 folders, got %d", len(folders))
		}

		// deleting again returns true
		if !db.DeleteFolder(folder.Id) {
			t.Error("expected true when deleting already-deleted folder")
		}
	})
}

func TestUpdateFolder(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		folder := db.CreateFolder("old title")
		if folder.IsExpanded != true {
			t.Fatal("expected folder to be expanded by default")
		}

		t.Run("rename only", func(t *testing.T) {
			newTitle := "new title"
			ok, err := db.UpdateFolder(folder.Id, model.UpdateFolderParams{
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
			ok, err := db.UpdateFolder(folder.Id, model.UpdateFolderParams{
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
			ok, err := db.UpdateFolder(folder.Id, model.UpdateFolderParams{
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
			ok, err := db.UpdateFolder(folder.Id, model.UpdateFolderParams{})
			if !ok || err != nil {
				t.Fatalf("UpdateFolder failed: %v", err)
			}

			folders := db.ListFolders()
			if len(folders) != 1 || folders[0].Title != "both" || folders[0].IsExpanded != true {
				t.Errorf("expected no changes, got title=%s expanded=%v", folders[0].Title, folders[0].IsExpanded)
			}
		})
	})
}
