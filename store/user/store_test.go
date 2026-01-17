package user

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"pcast-api/db"
)

var d *sql.DB
var us *Store

const testDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

// Helper function to create a pointer to a string
func strPtr(s string) *string {
	return &s
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setup() {
	d = db.NewTestDB(testDSN)

	// Run migrations
	runMigrations()

	us = New(d)
}

func tearDown() {
	// Clean up test data
	truncateTable()
	d.Close()
}

func runMigrations() {
	// Create users table if not exists
	// Split statements to avoid race conditions in parallel tests
	d.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255),
			google_id VARCHAR(255) UNIQUE
		)
	`)
	d.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
	d.Exec(`CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id)`)
}

func truncateTable() {
	_, err := d.Exec("TRUNCATE TABLE users CASCADE")
	if err != nil {
		// Table might not exist yet, ignore error
		return
	}
}

func TestCreateUser(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: strPtr("password")}
	err := us.Create(context.Background(), user)
	assert.NoError(t, err)

	truncateTable()
}

func TestFindUserByID(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: strPtr("password")}
	err := us.Create(context.Background(), user)
	assert.NoError(t, err)

	foundUser, err := us.FindByID(context.Background(), user.ID)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)

	truncateTable()
}

func TestFindUserByEmail(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: strPtr("password")}
	err := us.Create(context.Background(), user)
	assert.NoError(t, err)

	foundUser, err := us.FindByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)

	truncateTable()
}

func TestDeleteUser(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: strPtr("password")}
	err := us.Create(context.Background(), user)
	assert.NoError(t, err)

	err = us.Delete(context.Background(), user)
	assert.NoError(t, err)

	truncateTable()
}

func TestUpdateUser(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: strPtr("password")}
	err := us.Create(context.Background(), user)

	assert.NoError(t, err)
	user.Email = "bar@foo.com"
	err = us.Update(context.Background(), user)
	assert.NoError(t, err)

	foundUser, err := us.FindByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)

	truncateTable()
}

func TestCreateOAuthUser(t *testing.T) {
	user := &User{Email: "oauth@bar.com", GoogleID: strPtr("google123")}
	err := us.CreateOAuthUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Nil(t, user.Password)

	truncateTable()
}

func TestFindUserByGoogleID(t *testing.T) {
	googleID := "google123"
	user := &User{Email: "oauth@bar.com", GoogleID: &googleID}
	err := us.CreateOAuthUser(context.Background(), user)
	assert.NoError(t, err)

	foundUser, err := us.FindByGoogleID(context.Background(), googleID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, googleID, *foundUser.GoogleID)

	truncateTable()
}

func TestUpdateGoogleID(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: strPtr("password")}
	err := us.Create(context.Background(), user)
	assert.NoError(t, err)

	googleID := "google456"
	err = us.UpdateGoogleID(context.Background(), user.ID, googleID)
	assert.NoError(t, err)

	foundUser, err := us.FindByGoogleID(context.Background(), googleID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, googleID, *foundUser.GoogleID)

	truncateTable()
}
