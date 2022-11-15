package cache

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
)

type Cache struct {
	cache *lru.Cache
}

func New(capacity int) (*Cache, error) {
	cache, err := lru.New(capacity)
	if err != nil {
		return nil, err
	}
	return &Cache{
		cache: cache,
	}, nil
}

type Bucket map[string]*domain.Message

func (c *Cache) getBucket(userId int64) Bucket {
	bucket, ok := c.cache.Get(userId)
	if !ok {
		bucket = make(Bucket)
		c.cache.Add(userId, bucket)
	}
	return bucket.(Bucket)
}

func (c *Cache) Add(userId int64, signature string, response *domain.Message) {
	bucket := c.getBucket(userId)
	bucket[signature] = response
}

func (c *Cache) Contains(userId int64, signature string) bool {
	bucket, ok := c.cache.Get(userId)
	if !ok {
		return false
	}
	_, found := bucket.(Bucket)[signature]
	return found
}

func (c *Cache) Get(userId int64, signature string) (*domain.Message, bool) {
	bucket := c.getBucket(userId)
	response, ok := bucket[signature]
	return response, ok
}

func (c *Cache) Clear(userId int64) {
	c.cache.Remove(userId)
}

func (c *Cache) WithCache(cmd domain.CommandHandler) domain.CommandHandler {
	return func(msg *domain.Message) *domain.Message {
		var response *domain.Message

		span, ctx := opentracing.StartSpanFromContext(
			msg.Context,
			fmt.Sprintf("cache for user %d", msg.UserID),
		)
		defer span.Finish()
		msg.Context = ctx

		if !c.Contains(msg.UserID, msg.Command) {
			logger.Debug("cache miss")
			response = cmd(msg)
			c.Add(msg.UserID, msg.Command, response)
		} else {
			logger.Debug("cache hit")
		}

		response, ok := c.Get(msg.UserID, msg.Command)
		if !ok {
			response = &domain.Message{
				UserID: msg.UserID,
				Text:   "Cache error",
			}
		}
		return response
	}
}

func (c *Cache) Invalidate(cmd domain.CommandHandler) domain.CommandHandler {
	return func(msg *domain.Message) *domain.Message {
		logger.Debugf("Invalidating cache for user %d", msg.UserID)
		c.Clear(msg.UserID)
		return cmd(msg)
	}
}
