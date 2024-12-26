package Database

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func SignUp(name string, lastname string, nickname string, email string, password string, tel string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	file, err := ReadInfo(fmt.Sprintf("%s/Desktop/PostUploadApp/Database/config.json", homeDir))
	if err != nil {
		return err
	}
	con, err := Connection(file)
	if err != nil {
		return err
	}
	defer con.Close()
	var result string
	err = con.QueryRow("SELECT CheckUserInfo($1,$2,$3);", nickname, email, tel).Scan(&result)
	if err != nil {
		return errors.New("Sorgu çalıştırılırken hata ile karşılaşıldı," + err.Error())
	}
	var r RecvMessage
	err = json.Unmarshal([]byte(result), &r)
	if err != nil {
		return errors.New("Message unmarshal edilemedi." + err.Error())
	}
	if !r.Status {
		return errors.New(r.Message)
	}
	uuid, err := uuid.NewV7()
	if err != nil {
		return errors.New("uuid oluşturulurken hata ile karşılaşıldı.")
	}
	sEnc := b64.StdEncoding.EncodeToString([]byte(password))
	res, err := con.Exec("INSERT INTO public.user (uuid,name,lastname,nickname,email,password,tel) VALUES ($1,$2,$3,$4,$5,$6,$7);", uuid, name, lastname, nickname, email, sEnc, tel)
	if err != nil {
		return err
	}
	number, err := res.RowsAffected()
	if err != nil {
		return errors.New("Etkilenen satır sayısı alınamadı.")
	}

	if number == 0 {
		return errors.New("Hiçbir satır eklenmedi.")
	}

	return nil
}
func SignIn(telno string, email string, password string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	file, err := ReadInfo(fmt.Sprintf("%s/Desktop/PostUploadApp/Database/config.json", homeDir))
	if err != nil {
		return err
	}
	con, err := Connection(file)
	if err != nil {
		return err
	}
	defer con.Close()
	var result string
	err = con.QueryRow("SELECT IsUserExist($1,$2);", email, telno).Scan(&result)
	if err != nil {
		return errors.New("Sorgu çalıştırılamadı: " + err.Error())
	}

	var u RecvMessage
	err = json.Unmarshal([]byte(result), &u)
	if err != nil {
		return errors.New("Unmarshal hatası " + err.Error())
	}
	if !u.Status {
		return errors.New(u.Message)
	}
	var number int
	sEnc := b64.StdEncoding.EncodeToString([]byte(password))
	err = con.QueryRow("SELECT COUNT(*) FROM public.user WHERE (email=$1 OR tel=$2) AND password=$3;", email, telno, sEnc).Scan(&number)
	if err != nil {
		return errors.New("Sorgu çalıştırılamadı: " + err.Error())
	}
	if number != 1 {
		return errors.New("Kullanıcı bilgileriniz yanlış.")
	}
	return nil
}
