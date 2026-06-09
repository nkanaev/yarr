package sqlite

import (
	"cmp"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/nkanaev/yarr/src/storage/model"
)

// TODO: serialize/deserialize
func (m *MediaLinks) Scan(src any) error {
	switch data := src.(type) {
	case []byte:
		return json.Unmarshal(data, m)
	case string:
		return json.Unmarshal([]byte(data), m)
	default:
		return nil
	}
}

func (m MediaLinks) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (s *SQLiteStorage) CreateItems(items []Item) bool {
	tx, err := s.db.Begin()
	if err != nil {
		log.Print(err)
		return false
	}

	now := time.Now().UTC()

	slices.SortStableFunc(items, func(a, b model.Item) int {
		sa := a.Date.Format(time.RFC3339) + "::" + a.GUID
		sb := b.Date.Format(time.RFC3339) + "::" + b.GUID
		return cmp.Compare(sa, sb)
	})

	for _, item := range items {
		_, err = tx.Exec(`
			insert into items (
				guid, feed_id, title, link, date,
				content, media_links,
				date_arrived, last_arrived, status
			)
			values (
				:guid, :feed_id, :title, :link, strftime('%Y-%m-%d %H:%M:%f', :date),
				:content, :media_links,
				:date_arrived, :last_arrived, :status
			)
			on conflict (feed_id, guid) do update set
				last_arrived = :last_arrived`,
			sql.Named("guid", item.GUID),
			sql.Named("feed_id", item.FeedId),
			sql.Named("title", item.Title),
			sql.Named("link", item.Link),
			sql.Named("date", item.Date),
			sql.Named("content", item.Content),
			sql.Named("media_links", item.MediaLinks),
			sql.Named("date_arrived", now),
			sql.Named("last_arrived", now),
			sql.Named("status", model.UNREAD),
		)
		if err != nil {
			log.Print(err)
			if err = tx.Rollback(); err != nil {
				log.Print(err)
				return false
			}
			return false
		}
	}
	if err = tx.Commit(); err != nil {
		log.Print(err)
		return false
	}
	return true
}

func listQueryPredicate(filter ItemFilter, newestFirst bool) (string, []any) {
	cond := make([]string, 0)
	args := make([]any, 0)
	if filter.FolderID != nil {
		cond = append(cond, "i.feed_id in (select id from feeds where folder_id = :folder_id)")
		args = append(args, sql.Named("folder_id", *filter.FolderID))
	}
	if filter.FeedID != nil {
		cond = append(cond, "i.feed_id = :feed_id")
		args = append(args, sql.Named("feed_id", *filter.FeedID))
	}
	if filter.Status != nil {
		cond = append(cond, "i.status = :status")
		args = append(args, sql.Named("status", *filter.Status))
	}
	if filter.Search != nil {
		words := strings.Fields(*filter.Search)
		terms := make([]string, len(words))
		for idx, word := range words {
			terms[idx] = word + "*"
		}

		cond = append(
			cond,
			"i.id in (select rowid as id from search where search match :search)",
		)
		args = append(args, sql.Named("search", strings.Join(terms, " ")))
	}
	if filter.After != nil {
		compare := ">"
		if newestFirst {
			compare = "<"
		}
		cond = append(
			cond,
			fmt.Sprintf(
				"(i.date, i.id) %s (select date, id from items where id = :after_id)",
				compare,
			),
		)
		args = append(args, sql.Named("after_id", *filter.After))
	}
	if filter.IDs != nil && len(*filter.IDs) > 0 {
		qmarks := make([]string, len(*filter.IDs))
		for i, id := range *filter.IDs {
			name := fmt.Sprintf("id%d", i)
			qmarks[i] = ":" + name
			args = append(args, sql.Named(name, id))
		}
		cond = append(cond, "i.id in ("+strings.Join(qmarks, ",")+")")
	}
	if filter.SinceID != nil {
		cond = append(cond, "i.id > :since_id")
		args = append(args, sql.Named("since_id", filter.SinceID))
	}
	if filter.MaxID != nil {
		cond = append(cond, "i.id < :max_id")
		args = append(args, sql.Named("max_id", filter.MaxID))
	}
	if filter.Before != nil {
		cond = append(cond, "i.date < :before")
		args = append(args, sql.Named("before", filter.Before))
	}

	predicate := "1"
	if len(cond) > 0 {
		predicate = strings.Join(cond, " and ")
	}

	return predicate, args
}

func (s *SQLiteStorage) CountItems() int {
	var count int
	err := s.db.QueryRow(`select count(*) from items`).Scan(&count)
	if err != nil {
		log.Print(err)
		return 0
	}
	return count
}

