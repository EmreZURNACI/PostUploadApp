package main

import (
	Api "PostUploadApp/Api"
	db "PostUploadApp/Database"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

func main() {
	var address string = ":3838"
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Listener hatası.")
		return
	}
	defer ln.Close()
	server := grpc.NewServer()
	Api.RegisterPostAppServer(server, &PostAppServer{})
	fmt.Printf("Server %s adresinden ayağa kaldırıldı.\n", address)
	err = server.Serve(ln)
	if err != nil {
		fmt.Println("Server ayağa kaldırılamadı.")
	}
	defer server.Stop()

}

type PostAppServer struct {
	Api.UnimplementedPostAppServer
}

func (p *PostAppServer) UploadPost(stream Api.PostApp_UploadPostServer) error {

	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("Stream alınamadı: %v", err)
	}

	if req.GetToken() != "randomjwt" {
		return errors.New("Unauthorized: Token bulunamadı.")
	}

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Directorye erişilemedi: %v", err)
	}

	imagePath := fmt.Sprintf("%s/Desktop/PostUploadApp/Images/%d", homeDirectory, time.Now().UnixNano())

	configFile, err := db.ReadInfo("/home/emre/Desktop/PostUploadApp/Database/config.json")
	if err != nil {
		return fmt.Errorf("Config dosyası okunamadı: %v", err)
	}

	con, err := db.Connection(configFile)
	if err != nil {
		return fmt.Errorf("Veritabanı bağlantısı hatası: %v", err)
	}
	defer con.Close()

	// Veritabanı transaction'ını başlatıyoruz
	tx, err := con.Begin()
	if err != nil {
		return fmt.Errorf("Veritabanı transaction başlatılamadı: %v", err)
	}

	// Hata durumunda rollback işlemi yapacağız
	defer func() {
		if err != nil {
			// Transaction geri alınacak
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Printf("Rollback hatası: %v\n", rollbackErr)
			}
		}
	}()

	file, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("Dosya oluşturulamadı: %v", err)
	}

	err = db.CreatePost(tx, "d226e771-6eec-4861-a9b3-01dbce7a5a81", imagePath)
	if err != nil {
		file.Close()
		os.Remove(imagePath)
		return fmt.Errorf("Kullanıcı bilgileri veritabanına kaydedilemedi: %v", err)
	}

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			file.Close()
			os.Remove(imagePath)
			return fmt.Errorf("Stream'den veriler dosyaya yazılamadı: %v", err)
		}

		_, err = file.Write(req.GetChunk())
		if err != nil {
			file.Close()
			os.Remove(imagePath)
			return fmt.Errorf("Dosyaya yazma hatası: %v", err)
		}
	}

	err = file.Close()
	if err != nil {
		os.Remove(imagePath)
		return fmt.Errorf("Dosya kapatılamadı: %v", err)
	}

	// Transaction'ı commit ediyoruz
	err = tx.Commit()
	if err != nil {
		// Commit hatası olursa dosyayı silip transaction geri alınacak
		os.Remove(imagePath)
		return fmt.Errorf("Veritabanı işlemi commit edilemedi: %v", err)
	}

	return stream.SendAndClose(&Api.UploadPostRes{
		Status:     true,
		StatusCode: 200,
		Message:    fmt.Sprintf("Dosya başarıyla %s olarak kaydedildi", file.Name()),
	})
}
func (p *PostAppServer) SignIn(ctx context.Context, req *Api.SignInReq) (*Api.SignInRes, error) {
	err := db.SignIn(req.GetTelno(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &Api.SignInRes{
		Status:     true,
		StatusCode: 200,
		Message:    "Başarıyla giriş yapıldı.",
		Token:      "randomjwt",
	}, nil
}
func (p *PostAppServer) SignUp(ctx context.Context, req *Api.SignUpReq) (*Api.SignUpRes, error) {
	err := db.SignUp(req.GetName(), req.GetLastname(), req.GetNickname(), req.GetEmail(), req.GetPassword(), req.GetTel())
	if err != nil {
		return nil, err
	}
	return &Api.SignUpRes{
		Status:     true,
		StatusCode: 200,
		Message:    fmt.Sprintf("%s %s adlı kullanıcı başarıyla eklendi.", req.GetName(), req.GetLastname()),
	}, nil
}
func (p *PostAppServer) LikePost(ctx context.Context, req *Api.LikePostReq) (*Api.LikePostRes, error) {
	if req.GetToken() != "randomjwt" {
		return nil, errors.New("Unauthorized")
	}
	err := db.LikePost(req.GetUuid())
	if err != nil {
		return nil, err
	}
	return &Api.LikePostRes{
		Status:     true,
		StatusCode: 200,
		Message:    "Postu beğendiniz",
	}, nil
}
func (p *PostAppServer) DislikePost(ctx context.Context, req *Api.DislikePostReq) (*Api.DislikePostRes, error) {
	if req.GetToken() != "randomjwt" {
		return nil, errors.New("Unauthorized")
	}
	err := db.DislikePost(req.GetUuid())
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &Api.DislikePostRes{
		Status:     true,
		StatusCode: 200,
		Message:    "Postu dislikeladınız",
	}, nil
}
func (p *PostAppServer) CommentPost(ctx context.Context, req *Api.CommentPostReq) (*Api.CommentPostRes, error) {
	if req.GetToken() != "randomjwt" {
		return nil, errors.New("Unauthorized")
	}
	err := db.CommentPost(req.GetPostUuid(), req.GetUserUuid(), req.GetComment())
	if err != nil {
		return nil, err
	}
	return &Api.CommentPostRes{
		Status:     true,
		StatusCode: 200,
		Message:    fmt.Sprintf("%s uuid'li posta %s id'li kullanıcı tarafından yorum yapılmıştır.", req.GetPostUuid(), req.GetUserUuid()),
	}, nil
}
