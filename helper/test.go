package helper

import (
	"fmt"
	"gorm.io/gorm"
)

func TruncateTables(db *gorm.DB, table string) {
	query := fmt.Sprintf("DELETE FROM %s;", table)

	err := db.Exec(query).Error
	if err != nil {
		panic(err)
	}
}

func RemoveTable(db *gorm.DB, model interface{}) {
	err := db.Migrator().DropTable(model)
	if err != nil {
		panic(err)
	}
}
