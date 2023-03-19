package adapter

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

var enforcer *casbin.SyncedEnforcer

func InitializeCasbin(db *sqlx.DB) error {
	var err error
	pgSchema := os.Getenv("DB_SCHEMA")
	a := NewCasbinAdapter(db, pgSchema)
	enforcer, err = casbin.NewSyncedEnforcer("./model.conf", a)
	if err != nil {
		return fmt.Errorf("casbin initialization failed. Err:%w", err)
	}

	enforcer.StartAutoLoadPolicy(time.Second * 5)
	return nil
}

func Policy() *casbin.SyncedEnforcer {
	if enforcer == nil {
		panic(errors.New("casbin not initalized"))
	}
	return enforcer
}
