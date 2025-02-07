// *****************************************************************************
// 作者: lgdz
// 创建时间: 2025/2/7
// 描述：kingbase数据库配置
// *****************************************************************************

package pgsql

import "github.com/lgdzz/vingo-utils-v2/config"

type Config struct {
	config.Config
	Host         string `yaml:"host" json:"host"`
	Port         string `yaml:"port" json:"port"`
	Dbname       string `yaml:"dbname" json:"dbname"`
	Username     string `yaml:"username" json:"username"`
	Password     string `yaml:"password" json:"password"`
	Charset      string `yaml:"charset" json:"charset"`
	MaxIdleConns int    `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns" json:"maxOpenConns"`
}
