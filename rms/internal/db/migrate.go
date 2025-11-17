package db

import (
    "errors"
    "io/fs"
    "os"
    "path/filepath"

    "gorm.io/gorm"
)

func fileExists(p string) bool {
    _, err := os.Stat(p)
    return err == nil
}

func ensureDir(path string) error {
    dir := filepath.Dir(path)
    if dir == "." || dir == "" {
        return nil
    }
    return os.MkdirAll(dir, 0o755)
}

func ApplySchema(db *gorm.DB, schemaPath string) error {
    if schemaPath == "" {
        return errors.New("schema path is empty")
    }
    b, err := os.ReadFile(schemaPath)
    if err != nil {
        return err
    }
    return db.Exec(string(b)).Error
}

func EnsureDatabase(dbPath, schemaPath string, open func(string) (*gorm.DB, error)) (*gorm.DB, bool, error) {
    created := false
    if !fileExists(dbPath) {
        if err := ensureDir(dbPath); err != nil {
            return nil, created, err
        }
        created = true
    }
    gdb, err := open(dbPath)
    if err != nil {
        return nil, created, err
    }
    if created {
        if err := ApplySchema(gdb, schemaPath); err != nil {
            _ = os.Remove(dbPath)
            return nil, created, err
        }
    }
    return gdb, created, nil
}

func ReapplySchema(db *gorm.DB, schemaPath string) error {
    _, err := os.Stat(schemaPath)
    if err != nil && !errors.Is(err, fs.ErrNotExist) {
        return err
    }
    return ApplySchema(db, schemaPath)
}