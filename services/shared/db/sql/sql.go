package main

import (
    "github.com/go-pg/pg/v10"
    "github.com/go-pg/pg/v10/orm"
    "github.com/mats9693/unnamed_plan/services/shared/db/model"
    "log"
    "reflect"
    "time"
)

var dbConf = &pg.Options{
    Addr:         "106.75.174.54:5432",
    User:         "mario_test",
    Password:     "123456",
    Database:     "unnamed_plan_test",
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 5 * time.Second,
}

func main() {
    db := pg.Connect(dbConf)
    defer func() {
        _ = db.Close()
    }()

    models := []interface{}{
        (*model.Administrator)(nil),
        (*model.Version)(nil),
        (*model.Service)(nil),
        (*model.Config)(nil),
        (*model.User)(nil),
        (*model.CloudFile)(nil),
        (*model.ThinkingNote)(nil),
    }

    for i, m := range models {
        err := db.Model(m).CreateTable(&orm.CreateTableOptions{Temp: false, IfNotExists: true})
        if err != nil {
            log.Fatalf("create table index: %d, name: %s failed, error: %v\n", i, reflect.TypeOf(m), err)
        }
    }

    // built-in data
    {
        setDefaultUser(db)
        setConfigCenterUser(db)
    }

    log.Println("db init finish.")
}
