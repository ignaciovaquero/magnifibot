package archimadrid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetGospelFromCache(t *testing.T) {
	tests := []struct {
		name,
		key string
		cacheObject   interface{}
		expected      *Gospel
		errorExpected bool
	}{
		{
			name: "Valid gospel from Cache",
			key:  "key",
			cacheObject: &Gospel{
				Day: "today",
			},
			expected: &Gospel{
				Day: "today",
			},
			errorExpected: false,
		},
		{
			name: "Invalid object from Cache",
			key:  "key",
			cacheObject: Gospel{
				Day: "today",
			},
			errorExpected: true,
		},
		{
			name:          "No object from Cache",
			key:           "key",
			errorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			client := NewClient()
			if test.cacheObject != nil {
				client.Set(test.key, test.cacheObject)
			}
			actual, err := client.getGospelFromCache(test.key)
			if test.errorExpected {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
			assert.ObjectsAreEqual(test.expected, actual)
		})
	}
}

func TestGetResponseFromCache(t *testing.T) {
	tests := []struct {
		name,
		key string
		cacheObject   interface{}
		expected      *gospelResponse
		errorExpected bool
	}{
		{
			name: "Valid response from cache",
			key:  "key",
			cacheObject: &gospelResponse{
				PostTitle: "title",
			},
			expected: &gospelResponse{
				PostTitle: "title",
			},
			errorExpected: false,
		},
		{
			name: "Invalid response from cache",
			key:  "key",
			cacheObject: gospelResponse{
				PostTitle: "title",
			},
			errorExpected: true,
		},
		{
			name:          "No response from cache",
			key:           "key",
			errorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			client := NewClient()
			if test.cacheObject != nil {
				client.Set(test.key, test.cacheObject)
			}
			actual, err := client.getResponseFromCache(test.key)
			if test.errorExpected {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
			assert.ObjectsAreEqual(test.expected, actual)
		})
	}
}

func TestGetGospelOrLecture(t *testing.T) {
	tests := []struct {
		name string
		day  time.Time
		regexString,
		cachePrefix string
		psalm         bool
		key           string
		cache         interface{}
		response      string
		code          int
		expected      *Gospel
		errorExpected bool
	}{
		{},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.code)
				w.Write([]byte(test.response))
			}))
			defer server.Close()

			client := NewClient(SetURL(server.URL))

			if test.cache != nil {
				client.Set(test.key, test.cache)
			}
			actual, err := client.getGospelOrLecture(
				context.TODO(),
				test.day,
				test.regexString,
				test.cachePrefix,
				test.psalm,
			)
			if test.errorExpected {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
			assert.ObjectsAreEqual(test.expected, actual)

			today := test.day.Format("2006-01-02")
			object, err := client.Get(ResponsePrefix + today)
			_, ok := object.(*gospelResponse)
			assert.NoError(tt, err)
			assert.True(tt, ok)

			object, err = client.Get(test.cachePrefix + today)
			actual, ok = object.(*Gospel)
			assert.NoError(tt, err)
			assert.True(tt, ok)
			assert.ObjectsAreEqual(test.expected, actual)
		})
	}
}
