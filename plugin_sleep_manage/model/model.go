// Package model 睡眠管理数据库
package model

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/logoove/sqlite" // use sql
	log "github.com/sirupsen/logrus"
)

// SleepDB 睡眠数据库
type SleepDB gorm.DB

// Initialize 初始化
func Initialize(dbpath string) *SleepDB {
	var err error
	if _, err = os.Stat(dbpath); err != nil || os.IsNotExist(err) {
		// 生成文件
		f, err := os.Create(dbpath)
		if err != nil {
			return nil
		}
		defer f.Close()
	}
	gdb, err := gorm.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}
	gdb.AutoMigrate(&SleepManage{})
	return (*SleepDB)(gdb)
}

// Open 打开
func Open(dbpath string) (*SleepDB, error) {
	db, err := gorm.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	return (*SleepDB)(db), nil
}

// Close 关闭
func (sdb *SleepDB) Close() error {
	db := (*gorm.DB)(sdb)
	return db.Close()
}

// SleepManage 睡眠信息
type SleepManage struct {
	gorm.Model
	GroupID   int64     `gorm:"column:group_id"`
	UserID    int64     `gorm:"column:user_id"`
	SleepTime time.Time `gorm:"column:sleep_time"`
}

// TableName 表名
func (SleepManage) TableName() string {
	return "sleep_manage"
}

// Sleep 更新睡眠时间
func (sdb *SleepDB) Sleep(gid, uid int64) (position int, awakeTime time.Duration) {
	db := (*gorm.DB)(sdb)
	now := time.Now()
	var today time.Time
	if now.Hour() >= 21 {
		today = now.Add(-time.Hour*time.Duration(-21+now.Hour()) - time.Minute*time.Duration(now.Minute()) - time.Second*time.Duration(now.Second()))
	} else if now.Hour() <= 3 {
		today = now.Add(-time.Hour*time.Duration(3+now.Hour()) - time.Minute*time.Duration(now.Minute()) - time.Second*time.Duration(now.Second()))
	}
	st := SleepManage{
		GroupID:   gid,
		UserID:    uid,
		SleepTime: now,
	}
	if err := db.Debug().Model(&SleepManage{}).Where("group_id = ? and user_id = ?", gid, uid).First(&st).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			db.Debug().Model(&SleepManage{}).Create(&st) // newUser not user
		}
	} else {
		log.Println("sleeptime为", st)
		awakeTime = now.Sub(st.SleepTime)
		db.Debug().Model(&SleepManage{}).Where("group_id = ? and user_id = ?", gid, uid).Update(
			map[string]interface{}{
				"sleep_time": now,
			})
	}
	db.Debug().Model(&SleepManage{}).Where("group_id = ? and sleep_time <= ? and sleep_time >= ?", gid, now, today).Count(&position)
	return position, awakeTime
}

// GetUp 更新起床时间
func (sdb *SleepDB) GetUp(gid, uid int64) (position int, sleepTime time.Duration) {
	db := (*gorm.DB)(sdb)
	now := time.Now()
	today := now.Add(-time.Hour*time.Duration(-6+now.Hour()) - time.Minute*time.Duration(now.Minute()) - time.Second*time.Duration(now.Second()))
	st := SleepManage{
		GroupID:   gid,
		UserID:    uid,
		SleepTime: now,
	}
	if err := db.Debug().Model(&SleepManage{}).Where("group_id = ? and user_id = ?", gid, uid).First(&st).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			db.Debug().Model(&SleepManage{}).Create(&st) // newUser not user
		}
	} else {
		log.Println("sleeptime为", st)
		sleepTime = now.Sub(st.SleepTime)
		db.Debug().Model(&SleepManage{}).Where("group_id = ? and user_id = ?", gid, uid).Update(
			map[string]interface{}{
				"sleep_time": now,
			})
	}
	db.Debug().Model(&SleepManage{}).Where("group_id = ? and sleep_time <= ? and sleep_time >= ?", gid, now, today).Count(&position)
	return position, sleepTime
}
