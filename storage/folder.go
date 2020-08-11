package storage

import (
	"fmt"
)

type Folder struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	IsExpanded bool   `json:"is_expanded"`
}

func (s *Storage) CreateFolder(title string) *Folder {
	expanded := true
	result, err := s.db.Exec(`
		insert into folders (title, is_expanded) values (?, ?)
		on conflict (title) do nothing`,
		title, expanded,
	)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var id int64
	numrows, err := result.RowsAffected()
	if err != nil {
		s.log.Print(err)
		return nil
	}
	if numrows == 1 {
		id, err = result.LastInsertId()
		if err != nil {
			s.log.Print(err)
			return nil
		}
	} else {
		err = s.db.QueryRow(`select id, is_expanded from folders where title=?`, title).Scan(&id, &expanded)
		if err != nil {
			s.log.Print(err)
			return nil
		}
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

func (s *Storage) ListFolders() []Folder {
	result := make([]Folder, 0, 0)
	rows, err := s.db.Query(`
		select id, title, is_expanded
		from folders
		order by title collate nocase
	`)
	if err != nil {
		s.log.Print(err)
		return result
	}
	for rows.Next() {
		var f Folder
		err = rows.Scan(&f.Id, &f.Title, &f.IsExpanded)
		if err != nil {
			s.log.Print(err)
			return result
		}
		result = append(result, f)
	}
	return result
}
