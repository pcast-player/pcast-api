package user

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"pcast-api/db"
)

var d *sql.DB
var us *Store

const testDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

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
			password VARCHAR(255) NOT NULL
		)
	`)
	d.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
}

func truncateTable() {
	_, err := d.Exec("TRUNCATE TABLE users CASCADE")
	if err != nil {
		// Table might not exist yet, ignore error
		return
	}
}

func TestCreateUser(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: "password"}
	err := us.Create(user)
	assert.NoError(t, err)

	truncateTable()
}

func TestFindUserByID(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: "password"}
	err := us.Create(user)
	assert.NoError(t, err)

	foundUser, err := us.FindByID(user.ID)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)

	truncateTable()
}

func TestFindUserByEmail(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: "password"}
	err := us.Create(user)
	assert.NoError(t, err)

	foundUser, err := us.FindByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)

	truncateTable()
}

func TestDeleteUser(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: "password"}
	err := us.Create(user)
	assert.NoError(t, err)

	err = us.Delete(user)
	assert.NoError(t, err)

	truncateTable()
}

func TestUpdateUser(t *testing.T) {
	user := &User{Email: "foo@bar.com", Password: "password"}
	err := us.Create(user)

	assert.NoError(t, err)
	user.Email = "bar@foo.com"
	err = us.Update(user)
	assert.NoError(t, err)

	foundUser, err := us.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)

	truncateTable()
}
