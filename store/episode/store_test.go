package episode

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"pcast-api/db"
	"strconv"
	"testing"
)

var d *gorm.DB
var es *Store

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setup() {
	d = db.NewTestDB("./../../fixtures/test/pcast.db")
	es = New(d)
}

func tearDown() {
	es.RemoveTable()
}

func newEpisode() *Episode {
	id, _ := uuid.NewV7()
	guid := strconv.Itoa(rand.Intn(9999999999))

	return &Episode{FeedId: id, FeedGUID: fmt.Sprintf("tag:soundcloud,2010:tracks/%s", guid)}
}

func TestCreateEpisode(t *testing.T) {
	episode := newEpisode()
	err := es.Create(episode)
	assert.NoError(t, err)

	es.TruncateTables()
}

func TestFindEpisodeByID(t *testing.T) {
	episode := newEpisode()
	err := es.Create(episode)
	assert.NoError(t, err)

	foundEpisode, err := es.FindByID(episode.ID)

	assert.NoError(t, err)
	assert.Equal(t, episode.FeedId, foundEpisode.FeedId)

	es.TruncateTables()
}

func TestDeleteEpisode(t *testing.T) {
	episode := newEpisode()
	err := es.Create(episode)
	assert.NoError(t, err)

	err = es.Delete(episode)
	assert.NoError(t, err)

	es.TruncateTables()
}

func TestUpdateEpisode(t *testing.T) {
	episode := newEpisode()
	err := es.Create(episode)
	assert.NoError(t, err)

	episode.FeedId, _ = uuid.NewV7()
	err = es.Update(episode)
	assert.NoError(t, err)

	foundEpisode, err := es.FindByID(episode.ID)
	assert.NoError(t, err)
	assert.Equal(t, episode.FeedId, foundEpisode.FeedId)

	es.TruncateTables()
}