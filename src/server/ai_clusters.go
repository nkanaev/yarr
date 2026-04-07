package server

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/nkanaev/yarr/src/server/router"
	"github.com/nkanaev/yarr/src/storage"
)

// clusterRunPayload is the JSON body for POST /api/ai/clusters.
type clusterRunPayload struct {
	GeneratedAt string                 `json:"generated_at"`
	Algorithm   string                 `json:"algorithm"`
	NArticles   int                    `json:"n_articles"`
	NNoise      int                    `json:"n_noise"`
	Clusters    []storage.ClusterLabel `json:"clusters"`
	Centroids   []centroidPayload      `json:"centroids"`
}

// centroidPayload carries centroid data from Python (base64-encoded float64 bytes).
type centroidPayload struct {
	ClusterID int    `json:"cluster_id"`
	Label     string `json:"label"`
	Centroid  string `json:"centroid"` // base64-encoded raw float64 bytes
}

func (s *Server) handleAiClusters(c *router.Context) {
	if c.Req.Method == "GET" {
		summary, err := s.db.GetClusterSummary()
		if err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusInternalServerError)
			return
		}
		if summary == nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"clusters":     []interface{}{},
				"n_clusters":   0,
				"n_articles":   0,
				"generated_at": "",
				"message":      "No clusters yet. POST /api/ai/recluster to generate.",
			})
			return
		}
		c.JSON(http.StatusOK, summary)
	} else if c.Req.Method == "POST" {
		var payload clusterRunPayload
		if err := json.NewDecoder(c.Req.Body).Decode(&payload); err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}

		algorithm := payload.Algorithm
		if algorithm == "" {
			algorithm = "hdbscan"
		}

		// Decode base64 centroid blobs
		centroids := make([]storage.ClusterCentroid, 0, len(payload.Centroids))
		for _, cp := range payload.Centroids {
			blob, err := base64.StdEncoding.DecodeString(cp.Centroid)
			if err != nil {
				log.Printf("invalid centroid base64 for cluster %d: %v", cp.ClusterID, err)
				c.Out.WriteHeader(http.StatusBadRequest)
				return
			}
			centroids = append(centroids, storage.ClusterCentroid{
				ClusterID: cp.ClusterID,
				Label:     cp.Label,
				Centroid:  blob,
			})
		}

		_, err := s.db.SaveClusterRun(
			payload.GeneratedAt,
			algorithm,
			payload.NArticles,
			payload.NNoise,
			payload.Clusters,
			centroids,
		)
		if err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.Out.WriteHeader(http.StatusCreated)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleAiClusterCentroids(c *router.Context) {
	if c.Req.Method != "GET" {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	centroids, err := s.db.GetClusterCentroids()
	if err != nil {
		log.Print(err)
		c.Out.WriteHeader(http.StatusInternalServerError)
		return
	}
	if centroids == nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// Encode centroid blobs as base64 for JSON transport
	type centroidResponse struct {
		ClusterID int    `json:"cluster_id"`
		Label     string `json:"label"`
		Centroid  string `json:"centroid"` // base64-encoded
	}
	result := make([]centroidResponse, 0, len(centroids))
	for _, c2 := range centroids {
		result = append(result, centroidResponse{
			ClusterID: c2.ClusterID,
			Label:     c2.Label,
			Centroid:  base64.StdEncoding.EncodeToString(c2.Centroid),
		})
	}
	c.JSON(http.StatusOK, result)
}
