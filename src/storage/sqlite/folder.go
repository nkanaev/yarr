package sqlite

import (
	"database/sql"
	"log"

	"github.com/nkanaev/yarr/src/storage/model"
)

func (s *SQLiteStorage) CreateFolder(title string) *model.Folder {
	expanded := true
	row := s.db.QueryRow(`
		insert into folders (title, is_expanded) values (:title, :is_expanded)
		on conflict (title) do update set title = :title
        returning id`,
		sql.Named("title", title),
		sql.Named("is_expanded", expanded),
	)
	var id int64
	err := row.Scan(&id)

	if err != nil {
		log.Print(err)
		return nil
	}
	return &model.Folder{Id: id, Title: title, IsExpanded: expanded}
}

func (s *SQLiteStorage) DeleteFolder(folderId int64) bool {
	_, err := s.db.Exec(`delete from folders where id = :id`, sql.Named("id", folderId))
	if err != nil {
		log.Print(err)
	}
	return err == nil
}

func (s *SQLiteStorage) UpdateFolder(folderId int64, params model.UpdateFolderParams) (bool, error) {
	_, err := s.db.Exec(`
		update folders set
			title       = coalesce(:title, title),
			is_expanded = coalesce(:is_expanded, is_expanded)
		where id = :id
	`,
		sql.Named("id", folderId),
		sql.Named("title", params.Title),
		sql.Named("is_expanded", params.IsExpanded),
	)
	if err != nil {
		log.Print(err)
		return false, err
	}
	return true, nil
}

func (s *SQLiteStorage) ListFolders() []model.Folder {
	result := make([]model.Folder, 0)
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
		var f model.Folder
		err = rows.Scan(&f.Id, &f.Title, &f.IsExpanded)
		if err != nil {
			log.Print(err)
			return result
		}
		result = append(result, f)
	}
	return result
}
