package storage

import (
	"log"
)

type Folder struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	IsExpanded bool   `json:"is_expanded"`
}

func (s *Storage) CreateFolder(title string) *Folder {
	expanded := true
	row := s.db.QueryRow(`
		insert into folders (title, is_expanded) values (?, ?)
		on conflict (title) do update set title = ?
        returning id`,
		title, expanded,
		// provide title again so that we can extract row id
		title,
	)
	var id int64
	err := row.Scan(&id)

	if err != nil {
		log.Print(err)
		return nil
	}
	return &Folder{Id: id, Title: title, IsExpanded: expanded}
}

func (s *Storage) DeleteFolder(folderId int64) bool {
	_, err := s.db.Exec(`delete from folders where id = ?`, folderId)
	if err != nil {
		log.Print(err)
	}
	return err == nil
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
		log.Print(err)
		return result
	}
	for rows.Next() {
		var f Folder
		err = rows.Scan(&f.Id, &f.Title, &f.IsExpanded)
		if err != nil {
			log.Print(err)
			return result
		}
		result = append(result, f)
	}
	return result
}
