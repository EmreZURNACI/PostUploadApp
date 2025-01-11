package main

import (
	"PostUploadApp/Api"
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/grpc"
)

func main() {
	var address string = ":3838"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Bağlantı hatası: %v\n", err)
		return
	}
	defer conn.Close()

	client := Api.NewPostAppClient(conn)

	var ctx context.Context = context.Background()
	res, err := SignIn(client, ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	res2, err2 := LikePost(client, ctx, res.Token, "ae65db75-e2e6-459f-8bb7-77b8fed9cbfc")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println(res2)

}

func CreatePost(client Api.PostAppClient, ctx context.Context, token string, user_id string) {

	stream, err := client.UploadPost(ctx)
	if err != nil {
		fmt.Printf("Stream başlatma hatası: %v\n", err)
		return
	}
	directory, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Kullanıcı home directory alınamadı: %v\n", err)
		return
	}
	var filepath string = fmt.Sprintf("%s/Desktop/TürkBayrağı", directory)
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Dosya açma hatası: %v\n", err)
		return
	}
	defer file.Close()

	//İlk streamden token ve userId'yi alsın diye ayrı gönderdim.
	//Çünkü hepsini beraber gönderdiğimde dosyaya yazma kımını atladığım için dosyanın bir miktar chunkı boşa gitti
	err = stream.Send(&Api.UploadPostReq{
		Token:  token,
		UserId: user_id,
	})
	if err != nil {
		fmt.Printf("Stream gönderme hatası: %v\n", err)
		return
	}

	buffer := make([]byte, 1024)
	for {
		//Dosyanın sonu
		n, err := file.Read(buffer)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			fmt.Printf("Dosya okuma hatası: %v\n", err)
			return
		}

		err = stream.Send(&Api.UploadPostReq{
			Chunk: buffer[:n],
		})
		if err != nil {
			fmt.Printf("Stream gönderme hatası: %v\n", err)
			return
		}
	}

	err = stream.CloseSend()
	if err != nil {
		fmt.Printf("Stream kapatma hatası: %v\n", err)
		return
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("Stream cevap alma hatası: %v\n", err)
		return
	}

	fmt.Printf("Sunucudan gelen yanıt: %+v\n", res)
}
func LikePost(client Api.PostAppClient, ctx context.Context, token string, uuid string) (*Api.LikePostRes, error) {
	res, err := client.LikePost(ctx, &Api.LikePostReq{
		Token: token,
		Uuid:  uuid,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
func DislikePost(client Api.PostAppClient, ctx context.Context, token string, uuid string) (*Api.DislikePostRes, error) {
	res, err := client.DislikePost(ctx, &Api.DislikePostReq{
		Token: token,
		Uuid:  uuid,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
func CommentPost(client Api.PostAppClient, ctx context.Context, token string, comment string, post_uuid string, user_uuid string) (*Api.CommentPostRes, error) {
	res, err := client.CommentPost(ctx, &Api.CommentPostReq{
		Token:    token,
		UserUuid: user_uuid,
		PostUuid: post_uuid,
		Comment:  comment,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
func SignUp(client Api.PostAppClient, ctx context.Context) (*Api.SignUpRes, error) {
	res, err := client.SignUp(ctx, &Api.SignUpReq{
		Name:     "Emre",
		Lastname: "ZURNACI",
		Nickname: "EmreZ2",
		Email:    "Emrez2@gmail.com",
		Password: "123",
		Tel:      "+905061546782",
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
func SignIn(client Api.PostAppClient, ctx context.Context) (*Api.SignInRes, error) {
	res, err := client.SignIn(ctx, &Api.SignInReq{
		Email:    "Emrez2@gmail.com",
		Telno:    "+905061546782",
		Password: "123",
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
