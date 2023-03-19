package record

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRecords_Sort(t *testing.T) {
	recs := Records{
		{
			Weight: 50,
		}, {
			Weight: 100,
		}, {
			Weight: 1000,
		}, {
			Weight: 1,
		},
	}
	recs.Sort()
	var weights []uint64
	for _, r := range recs {
		weights = append(weights, r.Weight)
	}
	assert.ElementsMatch(t, weights, []uint64{1000, 100, 50, 1})
}

func TestRecords_Search(t *testing.T) {
	t.Run("all records should be removed", func(t *testing.T) {
		records := Records{{
			Id:      0,
			Command: "some command OA",
		}, {
			Id:      1,
			Command: "cw tail /kosmos/fargate/market-status-updater -b 5m -f",
		}}
		sorted := records.Search("y")
		assert.Empty(t, sorted)
	})

	t.Run("should prefer the newer record", func(t *testing.T) {
		records := Records{{
			Id:         0,
			LastUpdate: time.Now().Add(-5 * time.Second),
			Command:    "some command OA",
		}, {
			Id:         1,
			LastUpdate: time.Now().Add(-1 * time.Minute),
			Command:    "cw tail /kosmos/fargate/market-status-updater -b 5m -f",
		}}
		weighted := records.Search("s")
		assert.True(t, weighted[0].Command == records[0].Command)
	})

	t.Run("should prefer the record with the highest count", func(t *testing.T) {
		records := Records{{
			Id:      0,
			Count:   2,
			Command: "command 1",
		}, {
			Id:      1,
			Count:   100,
			Command: "command 2",
		}, {
			Id:      2,
			Count:   10,
			Command: "command 3",
		}}
		weighted := records.Search("command")
		var ids []int64
		for _, w := range weighted {
			ids = append(ids, w.Id)
		}
		assert.ElementsMatch(t, ids, []int64{1, 2, 0})
	})

	t.Run("should weight 0 if an letter does not occur", func(t *testing.T) {
		records := Records{{
			Id:      0,
			Command: "git push",
		}}
		weighted := records.Search("exp")
		assert.Empty(t, weighted)
	})
}
