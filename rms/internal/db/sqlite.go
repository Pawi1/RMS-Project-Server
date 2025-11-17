package db

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func OpenSQLite(path string) (*gorm.DB, error) {
    d, err := gorm.Open(sqlite.Open(path), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Warn),
    })
    if err != nil {
        return nil, err
    }
    sqlDB, err := d.DB()
    if err != nil {
        return nil, err
    }
    if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
        return nil, err
    }
    return d, nil
}