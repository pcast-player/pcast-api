package episode

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"pcast-api/db"
)

var d *sql.DB
var es *Store

const testDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setup() {
	d = db.NewTestDBSQL(testDSN)

	// Run migrations
	runMigrations()

	es = New(d)
}

func tearDown() {
	// Clean up test data
	truncateTable()
	d.Close()
}

func runMigrations() {
	// Create episodes table if not exists
	// Split statements to avoid race conditions in parallel tests
	d.Exec(`
		CREATE TABLE IF NOT EXISTS episodes (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			feed_id UUID NOT NULL,
			feed_guid VARCHAR(255) NOT NULL,
			current_position INTEGER,
			played BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)
	d.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_id ON episodes(feed_id)`)
	d.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_guid ON episodes(feed_guid)`)
}

func truncateTable() {
	_, err := d.Exec("TRUNCATE TABLE episodes")
	if err != nil {
		// Table might not exist yet, ignore error
		return
	}
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

	truncateTable()
}

func TestFindEpisodeByID(t *testing.T) {
	episode := newEpisode()
	err := es.Create(episode)
	assert.NoError(t, err)

	foundEpisode, err := es.FindByID(episode.ID)

	assert.NoError(t, err)
	assert.Equal(t, episode.FeedId, foundEpisode.FeedId)

	truncateTable()
}

func TestDeleteEpisode(t *testing.T) {
	episode := newEpisode()
	err := es.Create(episode)
	assert.NoError(t, err)

	err = es.Delete(episode)
	assert.NoError(t, err)

	truncateTable()
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

	truncateTable()
}
