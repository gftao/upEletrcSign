package dbModels

import (
	"time"
	"github.com/jinzhu/gorm"
	"mygolib/modules/myLogger"
)

type DbBase struct {
	//times
	RecUpdTs time.Time
	RecCrtTs time.Time
}

func (t *DbBase) BeforeCreate(scope *gorm.Scope) (err error) {
	myLogger.Debug( "get in BeforeCreate")
	scope.SetColumn("RecUpdTs", time.Now().UTC())
	scope.SetColumn("RecCrtTs", time.Now().UTC())
	return nil
}

func (t *DbBase) BeforeUpdate(scope *gorm.Scope) (err error) {
	myLogger.Debug("get in BeforeUpdate")
	scope.SetColumn("RecUpdTs", time.Now().UTC())
	return nil
}

