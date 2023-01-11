package storage

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestComposedShortener_SaveShortLink(t *testing.T) {
	cases := []struct {
		name          string
		description   string
		mainStoreErr  error
		cacheStoreErr error
		expectedErr   error
	}{
		{
			"correct",
			"if all main and cache saved correctly should return nil error",
			nil,
			nil,
			nil,
		},
		{
			"save_failed",
			"if both storages failed should return ErrSaveFailed",
			ErrSaveFailed,
			ErrSaveFailed,
			ErrSaveFailed,
		},
		{
			"save_failed",
			"if the main storage failed should return ErrSaveFailed even if cache worked well",
			ErrSaveFailed,
			nil,
			ErrSaveFailed,
		},
		{
			"partial_save",
			"if cache storage failed should return ErrPartialSave",
			nil,
			ErrSaveFailed,
			ErrPartialSave,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			cacheMock := new(ShortenerStorageMock)
			cacheMock.On("SaveShortLink", mock.Anything, mock.Anything).Return(c.cacheStoreErr)

			mainMock := new(ShortenerStorageMock)
			mainMock.On("SaveShortLink", mock.Anything, mock.Anything).Return(c.mainStoreErr)

			store := NewComposedShortener(mainMock, cacheMock)

			err := store.SaveShortLink(ctx, "abc", "https://localhost:80/")
			assert.ErrorIs(t, err, c.expectedErr,
				"if all main and cache saved correctly should return no error")
			mainMock.MethodCalled("SaveShortLink", "abc", "https://localhost:80/")
			cacheMock.MethodCalled("SaveShortLink", "abc", "https://localhost:80/")
			mainMock.AssertExpectations(t)
			cacheMock.AssertExpectations(t)
		})
	}

}

func TestComposedShortener_GetSource(t *testing.T) {
	cases := []struct {
		name                   string
		errMsg                 string
		valMsg                 string
		linkID                 string
		expectedURL            string
		expectedErr            error
		mainStorageErr         error
		cacheStorageErr        error
		shouldCallMainStorage  bool
		shouldCallCacheStorage bool
	}{
		{
			name:                   "correct",
			valMsg:                 "should return correct value",
			errMsg:                 "should return no error on correct return of both caches",
			linkID:                 "abc",
			expectedURL:            "https://localhost",
			expectedErr:            nil,
			mainStorageErr:         nil,
			cacheStorageErr:        nil,
			shouldCallMainStorage:  false,
			shouldCallCacheStorage: true,
		},
		{
			name:                   "cache_no_source_found",
			valMsg:                 "should fetch correct value from main storage",
			errMsg:                 "should return no error",
			linkID:                 "abc",
			expectedURL:            "https://localhost",
			expectedErr:            nil,
			mainStorageErr:         nil,
			cacheStorageErr:        ErrNoSourceUrl,
			shouldCallCacheStorage: true,
			shouldCallMainStorage:  true,
		},
		{
			name:                   "cache_error",
			valMsg:                 "should fetch correct value from main storage",
			errMsg:                 "should return no error",
			linkID:                 "abc",
			expectedURL:            "https://localhost",
			expectedErr:            nil,
			mainStorageErr:         nil,
			cacheStorageErr:        errors.New("test unexpected error"),
			shouldCallCacheStorage: true,
			shouldCallMainStorage:  true,
		},
		{
			name:                   "source_url_not_found",
			valMsg:                 "can't fetch source url from storage",
			errMsg:                 "should return ErrNoSourceUrl",
			linkID:                 "abc",
			expectedURL:            "",
			expectedErr:            ErrNoSourceUrl,
			mainStorageErr:         ErrNoSourceUrl,
			cacheStorageErr:        ErrNoSourceUrl,
			shouldCallCacheStorage: true,
			shouldCallMainStorage:  true,
		},
	}

	for _, c := range cases {

		t.Run(c.name, func(t *testing.T) {
			cacheMock := new(ShortenerStorageMock)
			if c.shouldCallCacheStorage {
				cacheMock.On("GetSource", c.linkID).Return(c.expectedURL, c.cacheStorageErr)
			}

			mainMock := new(ShortenerStorageMock)
			if c.shouldCallMainStorage {
				mainMock.On("GetSource", c.linkID).Return(c.expectedURL, c.mainStorageErr)
			}

			store := NewComposedShortener(mainMock, cacheMock)
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()
			linkID, err := store.GetSource(ctx, c.linkID)
			assert.Equal(t, c.expectedURL, linkID, c.valMsg)
			assert.ErrorIs(t, err, c.expectedErr, c.errMsg)
			cacheMock.AssertExpectations(t)
			mainMock.AssertExpectations(t)
		})
	}
}
