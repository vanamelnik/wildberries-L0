package postgres

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vanamelnik/wildberries-L0/storage"
)

func TestStore(t *testing.T) {
	defer cleanOrdersTable(t)
	t.Run("Store 10 records", func(t *testing.T) {
		fixtures := make([]storage.OrderDB, 0, 10)
		for i := 1; i <= 10; i++ {
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
			aaaaaaaa
		}
	})

}
