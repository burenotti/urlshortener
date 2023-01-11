package storage

import (
	"context"
	"errors"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	linkID     = "abc"
	url        = "https://localhost"
	linkTTL    = time.Hour
	redisError = errors.New("unexpected redis error")
)

func TestRedisStorage_SaveShortLink(t *testing.T) {

	db, mock := redismock.NewClientMock()
	mock.ExpectSetNX(linkID, url, linkTTL).SetVal(true)
	store := NewRedisStorage(db, linkTTL)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	err := store.SaveShortLink(ctx, linkID, url)
	assert.NoError(t, err, "error should be nil")
	assert.NoError(t, mock.ExpectationsWereMet(), "should match expectations")
}

func TestRedisStorage_SaveShortLink_RedisError(t *testing.T) {

	db, mock := redismock.NewClientMock()
	mock.ExpectSetNX(linkID, url, linkTTL).SetErr(redisError)
	store := NewRedisStorage(db, linkTTL)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	err := store.SaveShortLink(ctx, linkID, url)
	assert.ErrorIs(t, err, redisError, "should return error if something go wrong")
	assert.NoError(t, mock.ExpectationsWereMet(), "should match expectations")
}

func TestRedisStorage_GetSource(t *testing.T) {
	linkID := "abc"
	url := "https://localhost"
	linkTTL := time.Hour

	db, mock := redismock.NewClientMock()
	mock.ExpectGet(linkID).SetVal(url)
	store := NewRedisStorage(db, linkTTL)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	val, err := store.GetSource(ctx, linkID)

	assert.Equal(t, url, val)
	assert.ErrorIs(t, err, nil, "error should be nil")
	assert.NoError(t, mock.ExpectationsWereMet(), "should match expectations")
}

func TestRedisStorage_GetSource_LinkIDNotFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectGet(linkID).SetErr(redisError)
	store := NewRedisStorage(db, linkTTL)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	val, err := store.GetSource(ctx, linkID)

	assert.Equal(t, "", val)
	assert.ErrorIs(t, err, ErrNoSourceUrl, "should return ErrNoSourceUrl if redis doesn't contain given linkID")
	assert.NoError(t, mock.ExpectationsWereMet(), "should match expectations")
}
