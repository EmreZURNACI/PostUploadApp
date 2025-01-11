package Token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const SECRET_KEY string = "sa türk war mi"

func CreateToken(email string, telno string) (string, error) {
	// Gizli anahtar (Token imzalama için kullanılan key)
	// Production ortamında bu, güvenli bir şekilde saklanmalıdır.

	// Kullanıcıya özel veriler (payload) içeren token claims oluşturuyoruz.
	claims := jwt.MapClaims{
		"sub": email,                            // Subject (kimlik bilgisi olarak e-posta)
		"tel": telno,                            // Telefon numarası (örnek payload ek verisi)
		"iss": "PostUploadApp",                  // Issuer (bu token'ı oluşturan uygulama adı)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time (token'ın geçerlilik süresi)
		"iat": time.Now().Unix(),                // Issued at (token'ın oluşturulma zamanı)
	}

	// Token oluşturmak için claims ve imzalama yöntemi belirliyoruz.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Token'ı gizli anahtarla imzalıyoruz.
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err // Eğer hata oluşursa hatayı döndür
	}

	// Oluşturulan token'ı döndür
	return tokenString, nil
}
func ParseToken(tokenString string) (string, error) {
	// Token'ı çözümlemek ve geçerliliğini kontrol etmek
	//token, err istenirse token üzerinden claimler alınabilir
	// JWT'nin imzalanan algoritması doğrulandıktan sonra gizli anahtar döndürülüyor
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// JWT'nin imzalanan algoritması doğrulandıktan sonra gizli anahtar döndürülüyor
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token method invalid")
		}
		return []byte(SECRET_KEY), nil
	})

	// Eğer hata oluşursa, geçersiz token olduğunu belirtiyoruz
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	// Eğer token geçerliyse, payload kısmını çıkarabiliriz
	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

	// 	email := claims["sub"].(string)
	// 	telno := claims["tel"].(string)
	// }

	// Token geçerli olduğunda, "Token geçerli" mesajını döndürüyoruz
	return "Token geçerli", nil
}
