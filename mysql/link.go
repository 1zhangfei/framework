package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/1zhangfei/framework/config"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlInit(address string, hand func(cli *gorm.DB) error) error {
	err := config.ViperInit(address)
	if err != nil {
		return err
	}
	dataId := viper.GetString("Database.DataId")
	Group := viper.GetString("Database.Group")
	var Mysql struct {
		M struct {
			Username string
			Password string
			Host     string
			Port     string
			Database string
		} `json:"Mysql"`
	}
	getConfig, err := config.GetConfig(dataId, Group)
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(getConfig), &Mysql); err != nil {
		return err
	}

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		Mysql.M.Username,
		Mysql.M.Password,
		Mysql.M.Host,
		Mysql.M.Port,
		Mysql.M.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	d, _ := db.DB()
	defer d.Close()
	if err = hand(db); err != nil {
		return err
	}
	return nil
}
