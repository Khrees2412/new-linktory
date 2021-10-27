package controllers

import (
	"context"
	"fmt"

	"linktory/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var coll = database.GetCollection("Links")
// var Secretkey = "secretjunglekey"

type Link struct {
	// ID	primitive.ObjectID `bson:"_id"`
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name" validate:"required, min=2,max=100"`
	Url    string `json:"url"`
}

func CreateLink(c *fiber.Ctx) error {
	link := new(Link)

	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secretkey), nil
	})
	claims, _ := token.Claims.(*jwt.StandardClaims)
	userID := claims.Issuer
	if err != nil {
		c.JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}
	if err := c.BodyParser(link); err != nil {
		return err
	}
	
	link.UserID = userID
	link.ID  = uuid.New().String()
	
	coll.InsertOne(context.TODO(), link)

	c.Status(fiber.StatusCreated)
	return c.JSON(
		fiber.Map{
			"message" : "link created",
			"data" : link,
		})

}


func GetLink(c *fiber.Ctx) error {
 	link := new(Link)
 	linkCopy := new(Link)


	if err := c.BodyParser(link); err != nil {
		return err
	}
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secretkey), nil
	})
	claims, _ := token.Claims.(*jwt.StandardClaims)
	userID := claims.Issuer
	if err != nil {
		c.JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	filter := bson.M{"id": link.ID, "user_id": userID}
	geterr := coll.FindOne(context.TODO(), filter).Decode(&linkCopy)
	
	if geterr != nil{
		c.Status(fiber.StatusNotFound)
	return c.JSON(
		fiber.Map{
			"message" : "link doesn't exist ",
			"data" : "",
		})
	}
	c.Status(fiber.StatusFound)
		return c.JSON(
		fiber.Map{
			"message" : "link found",
			"data" : linkCopy,
		})
}


func GetLinks(c *fiber.Ctx) error {
	link := new(Link)
	var links []*Link

	if err := c.BodyParser(link); err != nil {
		return err
	}
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secretkey), nil
	})
	claims, _ := token.Claims.(*jwt.StandardClaims)
	userID := claims.Issuer
	if err != nil {
		c.JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}
	filter := bson.M{"id": link.ID, "user_id": userID}
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil{
			return c.JSON(
		fiber.Map{
			"message" : "links were not found",
		})
	}
	   for cur.Next(context.TODO()) {
		   linkCopy := new(Link)
        err = cur.Decode(&linkCopy)
        if err != nil {
				return c.JSON(fiber.Map{
					"message" : "error ocurred getting links",
				})
		}
        links = append(links, linkCopy)
    }
		return c.JSON(
		fiber.Map{
			"message" : fmt.Sprintf("Found %d links",len(links)),
			"data" : links,
		})
}


func DeleteLink(c *fiber.Ctx) error {
	link := new(Link)

	if err := c.BodyParser(link); err != nil {
		return err
	}
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secretkey), nil
	})
	claims, _ := token.Claims.(*jwt.StandardClaims)
	userID := claims.Issuer
	if err != nil {
		c.JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}
	filter := bson.M{"id": link.ID, "user_id": userID}
	delerr := coll.FindOneAndDelete(context.TODO(), filter)
	
	if delerr != nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"messaage" : "link with given id does not exist",
			"link" : link,
		})
	}
	c.Status(fiber.StatusFound)
	return c.JSON(fiber.Map{
			"messaage" : "link deleted",
			
		})
}

func UpdateLink(c *fiber.Ctx) error {
	link := new(Link)

	if err := c.BodyParser(link); err != nil {
		return err
	}
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secretkey), nil
	})
	claims, _ := token.Claims.(*jwt.StandardClaims)
	userID := claims.Issuer
	if err != nil {
		c.JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}
	filter := bson.M{"id": link.ID, "user_id": userID}
	update := bson.M{
		"name": link.Name,
		"url": link.Url,
		"user_id" : userID,
		 }
		
	_, reperr := coll.ReplaceOne(context.TODO(), filter, update)
	
	if reperr != nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message" : "link with given id does not exist",
			"data":link,
		})
	
	}
	c.Status(fiber.StatusFound)
	return c.JSON(fiber.Map{
			"message" : "link updated",
		})
		
}