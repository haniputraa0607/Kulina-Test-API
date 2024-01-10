package utils

import (
	"fmt"
	"rest-api/model"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = "SECRET_KEY"

type UserMiddleware struct {
	ID uint64
	Email string
	Name string
}

func GenerateTokenUser(user model.User) (string, error) {

	payloads := jwt.MapClaims{
		"id"   : user.ID,
		"email": user.Email,
		"name" : user.Name,
		"exp"  : time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payloads)

	webToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return webToken, nil
}

func GenerateTokenSupplier(user model.Supplier) (string, error) {

	payloads := jwt.MapClaims{
		"id"   : user.ID,
		"email": user.Email,
		"name" : user.Name,
		"exp"  : time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payloads)

	webToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return webToken, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {

	tokenJwt, errJwt := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, isValid := t.Method.(*jwt.SigningMethodHMAC)
		if !isValid {
			return nil, fmt.Errorf("unexpected singing method: %v", t.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if errJwt != nil {
		return nil, errJwt
	}

	return tokenJwt, nil
}

func DecodeToken(tokenString string) (jwt.MapClaims, error) {

	VerifyToken, err := VerifyToken(tokenString)

	if err != nil {
		return nil, err
	}

	claims, isOk := VerifyToken.Claims.(jwt.MapClaims)

	if isOk && VerifyToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func UserMid(decodeToken jwt.MapClaims)  (*UserMiddleware, bool){
	
	user := new(UserMiddleware)

	UserID, ok := decodeToken["id"]
	if !ok {
		return user, true
	}

	UserEmail, ok := decodeToken["email"]
	if !ok {
		return user, true
	}

	UserName, ok := decodeToken["name"]
	if !ok {
		return user, true
	}

	
	user.ID = uint64(UserID.(float64))
	user.Email = UserEmail.(string)
	user.Name = UserName.(string)
	
	return user, false


}