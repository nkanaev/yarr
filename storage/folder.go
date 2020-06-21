package storage

import "fmt"

type Folder struct {
	Id int64
	Title string
	IsExpanded bool
}

func (s *Storage) CreateFolder(title string) *Folder {
	expanded := true
	result, err := s.db.Exec(`
		insert into folders (title, is_expanded) values (?, ?)`,
		title, expanded,
	)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	id, idErr := result.LastInsertId()
	if idErr != nil {
		return nil
	}
	return &Folder{Id: id, Title: title, IsExpanded: expanded}
}

func (s *Storage) DeleteFolder(folderId int64) bool {
	_, err1 := s.db.Exec(`update feeds set folder_id = null where folder_id = ?`, folderId)
	_, err2 := s.db.Exec(`delete from folders where id = ?`, folderId)
	return err1 == nil && err2 == nil
}

func (s *Storage) RenameFolder(folderId int64, newTitle string) bool {
	_, err := s.db.Exec(`update folders set title = ? where id = ?`, newTitle, folderId)
	return err == nil
}

func (s *Storage) ToggleFolderExpanded(folderId int64, isExpanded bool) bool {
	_, err := s.db.Exec(`update folders set is_expanded = ? where id = ?`, isExpanded, folderId)
	return err == nil
}