func (s *SQLiteStorage) ListItems(
	filter ItemFilter,
	limit int,
	newestFirst bool,
	withContent bool,
) []Item {
	predicate, args := listQueryPredicate(filter, newestFirst)
	result := make([]Item, 0)

	order := "date desc, id desc"
	if !newestFirst {
		order = "date asc, id asc"
	}
	if filter.IDs != nil || filter.SinceID != nil {
		order = "i.id asc"
	}
	if filter.MaxID != nil {
		order = "i.id desc"
	}

	selectCols := "i.id, i.guid, i.feed_id, i.title, i.link, i.date, i.status, i.media_links"
	if withContent {
		selectCols += ", i.content"
	} else {
		selectCols += ", '' as content"
	}
	query := fmt.Sprintf(`
		select %s
		from items i
		where %s
		order by %s
		limit %d
		`, selectCols, predicate, order, limit)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Print(err)
		return result
	}
	for rows.Next() {
		var x Item
		err = rows.Scan(
			&x.Id, &x.GUID, &x.FeedId,
			&x.Title, &x.Link, &x.Date,
			&x.Status, &x.MediaLinks, &x.Content,
		)
		if err != nil {
			log.Print(err)
			return result
		}
		result = append(result, x)
	}
	return result
}

func (s *SQLiteStorage) GetItem(id int64) *Item {
	i := &Item{}
	err := s.db.QueryRow(`
		select
			i.id, i.guid, i.feed_id, i.title, i.link, i.content,
			i.date, i.status, i.media_links
		from items i
		where i.id = :id
	`, sql.Named("id", id)).Scan(
		&i.Id, &i.GUID, &i.FeedId, &i.Title, &i.Link, &i.Content,
		&i.Date, &i.Status, &i.MediaLinks,
	)
	if err != nil {
		log.Print(err)
		return nil
	}
	return i
}

func (s *SQLiteStorage) UpdateItemStatus(item_id int64, status ItemStatus) bool {
	_, err := s.db.Exec(`update items set status = :status where id = :id`,
		sql.Named("status", status),
		sql.Named("id", item_id),
	)
	return err == nil
}

func (s *SQLiteStorage) MarkItemsRead(filter MarkFilter) bool {
	predicate, args := listQueryPredicate(ItemFilter{
		FolderID: filter.FolderID,
		FeedID:   filter.FeedID,
		Before:   filter.Before,
	}, false)
	query := fmt.Sprintf(`
		update items as i set status = %d
		where %s and i.status != %d
		`, READ, predicate, STARRED)
	_, err := s.db.Exec(query, args...)
	if err != nil {
		log.Print(err)
	}
	return err == nil
}

type FeedStat struct {
	FeedId       int64 `json:"feed_id"`
	UnreadCount  int64 `json:"unread"`
	StarredCount int64 `json:"starred"`
}

func (s *SQLiteStorage) FeedStats() []FeedStat {
	result := make([]FeedStat, 0)
	rows, err := s.db.Query(fmt.Sprintf(`
		select
			feed_id,
			sum(case status when %d then 1 else 0 end),
			sum(case status when %d then 1 else 0 end)
		from items
		group by feed_id
	`, UNREAD, STARRED))
	if err != nil {
		log.Print(err)
		return result
	}
	for rows.Next() {
		stat := FeedStat{}
		rows.Scan(&stat.FeedId, &stat.UnreadCount, &stat.StarredCount)
		result = append(result, stat)
	}
	return result
}

var (
	itemsKeepSize = 50
	itemsKeepDays = 90
)

// Delete old articles from the database to cleanup space.
//
// The rules:
//   - Never delete starred entries.
//   - Keep at least 50 latest items for each feed.
//   - Delete entries older than 90 days relative to the latest arrived item in the same feed.
func (s *SQLiteStorage) DeleteOldItems() {
	result, err := s.db.Exec(`
		delete from items
		where id in (
			select id
			from (
				select
					id,
					row_number() over (partition by feed_id order by date desc) as rn,
					last_arrived,
					max(last_arrived) over (partition by feed_id) as max_la
				from items
				where status != :starred_status
			)
			where rn > :keep_size
			  and last_arrived < datetime(max_la, :keep_days_limit)
		)`,
		sql.Named("starred_status", STARRED),
		sql.Named("keep_size", itemsKeepSize),
		sql.Named("keep_days_limit", fmt.Sprintf("-%d days", itemsKeepDays)),
	)
	if err != nil {
		log.Print(err)
		return
	}
	numDeleted, err := result.RowsAffected()
	if err == nil && numDeleted > 0 {
		log.Printf("Deleted %d old items", numDeleted)

		if _, err := s.db.Exec("vacuum"); err != nil {
			log.Print(err)
		}
	}
}
