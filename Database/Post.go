package Database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func CreatePost(tx *sql.Tx, user_id string, imagePath string) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return errors.New("Uuid generate edilemedi.")
	}
	var post_number int
	err = tx.QueryRow("SELECT post_ordered FROM public.post WHERE user_id =$1 ORDER BY post_ordered DESC LIMIT 1", user_id).Scan(&post_number)
	if err == sql.ErrNoRows {
		post_number = 0
	} else if err != nil {
		return errors.New("Postların sırası elde edilemedi.")
	}
	_, err = tx.Exec("INSERT INTO public.post(uuid, user_id, post_ordered, comment_count, like_count, dislike_count, image_path) VALUES ($1,$2,$3,$4,$5,$6,$7);", uuid, user_id, post_number+1, 0, 0, 0, imagePath)
	if err != nil {
		return errors.New("Veri tabanı sorgusu çalıştırılırken hata ile karşılaşıldı.")
	}
	return nil
}
func LikePost(uuid string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	file, err := ReadInfo(fmt.Sprintf("%s/Desktop/PostUploadApp/Database/config.json", homeDir))
	if err != nil {
		return errors.New("Config dosyası okunamadı.")
	}
	con, err := Connection(file)
	if err != nil {
		return errors.New("Veri tabanı ile bağlantı sağlanamadı.")
	}
	defer con.Close()
	res, err := con.Exec("UPDATE public.post SET like_count = like_count + 1 WHERE uuid = $1", uuid)
	if err != nil {
		return errors.New("Sorgu çalıştırılamadı.")
	}

	number, err := res.RowsAffected()
	if err != nil {
		return errors.New("Etkilenen satır sayısı alınamadı.")
	}

	if number == 0 {
		return errors.New("Hiçbir satır güncellenmedi.")
	}

	return nil
}
func DislikePost(uuid string) error {
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
	res, err := con.Exec("UPDATE public.post SET dislike_count=dislike_count+1 WHERE uuid=$1;", uuid)
	if err != nil {
		return err
	}
	number, err := res.RowsAffected()
	if err != nil {
		return errors.New("Etkilenen satır sayısı alınamadı.")
	}

	if number == 0 {
		return errors.New("Hiçbir satır güncellenmedi.")
	}

	return nil
}
func CommentPost(post_id string, user_id string, text string) error {
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
	uuid, err := uuid.NewV7()
	if err != nil {
		return errors.New("Yeni uuid generate edilemedi.")
	}

	res, err := con.Exec("INSERT INTO public.comment (uuid,post_id,user_id,text,datetime) VALUES ($1,$2,$3,$4,now()::timestamp);", uuid, post_id, user_id, text)
	if err != nil {
		return errors.New("Yorum veritabanına eklenirken sorgu hatası ile karşılaşıldı." + err.Error())
	}
	number, err := res.RowsAffected()
	if err != nil {
		return errors.New("Etkilenen satır sayısı alınamadı.")
	}

	if number == 0 {
		return errors.New("Hiçbir satır güncellenmedi.")
	}

	return nil
}
