package controllers

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	// "go.mongodb.org/mongo-driver/bson/primitive"

	"linktory/database"
)


type Data struct {
	ID string	`json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"-"`
	CreatedAt    time.Time          `json:"created_at"`
    UpdatedAt    time.Time          `json:"updated_at"`
}

var Secretkey = os.Getenv("SECRET_KEY")
var DB_NAME = os.Getenv("DB_NAME")
var collection = database.GetCollection("Users")

func Register(c *fiber.Ctx) error {
	 data := new(Data)

	if err := c.BodyParser(data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	data.Password = string(password)
	data.ID = uuid.New().String()
	data.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    data.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	filter := bson.M{"email": data.Email}
	count, _ := collection.CountDocuments(context.TODO(), filter)

	if(count > 0) {
		return c.JSON(fiber.Map{
			"message":"user with this email exists",
		})
	}
		
	
	collection.InsertOne(context.TODO(), data)
	

	return c.JSON(fiber.Map{
		"message": "User account successfully created",
		"data": data,
	})
}

func Login(c *fiber.Ctx) error {
	 data := new(Data)
	 user := new(Data)
	 
	if err := c.BodyParser(data); err != nil {
		return err
	}
	filter := bson.M{"email": data.Email}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if(err != nil){
			return c.JSON(fiber.Map{
			"message":	"user with this email does not exist",
		})
	}
	userPassword  := []byte(user.Password)
	dataPassword := []byte(data.Password)


	if err := bcrypt.CompareHashAndPassword(userPassword, dataPassword); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}
	id := user.ID
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    id,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(Secretkey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		SameSite: "strict",
		// Secure: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
		"data": user,
	})

}


func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}



func ChangePassword(c *fiber.Ctx) error {
	data := new(Data)
	
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secretkey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "user is unauthenticated",
		})
	}

	claims, _ := token.Claims.(*jwt.StandardClaims)
	userID := claims.Issuer

	if err := c.BodyParser(data); err != nil {
		return err
	}
	password := data.Password

	if userID == "" {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	newpassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	// Here we might choose to send a link for the user to update the password,
	// since I'm not working with a mailing service yet we do this (below) instead.
		
	filter := bson.M{"id": userID}
	update := bson.M{"password": newpassword}
		
	_, replaceerr := coll.ReplaceOne(context.TODO(), filter, update)
	
	if replaceerr != nil{
		c.Status(fiber.StatusNotFound)
		return	c.JSON(fiber.Map{
				"message":	"Couldn't update",
			})
	}

	c.Status(fiber.StatusFound)
		return	c.JSON(
			fiber.Map{
				"message":"Password updated successfully",
			})
	
}

func DeleteAccount(c *fiber.Ctx) error {
	 data := new(Data)
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secretkey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "user is unauthenticated",
		})
	}

	claims, _ := token.Claims.(*jwt.StandardClaims)
	user_ID := claims.Issuer

	if err := c.BodyParser(data); err != nil {
		return err
	}

	userID, _ := strconv.Atoi(user_ID)

	if userID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	update_cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&update_cookie)
	return c.JSON(fiber.Map{
		"message": "User Account deleted successfully",
	})
}

func ResetPassword( c *fiber.Ctx) error{
	return c.JSON("his")
}
