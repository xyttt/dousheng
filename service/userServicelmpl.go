package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"dousheng/config"
	"dousheng/dao"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserServiceImpl struct {
}

// GetUserById 未登录情况下,根据user_id获得User对象
func (usi *UserServiceImpl) GetUserById(id int64) (User, error) {
	user := User{
		Id:             0,
		Name:           "",
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		TotalFavorited: 0,
		FavoriteCount:  0,
	}
	tableUser, err := dao.GetTableUserById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user, err
	}
	log.Println("Query User Success")
	// followCount, _ := usi.GetFollowingCnt(id)
	// if err != nil {
	// 	log.Println("Err:", err.Error())
	// }
	// followerCount, _ := usi.GetFollowerCnt(id)
	// if err != nil {
	// 	log.Println("Err:", err.Error())
	// }
	// u := GetLikeService() //解决循环依赖
	// totalFavorited, _ := u.TotalFavourite(id)
	// favoritedCount, _ := u.FavouriteVideoCount(id)
	user = User{
		Id:   id,
		Name: tableUser.Name,
		// FollowCount:    followCount,
		// FollowerCount:  followerCount,
		IsFollow: false,
		// TotalFavorited: totalFavorited,
		// FavoriteCount:  favoritedCount,
	}
	return user, nil
}

// GenerateToken 根据username生成一个token
func GenerateToken(username string) string {
	u, _ := dao.GetTableUserByUsername(username)
	fmt.Printf("generatetoken: %v\n", u)
	expiresTime := time.Now().Unix() + int64(config.OneDayOfHours)
	fmt.Printf("expiresTime: %v\n", expiresTime)
	id64 := u.Id
	fmt.Printf("id: %v\n", strconv.FormatInt(id64, 10))
	claims := jwt.StandardClaims{
		Audience:  u.Name,
		ExpiresAt: expiresTime,
		Id:        strconv.FormatInt(id64, 10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "tiktok",
		NotBefore: time.Now().Unix(),
		Subject:   "token",
	}
	var jwtSecret = []byte(config.Secret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
		token := "Bearer " + token
		println("generate token success!\n")
		println(token)
		return token
	} else {
		println("generate token fail\n")
		return "fail"
	}
}

// type MyClaims struct {
// 	Username             string `json:"username"`
// 	Id                   string `json:"userId"`
// 	jwt.RegisteredClaims        // 注意!这是jwt-go的v4版本新增的，原先是jwt.StandardClaims
// }

// var MySecret = []byte("dousheng") // 定义secret

// // GenerateToken 根据username生成一个token
// func GenerateToken(username string) string {
// 	u, _ := dao.GetTableUserByUsername(username)
// 	ID := u.Id
// 	fmt.Printf("generatetoken: %v\n", u)

// 	claim := MyClaims{
// 		Username: username,
// 		Id:       strconv.FormatInt(ID, 10),
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
// 			NotBefore: jwt.NewNumericDate(time.Now()),                                       // 生效时间
// 		}}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
// 	tokenString, err := token.SignedString(MySecret)
// 	if err != nil {
// 		println("tocken generate error")
// 	}
// 	return tokenString
// }

// func Secret() jwt.Keyfunc {
// 	return func(token *jwt.Token) (interface{}, error) {
// 		return []byte("dousheng"), nil // secret
// 	}
// }

// func ParseToken(tokenss string) (*MyClaims, error) {
// 	token, err := jwt.ParseWithClaims(tokenss, &MyClaims{}, Secret())
// 	if err != nil {
// 		if ve, ok := err.(*jwt.ValidationError); ok {
// 			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
// 				return nil, errors.New("that's not even a token")
// 			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
// 				return nil, errors.New("token is expired")
// 			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
// 				return nil, errors.New("token not active yet")
// 			} else {
// 				return nil, errors.New("couldn't handle this token")
// 			}
// 		}
// 	}
// 	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
// 		return claims, nil
// 	}
// 	return nil, errors.New("couldn't handle this token")
// }

// EnCoder 密码加密
func EnCoder(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Result: " + sha)
	return sha
}
