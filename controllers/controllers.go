package controllers

import (
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	"time"
	"net/http"
	"scoutiq_server/db"
	_ "strconv"
	"github.com/dgrijalva/jwt-go"
	_ "fmt"
)

type AuthRequest struct {
	Key string `json:"key"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	ExpiresIn string  `json:"expiresIn"`
}

func Register(route *echo.Group) {
	//route.GET("", readMainHandler)
	//route.POST("/event", createEventHandler)
	//route.GET("/:id", readByIDHandler)
	//route.PUT("/:id", updateByIDHandler)
	//route.DELETE("/:id", deleteByIDHandler)
	route.POST("/token", getTokenHandler)
}

func DatabaseAccess(c echo.Context) *gorm.DB {
	database := c.Get("database").(*gorm.DB)
	return database
}
//
//func readMainHandler(c echo.Context) (err error) {
//	database := DatabaseAccess(c)
//	user := db.User{}
//	database.First(&user)
//	return c.JSON(http.StatusOK, user)
//}
//
//func createEventHandler(c echo.Context) error {
//	database := DatabaseAccess(c)
//	event := &db.Event{
//		CreatedAt: time.Now().Unix(),
//	}
//	if err := c.Bind(event); err != nil {
//		return err
//	}
//	database.Create(&event)
//	return c.JSON(http.StatusCreated, event) // <-- 201
//}
//
//func readByIDHandler(c echo.Context) error {
//	database := DatabaseAccess(c)
//	id, _ := strconv.Atoi(c.Param("id"))
//	event := &db.Event{
//		ID: id,
//	}
//	database.Find(&event)
//	return c.JSON(http.StatusOK, event) // <-- 200
//}
//
//func updateByIDHandler(c echo.Context) error {
//	database := DatabaseAccess(c)
//	id, _ := strconv.Atoi(c.Param("id"))
//	event := &db.Event{
//		ID: id,
//		UpdatedAt: time.Now().Unix(),
//	}
//	if err := c.Bind(event); err != nil {
//		return err
//	}
//	database.Model(&db.Event{}).UpdateColumns(&event)
//	return c.JSON(http.StatusOK, event) // <-- 200
//}
//
//func deleteByIDHandler(c echo.Context) error {
//	database := DatabaseAccess(c)
//	id, _ := strconv.Atoi(c.Param("id"))
//	event := &db.Event{
//		ID: id,
//	}
//	database.Delete(&event)
//	return c.JSON(http.StatusOK, event) // <-- 200
//}

func getTokenHandler(c echo.Context) (err error) {
	u := new(db.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	database := DatabaseAccess(c)
	user := new(db.User)

	if err := database.Where("email = ?", u.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err.Error()) // <-- 400 Not found username
	}

	if u.Password == user.Password {
		exp := time.Now().Add(time.Second * 30).UTC()
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["exp"] = exp.Unix()

		cfgSecretKey := c.Get("SecretKey").(string)
		authRes := new(AuthResponse)
		authRes.ExpiresIn = exp.Format(time.RFC3339)
		authRes.TokenType = "Bearer"
		authRes.Token, err = token.SignedString([]byte(cfgSecretKey))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, authRes) // <-- 200
	}
	return c.JSON(http.StatusBadRequest, "bad password")

}

//func getTestTokenHandler(c echo.Context) (err error) {
//	user := c.Get("user").(*jwt.Token)
//	claims := user.Claims.(jwt.MapClaims)
//	id := claims["appKey"].(string)
//
//	head := c.Request().Header
//	claim := head.Get("Authorization")
//	fmt.Println(claim)
//
//	return c.String(http.StatusOK, fmt.Sprintf("Welcome %f !", id))
//}