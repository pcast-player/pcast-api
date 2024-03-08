package user

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"os"
	"pcast-api/db"
	"pcast-api/helper"
	"testing"
)

var d *gorm.DB
var us *Store

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setup() {
	d = db.NewTestDB("./../../fixtures/test/store_user.db")
	us = New(d)
}

func tearDown() {
	helper.RemoveTable(d, "users")
}

func truncateTable() {
	helper.TruncateTables(d, "users")
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
