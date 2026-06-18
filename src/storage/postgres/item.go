package postgres

import (
	"cmp"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/nkanaev/yarr/src/content/htmlutil"
	"github.com/nkanaev/yarr/src/storage/model"
)

type MediaLinks model.MediaLinks

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

func (s *PostgresStorage) CreateItems(items []model.Item) bool {
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
		searchText := item.Title + " " + htmlutil.ExtractText(item.Content)
		_, err = tx.Exec(`
			insert into items (
				guid, feed_id, title, link, date,
				content, media_links,
				date_arrived, last_arrived, status,
				search
			)
			values (
				$1, $2, $3, $4, $5,
				$6, $7,
				$8, $9, $10,
				to_tsvector('simple', $11)
			)
			on conflict (feed_id, guid) do update set
				last_arrived = excluded.last_arrived`,
			item.GUID,
			item.FeedId,
			item.Title,
			item.Link,
			item.Date,
			item.Content,
			MediaLinks(item.MediaLinks),
			now,
			now,
			model.UNREAD,
			searchText,
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

func listQueryPredicate(filter model.ItemFilter, newestFirst bool) (string, []any) {
	cond := make([]string, 0)
	args := make([]any, 0)
	n := 0

	next := func() int {
		n++
		return n
	}

	if filter.FolderID != nil {
		cond = append(cond, fmt.Sprintf("i.feed_id in (select id from feeds where folder_id = $%d)", next()))
		args = append(args, *filter.FolderID)
	}
	if filter.FeedID != nil {
		cond = append(cond, fmt.Sprintf("i.feed_id = $%d", next()))
		args = append(args, *filter.FeedID)
	}
	if filter.Status != nil {
		cond = append(cond, fmt.Sprintf("i.status = $%d", next()))
		args = append(args, *filter.Status)
	}
	if filter.Search != nil {
		words := strings.Fields(*filter.Search)
		terms := make([]string, len(words))
		for idx, word := range words {
			terms[idx] = word + ":*"
		}

		cond = append(cond, fmt.Sprintf(
			"i.search @@ to_tsquery('english', $%d)", next(),
		))
		args = append(args, strings.Join(terms, " & "))
	}
	if filter.After != nil {
		compare := ">"
		if newestFirst {
			compare = "<"
		}
		cond = append(cond, fmt.Sprintf(
			"(i.date, i.id) %s (select date, id from items where id = $%d)",
			compare, next(),
		))
		args = append(args, *filter.After)
	}
	if filter.IDs != nil && len(*filter.IDs) > 0 {
		placeholders := make([]string, len(*filter.IDs))
		for i, id := range *filter.IDs {
			placeholders[i] = fmt.Sprintf("$%d", next())
			args = append(args, id)
		}
		cond = append(cond, "i.id in ("+strings.Join(placeholders, ",")+")")
	}
	if filter.SinceID != nil {
		cond = append(cond, fmt.Sprintf("i.id > $%d", next()))
		args = append(args, filter.SinceID)
	}
	if filter.MaxID != nil {
		cond = append(cond, fmt.Sprintf("i.id < $%d", next()))
		args = append(args, filter.MaxID)
	}
	if filter.Before != nil {
		cond = append(cond, fmt.Sprintf("i.date < $%d", next()))
		args = append(args, filter.Before)
	}

	predicate := "1"
	if len(cond) > 0 {
		predicate = strings.Join(cond, " and ")
	}

	return predicate, args
}

func (s *PostgresStorage) CountItems() int {
	var count int
	err := s.db.QueryRow(`select count(*) from items`).Scan(&count)
	if err != nil {
		log.Print(err)
		return 0
	}
	return count
}

func (s *PostgresStorage) ListItems(
	filter model.ItemFilter,
	limit int,
	newestFirst bool,
	withContent bool,
) []model.Item {
	predicate, args := listQueryPredicate(filter, newestFirst)
	result := make([]model.Item, 0)

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
	defer rows.Close()

	for rows.Next() {
		var x model.Item
		err = rows.Scan(
			&x.Id, &x.GUID, &x.FeedId,
			&x.Title, &x.Link, &x.Date,
			&x.Status, (*MediaLinks)(&x.MediaLinks), &x.Content,
		)
		if err != nil {
			log.Print(err)
			return result
		}
		result = append(result, x)
	}
	return result
}

func (s *PostgresStorage) GetItem(id int64) *model.Item {
	i := &model.Item{}
	err := s.db.QueryRow(`
		select
			i.id, i.guid, i.feed_id, i.title, i.link, i.content,
			i.date, i.status, i.media_links
		from items i
		where i.id = $1
	`, id).Scan(
		&i.Id, &i.GUID, &i.FeedId, &i.Title, &i.Link, &i.Content,
		&i.Date, &i.Status, (*MediaLinks)(&i.MediaLinks),
	)
	if err != nil {
		log.Print(err)
		return nil
	}
	return i
}

func (s *PostgresStorage) UpdateItem(id int64, params model.UpdateItemParams) bool {
	sets := make([]string, 0)
	args := make([]any, 0)
	n := 0

	if params.Title != nil {
		n++
		sets = append(sets, fmt.Sprintf("title = $%d", n))
		args = append(args, *params.Title)
		n++
		sets = append(sets, fmt.Sprintf("search = to_tsvector('simple', $%d || ' ' || coalesce((select i2.content from items i2 where i2.id = $%d), ''))", n-1, n))
		args = append(args, id)
	}
	if params.Status != nil {
		n++
		sets = append(sets, fmt.Sprintf("status = $%d", n))
		args = append(args, *params.Status)
	}
	if params.LastArrived != nil {
		n++
		sets = append(sets, fmt.Sprintf("last_arrived = $%d", n))
		args = append(args, *params.LastArrived)
	}
	if len(sets) == 0 {
		return true
	}

	n++
	args = append(args, id)
	query := fmt.Sprintf("update items set %s where id = $%d", strings.Join(sets, ", "), n)
	_, err := s.db.Exec(query, args...)
	return err == nil
}

func (s *PostgresStorage) DeleteItem(id int64) bool {
	_, err := s.db.Exec(`delete from items where id = $1`, id)
	return err == nil
}

func (s *PostgresStorage) UpdateItemStatus(item_id int64, status model.ItemStatus) bool {
	_, err := s.db.Exec(`update items set status = $2 where id = $1`,
		item_id,
		status,
	)
	return err == nil
}

func (s *PostgresStorage) MarkItemsRead(filter model.MarkFilter) bool {
	predicate, args := listQueryPredicate(model.ItemFilter{
		FolderID: filter.FolderID,
		FeedID:   filter.FeedID,
		Before:   filter.Before,
	}, false)
	query := fmt.Sprintf(`
		update items as i set status = %d
		where %s and i.status != %d
		`, model.READ, predicate, model.STARRED)
	_, err := s.db.Exec(query, args...)
	if err != nil {
		log.Print(err)
	}
	return err == nil
}

func (s *PostgresStorage) FeedStats() []model.FeedStat {
	result := make([]model.FeedStat, 0)
	rows, err := s.db.Query(fmt.Sprintf(`
		select
			feed_id,
			sum(case status when %d then 1 else 0 end),
			sum(case status when %d then 1 else 0 end)
		from items
		group by feed_id
	`, model.UNREAD, model.STARRED))
	if err != nil {
		log.Print(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		stat := model.FeedStat{}
		rows.Scan(&stat.FeedId, &stat.UnreadCount, &stat.StarredCount)
		result = append(result, stat)
	}
	return result
}

var (
	itemsKeepSize = 50
	itemsKeepDays = 90
)

func (s *PostgresStorage) DeleteOldItems() {
	keepDaysLimit := fmt.Sprintf("-%d days", itemsKeepDays)
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
				where status != $1
			) sub
			where rn > $2
			  and last_arrived < max_la + $3::interval
		)`,
		model.STARRED,
		itemsKeepSize,
		keepDaysLimit,
	)
	if err != nil {
		log.Print(err)
		return
	}
	numDeleted, err := result.RowsAffected()
	if err == nil && numDeleted > 0 {
		log.Printf("Deleted %d old items", numDeleted)
	}
}
