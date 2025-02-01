package ai_service

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	cache_service "github.com/wapikit/wapikit/services/redis_service"
)

type RedisVectorCache struct {
	client   *cache_service.RedisClient
	prefix   string
	indexKey string
}

// NewRedisVectorCache initializes a new RedisVectorCache.
func NewRedisVectorCache(client *cache_service.RedisClient) *RedisVectorCache {
	return &RedisVectorCache{
		client:   client,
		prefix:   "vector_cache:",
		indexKey: "vector_cache:keys",
	}
}

// StoreQueryEmbedding stores the embedding for a query in Redis.
// The embedding is stored as a JSON string under key "vector_cache:<query>".
func (vc *RedisVectorCache) StoreQueryEmbedding(query string, embedding []float64, ttl time.Duration) error {
	data, err := json.Marshal(embedding)
	if err != nil {
		return err
	}
	key := vc.prefix + query
	// Cache the embedding.
	if err = vc.client.CacheData(key, string(data), ttl); err != nil {
		return err
	}
	// Also add the query to an index set so we can later iterate over all embeddings.
	ctx := context.Background()
	return vc.client.Client.SAdd(ctx, vc.indexKey, query).Err()
}

// GetEmbedding retrieves the stored embedding for a given query.
// It returns the embedding, a boolean indicating whether it was found, and an error if any.
func (vc *RedisVectorCache) GetEmbedding(query string) ([]float64, bool, error) {
	key := vc.prefix + query
	data, err := vc.client.GetCachedData(key)
	if err != nil {
		// If key not found, redis returns a redis.Nil error.
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, err
	}
	if data == "" {
		return nil, false, nil
	}
	var embedding []float64
	if err := json.Unmarshal([]byte(data), &embedding); err != nil {
		return nil, false, err
	}
	return embedding, true, nil
}

// FindSimilarEmbeddings performs a simple similarity search.
// It iterates over all keys in the index set, computes cosine similarity with the target embedding,
// and returns the queries whose embeddings have a cosine similarity above the threshold.
func (vc *RedisVectorCache) FindSimilarEmbeddings(target []float64, threshold float64) ([]string, error) {
	ctx := context.Background()
	queries, err := vc.client.Client.SMembers(ctx, vc.indexKey).Result()
	if err != nil {
		return nil, err
	}
	var similarQueries []string
	for _, q := range queries {
		emb, found, err := vc.GetEmbedding(q)
		if err != nil || !found {
			continue
		}
		sim := cosineSimilarity(target, emb)
		if sim >= threshold {
			similarQueries = append(similarQueries, q)
		}
	}
	return similarQueries, nil
}

// cosineSimilarity computes the cosine similarity between two vectors.
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
