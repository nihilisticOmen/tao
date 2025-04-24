package jwts

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type JwtToken struct {
	AccessToken  string
	RefreshToken string
	AccessExp    int64
	RefreshExp   int64
}

func CreateToken(val string, exp time.Duration, secret string, rf time.Duration, rfSecret string) *JwtToken {
	aExp := time.Now().Add(exp).Unix()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tokenKey": val,
		"exp":      aExp,
	})
	atoken, err := claims.SignedString([]byte(secret))
	if err != nil {
		log.Println(err)
	}
	rExp := time.Now().Add(rf).Unix()
	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"key": val,
		"exp": rExp,
	})
	rToken, err := refreshClaims.SignedString([]byte(rfSecret))
	if err != nil {
		log.Println(err)
	}
	return &JwtToken{
		AccessToken:  atoken,
		RefreshToken: rToken,
		AccessExp:    aExp,
		RefreshExp:   rExp,
	}
}

func ParseToken(tokenString string, secret string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("%v \n", claims)
	} else {
		fmt.Println(err)
	}
}
