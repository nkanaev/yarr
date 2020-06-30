package storage

type Feed struct {
	Id int64
	FolderId int64
	Title string
	Description string
	Link string
	FeedLink string
	Icon string
}

func (s *Storage) CreateFeed(title, description, link, feedLink, icon string, folderId int64) *Feed {
	result, err := s.db.Exec(`
		insert into feeds (title, description, link, feed_link, icon, folder_id) 
		values (?, ?, ?, ?, ?, ?)
		on conflict (feed_link) do update set folder_id=?`,
		title, description, link, feedLink, icon, intOrNil(folderId),
		intOrNil(folderId),
	)
	if err != nil {
		return nil
	}
	id, idErr := result.LastInsertId()
	if idErr != nil {
		return nil
	}
	return &Feed{
		Id: id,
		Title: title,
		Description: description,
		Link: link,
		FeedLink: feedLink,
		Icon: icon,
		FolderId: folderId,
	}
}

func (s *Storage) DeleteFeed(feedId int64) bool {
	_, err1 := s.db.Exec(`delete from items where feed_id = ?`, feedId)
	_, err2 := s.db.Exec(`delete from feeds where id = ?`, feedId)
	return err1 == nil && err2 == nil
}

func (s *Storage) RenameFeed(feedId int64, newTitle string) bool {
	_, err := s.db.Exec(`update feeds set title = ? where id = ?`, newTitle, feedId)
	return err == nil
}

func (s *Storage) UpdateFeedFolder(feedId int64, newFolderId int64) bool {
	_, err := s.db.Exec(`update feeds set folder_id = ? where id = ?`, intOrNil(newFolderId), feedId)
	return err == nil
}

func (s *Storage) ListFeeds() []Feed {
	result := make([]Feed, 0, 0)
	rows, err := s.db.Query(`
		select id, folder_id, title, description, link, feed_link, icon
		from feeds
	`)
	if err != nil {
		s.log.Print(err)
		return result
	}
	for rows.Next() {
		var f Feed
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
			s.log.Print(err)
			return result
		}
		result = append(result, f)
	}
	return result
}
