package Database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ReadInfo(filePath string) (*ConfigFileInfo, error) {
	bs, err := os.ReadFile(fmt.Sprintf("%s", filePath))
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
	var dsn string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", c.Host, c.Port, c.User, c.Password, c.Dbname)
	con, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.New("Bağlantı açılamadı.")
	}
	err = con.Ping()
	if err != nil {
		return nil, errors.New("Bağlantı kurulamadı.")
	}
	return con, nil
}
