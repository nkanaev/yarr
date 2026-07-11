package sqlite

import (
	"database/sql"
	"log"

	"github.com/nkanaev/yarr/src/storage/model"
)

func (s *SQLiteStorage) CreateFeed(params model.CreateFeedParams) *model.Feed {
	title := params.Title
	if title == "" {
		title = params.FeedLink
	}
	row := s.db.QueryRow(`
		insert into feeds (title, description, link, feed_link, folder_id)
		values (:title, :description, :link, :feed_link, :folder_id)
		on conflict (feed_link) do update set folder_id = :folder_id
        returning id`,
		sql.Named("title", title),
		sql.Named("description", params.Description),
		sql.Named("link", params.Link),
		sql.Named("feed_link", params.FeedLink),
		sql.Named("folder_id", params.FolderID),
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

func (s *SQLiteStorage) DeleteFeed(feedId int64) bool {
	result, err := s.db.Exec(`delete from feeds where id = :id`, sql.Named("id", feedId))
	if err != nil {
		log.Print(err)
		return false
	}
	nrows, err := result.RowsAffected()
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return false
	}
	return nrows == 1
}

func (s *SQLiteStorage) UpdateFeed(feedId int64, params model.UpdateFeedParams) (bool, error) {
	_, err := s.db.Exec(`
		update feeds set
			title     = coalesce(:title, title),
			feed_link = coalesce(:feed_link, feed_link),
			folder_id = case when :update_folder_id then :folder_id else folder_id end,
			icon      = case when :update_icon then :icon else icon end
		where id = :id
	`,
		sql.Named("id", feedId),
		sql.Named("title", params.Title),
		sql.Named("feed_link", params.FeedLink),
		sql.Named("update_folder_id", params.FolderID.Set),
		sql.Named("folder_id", params.FolderID.Value),
		sql.Named("update_icon", params.Icon.Set),
		sql.Named("icon", params.Icon.Value),
	)
	if err != nil {
		log.Print(err)
		return false, err
	}
	return true, nil
}

func (s *SQLiteStorage) ListFeeds() []model.Feed {
	result := make([]model.Feed, 0)
	rows, err := s.db.Query(`
		select id, folder_id, title, description, link, feed_link, icon
		from feeds
		order by title collate nocase
	`)
	if err != nil {
		log.Print(err)
		return result
	}
	for rows.Next() {
		var f model.Feed
		err = rows.Scan(
			&f.Id,
			&f.FolderId,
			&f.Title,
			&f.Description,
			&f.Link,
			&f.FeedLink,
			&f.Icon,
		)
		if err != nil {
			log.Print(err)
			return result
		}
		result = append(result, f)
	}
	return result
}

func (s *SQLiteStorage) GetFeed(id int64) *model.Feed {
	var f model.Feed
	err := s.db.QueryRow(`
		select
			id, folder_id, title, link, feed_link,
			icon
		from feeds where id = :id
	`, sql.Named("id", id)).Scan(
		&f.Id, &f.FolderId, &f.Title, &f.Link, &f.FeedLink,
		&f.Icon,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return nil
	}
	return &f
}
