package postgres

import (
	"database/sql"
	"log"

	"github.com/nkanaev/yarr/src/storage/model"
)

func (s *PostgresStorage) CreateFeed(params model.CreateFeedParams) *model.Feed {
	title := params.Title
	if title == "" {
		title = params.FeedLink
	}
	row := s.db.QueryRow(`
		insert into feeds (title, description, link, feed_link, folder_id)
		values ($1, $2, $3, $4, $5)
		on conflict (feed_link) do update set folder_id = $5
		returning id`,
		title,
		params.Description,
		params.Link,
		params.FeedLink,
		params.FolderID,
	)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &model.Feed{
		Id:          id,
		Title:       title,
		Description: params.Description,
		Link:        params.Link,
		FeedLink:    params.FeedLink,
		FolderId:    params.FolderID,
	}
}

func (s *PostgresStorage) DeleteFeed(feedId int64) bool {
	result, err := s.db.Exec(`delete from feeds where id = $1`, feedId)
	if err != nil {
		log.Print(err)
		return false
	}
	nrows, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		return false
	}
	return nrows == 1
}

func (s *PostgresStorage) UpdateFeed(feedId int64, params model.UpdateFeedParams) (bool, error) {
	_, err := s.db.Exec(`
		update feeds set
			title     = coalesce($2, title),
			feed_link = coalesce($3, feed_link),
			folder_id = case when $4 then $5 else folder_id end,
			icon      = case when $6 then $7 else icon end
		where id = $1
	`,
		feedId,
		params.Title,
		params.FeedLink,
		params.FolderID.Set,
		params.FolderID.Value,
		params.Icon.Set,
		params.Icon.Value,
	)
	if err != nil {
		log.Print(err)
		return false, err
	}
	return true, nil
}

func (s *PostgresStorage) ListFeeds() []model.Feed {
	result := make([]model.Feed, 0)
	rows, err := s.db.Query(`
		select id, folder_id, title, description, link, feed_link,
		       coalesce(length(icon), 0) > 0 as has_icon
		from feeds
		order by lower(title)
	`)
	if err != nil {
		log.Print(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var f model.Feed
		err = rows.Scan(
			&f.Id,
			&f.FolderId,
			&f.Title,
			&f.Description,
			&f.Link,
			&f.FeedLink,
			&f.HasIcon,
		)
		if err != nil {
			log.Print(err)
			return result
		}
		result = append(result, f)
	}
	return result
}

func (s *PostgresStorage) GetFeed(id int64) *model.Feed {
	var f model.Feed
	err := s.db.QueryRow(`
		select
			id, folder_id, title, link, feed_link,
			icon, coalesce(length(icon), 0) > 0 as has_icon
		from feeds where id = $1
	`, id).Scan(
		&f.Id, &f.FolderId, &f.Title, &f.Link, &f.FeedLink,
		&f.Icon, &f.HasIcon,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return nil
	}
	return &f
}
