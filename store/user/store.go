package user

import "gorm.io/gorm"

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	autoMigrateUserModel(db)

	return &Store{db: db}
}

func autoMigrateUserModel(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		panic("Failed to migrate database!")
	}
}

func (s *Store) FindAll() ([]User, error) {
	var users []User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Store) FindByID(id string) (*User, error) {
	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) Create(user *User) error {
	return s.db.Create(user).Error
}

func (s *Store) Update(user *User) error {
	return s.db.Save(user).Error
}

func (s *Store) Delete(user *User) error {
	return s.db.Delete(user).Error
}

func (s *Store) TruncateTables() {
	err := s.db.Exec("DELETE FROM users;").Error
	if err != nil {
		panic(err)
	}
}

func (s *Store) RemoveTables() {
	err := s.db.Migrator().DropTable(&User{})
	if err != nil {
		panic(err)
	}
}
