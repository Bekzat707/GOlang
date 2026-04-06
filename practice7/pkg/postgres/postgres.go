package postgres

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Postgres struct {
    Conn *gorm.DB
}

func New(dsn string) (*Postgres, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return &Postgres{Conn: db}, nil
}
