package Database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	_ "github.com/lib/pq"
)

func ReadInfo(filePath string) (*ConfigFileInfo, error) {
	var file string
	if runtime.GOOS == "windows" {
		file = strings.ReplaceAll(filePath, "/", "\\")
	}
	bs, err := os.ReadFile(fmt.Sprintf("%s", file))
	if err != nil {
		return nil, errors.New("Dosya Okunamadı.")
	}
	var c ConfigFileInfo
	err = json.Unmarshal(bs, &c)
	if err != nil {
		return nil, errors.New("Unmarshal edilemedi.")
	}
	return &c, nil
}
func Connection(c *ConfigFileInfo) (*sql.DB, error) {
	var dsn string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Dbname)
	con, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.New("Bağlantı açılamadı.")
	}
	err = con.Ping()
	if err != nil {
		return nil, errors.New("Veri tabanı bağlantısı kurulamadı.")
	}
	return con, nil
}
