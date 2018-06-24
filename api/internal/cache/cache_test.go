package cache

import (
	"testing"
	"github.com/gambarini/cabapi/api/internal/model"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/gambarini/cabapi/tstutils"
)

func TestCache_Set(t *testing.T) {

	err := tstutils.StartRedis()
	defer tstutils.StopRedis()

	if err != nil {
		t.Errorf("Failed to start Redis, %s", err)
		return
	}

	cache, err := NewCache()

	if err != nil {
		t.Errorf("Failed to create cache, %s", err)
		return
	}

	defer cache.Close()

	expected := model.Trips{
		Medallion: "M",
		Date:      "2010-10-10",
		Total:     10,
	}

	err = cache.Set(expected)

	assert.Nil(t, err)

	actual, err := cache.Get("M", time.Date(2010, 10, 10, 0, 0, 0, 0, time.Local))

	assert.Nil(t, err)

	assert.EqualValues(t, expected.Medallion, actual.Medallion)
	assert.EqualValues(t, expected.Total, actual.Total)
	assert.EqualValues(t, expected.Date, actual.Date)

}
