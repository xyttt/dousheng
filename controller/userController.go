package controller

import (
	"dousheng/dao"
	"dousheng/data"
	"dousheng/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	data.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	data.Response
	User service.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	u, err := dao.GetTableUserByUsername(username)

	if err != nil {
		fmt.Println("find name error")
	}

	if u.Name == username {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: data.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser, _ := dao.WriteTable(username, service.EnCoder(password))
		if dao.InsertTableUser(&newUser) != true {
			fmt.Println("insert data fail")
		}
		u, err := dao.GetTableUserByUsername(username)
		token := service.GenerateToken(username)
		if err != nil {
			fmt.Println("generate tocken error")
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: data.Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	encoderPassword := service.EnCoder(password)
	println(encoderPassword)

	u, err := dao.GetTableUserByUsername(username)

	if err != nil {
		fmt.Println("find name error")
	}

	if u.Password == encoderPassword {
		token := service.GenerateToken(username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: data.Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: data.Response{StatusCode: 1, StatusMsg: "username or password error"},
		})
	}
}

func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	id, _ := strconv.ParseInt(user_id, 10, 64)
	usi := service.UserServiceImpl{}
	if u, err := usi.GetUserById(id); err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: data.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: data.Response{StatusCode: 0},
			User:     u,
		})
	}
}
