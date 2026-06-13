package postgres

import (
	"log"

	"github.com/nkanaev/yarr/src/storage/model"
)

func (s *PostgresStorage) CreateFolder(title string) *model.Folder {
	expanded := true
	row := s.db.QueryRow(`
		insert into folders (title, is_expanded) values ($1, $2)
		on conflict (title) do update set title = $1
		returning id`,
		title,
		expanded,
	)
	var id int64
	err := row.Scan(&id)

	if err != nil {
		log.Print(err)
		return nil
	}
	return &model.Folder{Id: id, Title: title, IsExpanded: expanded}
}

func (s *PostgresStorage) DeleteFolder(folderId int64) bool {
	_, err := s.db.Exec(`delete from folders where id = $1`, folderId)
	if err != nil {
		log.Print(err)
	}
	return err == nil
}

func (s *PostgresStorage) UpdateFolder(folderId int64, params model.UpdateFolderParams) (bool, error) {
	_, err := s.db.Exec(`
		update folders set
			title       = coalesce($2, title),
			is_expanded = coalesce($3, is_expanded)
		where id = $1
	`,
		folderId,
		params.Title,
		params.IsExpanded,
	)
	if err != nil {
		log.Print(err)
		return false, err
	}
	return true, nil
}

func (s *PostgresStorage) ListFolders() []model.Folder {
	result := make([]model.Folder, 0)
	rows, err := s.db.Query(`
		select id, title, is_expanded
		from folders
		order by lower(title)
	`)
	if err != nil {
		log.Print(err)
		return result
	}
	defer rows.Close()

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
