package main

import (
	"context"
	"exam_fpt/database"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func init() {
	database.MongoInit()

}

func addCustomer(c echo.Context) error {
	customer := echo.Map{}
	if err := c.Bind(&customer); err != nil {
		return err
	}
	dbo := database.GetMgoDB()
	res, err := dbo.Collection("customers").InsertOne(context.TODO(), customer)
	if err != nil {
		return err
	}
	return c.JSON(200, res)
}

func getCustomers(c echo.Context) error {
	filterParams := c.QueryParam("filter")
	page, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if err != nil || page == 0 {
		page = 1
	}
	pageSize, err := strconv.ParseInt(c.QueryParam("pageSize"), 10, 64)

	if pageSize == 0 || err != nil {
		pageSize = 30
	}
	query := bson.M{}
	var customers []interface{}
	if filterParams != "" {
		filters := strings.Split(filterParams, ";")
		for _, filter := range filters {
			pairs := strings.Split(filter, ":")
			if pairs[0] != "" && pairs[1] != "" {
				query[pairs[0]] = pairs[1]
			}
		}
	}
	dbo := database.GetMgoDB()
	total, err := dbo.Collection("customers").CountDocuments(context.TODO(), query)
	if total == 0 {
		return c.JSON(200, map[string]interface{}{"total": 0, "docs": ""})
	}
	opts := options.Find()
	opts.SetLimit(pageSize)
	opts.SetSkip(pageSize * (page - 1))
	cur, err := dbo.Collection("customers").Find(context.TODO(), query, opts)
	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var customer interface{}
		if err = cur.Decode(&customer); err != nil {
			return err
		}
		customers = append(customers, &customer)
	}

	return c.JSON(200, map[string]interface{}{"total": total, "docs": customers, "page": page, "pageSize": pageSize})

}

func deleteCustomer(c echo.Context) error {
	id := c.Param("customerId")
	objectIDS, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	doc := bson.D{{"_id", objectIDS}}
	dbo := database.GetMgoDB()
	res, err := dbo.Collection("customers").DeleteOne(context.TODO(), doc)

	if err != nil {
		return err
	}
	return c.JSON(200, res)

}

func getCustomer(c echo.Context) error {
	id := c.Param("customerId")
	objectIDS, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	idDoc := bson.D{{"_id", objectIDS}}
	dbo := database.GetMgoDB().Collection("customers")
	var customer interface{}
	if err = dbo.FindOne(context.TODO(), idDoc).Decode(&customer); err != nil {
		return err
	}
	return c.JSON(200, customer)
}

func updateCustomer(c echo.Context) error {
	id := c.Param("customerId")
	objectIDS, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	idDoc := bson.D{{"_id", objectIDS}}
	var updates = echo.Map{}
	if err = c.Bind(&updates); err != nil {
		return err
	}
	dbo := database.GetMgoDB().Collection("customers")
	var customer interface{}
	err = dbo.FindOne(context.TODO(), idDoc).Decode(&customer)
	if err != nil {
		return err
	}

	_, err = dbo.UpdateOne(context.TODO(), idDoc, bson.D{{"$set", updates}, {"$currentDate", bson.D{{"modifiedAt", true}}}})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.POST("/customer", addCustomer)
	e.DELETE("/customer/:customerId", deleteCustomer)
	e.PATCH("/customer/:customerId", updateCustomer)
	e.GET("/customer/:customerId", getCustomer)
	e.GET("/customers", getCustomers)
	e.Logger.Fatal(e.Start(":8080"))
}
