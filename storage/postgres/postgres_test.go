package postgres

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vanamelnik/wildberries-L0/storage"
)

func TestPostgres(t *testing.T) {
	defer cleanOrdersTable(t)
	fixtures := make([]storage.OrderDB, 0, 10)
	for i := 1; i <= 1000; i++ {
		order := struct {
			ID   int
			Data string
		}{
			ID:   i,
			Data: fmt.Sprintf("%d*%d=%d", i, i, i*i),
		}
		jsonOrder, err := json.Marshal(order)
		require.NoError(t, err)
		fixtures = append(fixtures, storage.OrderDB{
			OrderUID:  fmt.Sprint(i),
			JSONOrder: string(jsonOrder),
		})
	}
	t.Run("Store 1000 records", func(t *testing.T) {
		for _, o := range fixtures {
			assert.NoError(t, pgMockStorage.Store(o.OrderUID, o.JSONOrder))
		}
		time.Sleep(50 * time.Millisecond)
	})
	t.Run("Get all", func(t *testing.T) {
		got, err := pgMockStorage.GetAll()
		assert.NoError(t, err)
		assert.Equal(t, 1000, len(got))
		for i := range fixtures {
			assert.Equal(t, fixtures[i].OrderUID, got[i].OrderUID)
			assert.JSONEq(t, fixtures[i].JSONOrder, got[i].JSONOrder)
		}
	})
	t.Run("Test Get()", func(t *testing.T) {
		for i := range fixtures {
			order, err := pgMockStorage.Get(fmt.Sprint(i + 1))
			assert.NoError(t, err)
			assert.JSONEq(t, fixtures[i].JSONOrder, order)
		}
	})
}

func TestGetError(t *testing.T) {
	t.Run("Get non-existing order", func(t *testing.T) {
		_, err := pgMockStorage.Get("nihil")
		assert.ErrorIs(t, err, storage.ErrNotFound)
	})
}
