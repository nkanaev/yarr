package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/nkanaev/yarr/src/server/router"
	"github.com/nkanaev/yarr/src/storage"
)

type reactionPayload struct {
	ItemID   int64  `json:"item_id"`
	Reaction string `json:"reaction"`
}

type clickThroughPayload struct {
	ItemID int64 `json:"item_id"`
}

type readHerePayload struct {
	ItemID int64 `json:"item_id"`
}

func (s *Server) handleReactions(c *router.Context) {
	switch c.Req.Method {
	case "GET":
		itemID, err := c.QueryInt64("item_id")
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "item_id required"})
			return
		}
		reaction, err := s.db.GetReaction(itemID)
		if err != nil {
			log.Print("GetReaction error:", err)
			c.Out.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, map[string]string{"reaction": reaction})

	case "POST":
		var body reactionPayload
		if err := json.NewDecoder(c.Req.Body).Decode(&body); err != nil {
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.ItemID == 0 {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "item_id required"})
			return
		}
		if err := s.db.SetReaction(body.ItemID, body.Reaction); err != nil {
			log.Print("SetReaction error:", err)
			c.Out.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.Out.WriteHeader(http.StatusOK)

	default:
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleClickThroughs(c *router.Context) {
	if c.Req.Method != "POST" {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body clickThroughPayload
	if err := json.NewDecoder(c.Req.Body).Decode(&body); err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}
	if body.ItemID == 0 {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "item_id required"})
		return
	}
	if err := s.db.LogClickThrough(body.ItemID); err != nil {
		log.Print("LogClickThrough error:", err)
		c.Out.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Out.WriteHeader(http.StatusOK)
}

func (s *Server) handleRankedItems(c *router.Context) {
	if c.Req.Method != "GET" {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	limit := 20
	offset := 0
	query := c.Req.URL.Query()

	if pageStr := query.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 1 {
			offset = (page - 1) * limit
		}
	}
	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	filter := storage.ItemFilter{}
	if folderID, err := c.QueryInt64("folder_id"); err == nil {
		filter.FolderID = &folderID
	}
	if feedID, err := c.QueryInt64("feed_id"); err == nil {
		filter.FeedID = &feedID
	}
	if status := query.Get("status"); len(status) != 0 {
		statusValue := storage.StatusValues[status]
		filter.Status = &statusValue
	}

	items, hasMore, err := s.db.GetRankedItems(filter, limit, offset)
	if err != nil {
		log.Print("GetRankedItems error:", err)
		c.Out.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"list":     items,
		"has_more": hasMore,
	})
}

func (s *Server) handleReadHeres(c *router.Context) {
	if c.Req.Method != "POST" {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body readHerePayload
	if err := json.NewDecoder(c.Req.Body).Decode(&body); err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}
	if body.ItemID == 0 {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "item_id required"})
		return
	}
	if err := s.db.LogReadHere(body.ItemID); err != nil {
		log.Print("LogReadHere error:", err)
		c.Out.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Out.WriteHeader(http.StatusOK)
}

func (s *Server) handlePreferences(c *router.Context) {
	switch c.Req.Method {
	case "GET":
		stats, err := s.db.GetPreferenceStats()
		if err != nil {
			log.Print("GetPreferenceStats error:", err)
			c.Out.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, stats)

	case "DELETE":
		if err := s.db.DeleteAllPreferences(); err != nil {
			log.Print("DeleteAllPreferences error:", err)
			c.Out.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.Out.WriteHeader(http.StatusOK)

	default:
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}
