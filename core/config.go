package core

import (
	"net"
	"os"
	"strconv"

	"github.com/alecthomas/kong"
	"github.com/go-sql-driver/mysql"
)

type Config struct {
	AppAddr string `env:"APP_ADDR" default:":8080"`
	DBUser  string `env:"DB_USER" default:"root"`
	DBPass  string `env:"DB_PASS" default:"pass"`
	DBHost  string `env:"DB_HOST" default:"localhost"`
	DBPort  int    `env:"DB_PORT" default:"3306"`
	DBName  string `env:"DB_NAME" default:"app"`
}

func (c *Config) Parse() {
	kong.Parse(c)
}

func (c Config) MySQLConfig() *mysql.Config {
	mc := mysql.NewConfig()

	// NeoShowcase環境かどうかをチェック
	_, isDeployedOnNeoShowcase := os.LookupEnv("NS_MARIADB_DATABASE")
	if isDeployedOnNeoShowcase {
		// NeoShowcase環境の場合はNS_MARIADB_*環境変数を使用
		mc.User = os.Getenv("NS_MARIADB_USER")
		mc.Passwd = os.Getenv("NS_MARIADB_PASSWORD")
		mc.Net = "tcp"
		mc.Addr = net.JoinHostPort(
			os.Getenv("NS_MARIADB_HOSTNAME"),
			os.Getenv("NS_MARIADB_PORT"),
		)
		mc.DBName = os.Getenv("NS_MARIADB_DATABASE")
	} else {
		// 通常環境の場合はDB_*環境変数を使用
		mc.User = c.DBUser
		mc.Passwd = c.DBPass
		mc.Net = "tcp"
		mc.Addr = net.JoinHostPort(c.DBHost, strconv.Itoa(c.DBPort))
		mc.DBName = c.DBName
	}
	mc.Collation = "utf8mb4_general_ci"
	mc.AllowNativePasswords = true

	return mc
}
