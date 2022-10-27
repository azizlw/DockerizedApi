package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/azizlw/FinalProject/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://localhost:27017"
const dbName = "inventory"
const colName1 = "item"
const colName2 = "users"
const colName3 = "cart"

// (taking reference of mongodb collection)
var collection *mongo.Collection  // collection of items
var collection2 *mongo.Collection // collection of user credentials
var collection3 *mongo.Collection // collection of user cart

var jwtKey = []byte("secret_key")

// Connect with mongoDB

func init() {
	// client option
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")

	collection = client.Database(dbName).Collection(colName1)
	collection2 = client.Database(dbName).Collection(colName2)
	collection3 = client.Database(dbName).Collection(colName3)

	//collection instance
	fmt.Println("Collection instance is ready")
}

// MongoDB helpers - file

// registering the user in database
func register(credential model.Credentials) {
	// users[credentials.Username] = credentials.Password
	inserted, err := collection2.InsertOne(context.Background(), credential)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted 1 user in db with id: ", inserted.InsertedID)
}

// insert 1 record
func insertOneItem(item model.Inventory) {
	inserted, err := collection.InsertOne(context.Background(), item)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted 1 item in db with id: ", inserted.InsertedID)
}

// update 1 record
func updateOneItem(itemId string, quant int) {
	id, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"quantity": quant}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Modified count: ", result.ModifiedCount)
}

// delete 1 item
func deleteOneItem(itemId string) {
	id, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": id}

	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Item got delete with delete count: ", deleteCount)
}

// delete all items from mongodb
func deleteAllitems() int64 {
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Number of Items deleted: ", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

// get all items from database
func getAllItems() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var items []primitive.M

	for cur.Next(context.Background()) {
		var item bson.M
		err := cur.Decode(&item)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}
	defer cur.Close(context.Background())
	return items
}

// adding data to cart
func addToCart(itemId string, item model.Cart) {
	item.ItemId = itemId
	// need to check whether it is available in inventory and if available then add to cart.

	// not updating inventory now will update when user purchase the item from inventory.
	inserted, err := collection3.InsertOne(context.Background(), item)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted 1 item in cart with id: ", inserted.InsertedID)
}

// getting data from cart
func getcart(userName string) []model.Cart {

	cur, err := collection3.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var data []model.Cart
	for cur.Next(context.Background()) {
		var cartData model.Cart
		err := cur.Decode(&cartData)
		if err != nil {
			log.Fatal(err)
		}

		if userName == cartData.Username {
			data = append(data, cartData)
		}

	}
	return data
}

// Actual controllers file

// Authorizing user and generating token
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials model.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cur, err := collection2.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	flag := false

	for cur.Next(context.Background()) {
		var credential model.Credentials
		err := cur.Decode(&credential)
		if err != nil {
			log.Fatal(err)
		}

		if credentials.Username == credential.Username && credentials.Password == credential.Password {
			flag = true
			break
		}

	}
	defer cur.Close(context.Background())

	if !flag {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Hour * 2)
	claims := &model.Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })

	// w.Header().Set("Content-Type", "application/json")
	type stoken struct {
		Name    string
		Value   string
		Expires time.Time
	}

	var a stoken
	a.Name = "token"
	a.Value = tokenString
	a.Expires = expirationTime

	json.NewEncoder(w).Encode(a)
	fmt.Println()
}

func Register(w http.ResponseWriter, r *http.Request) {
	var credentials model.Credentials
	json.NewDecoder(r.Body).Decode(&credentials)
	register(credentials)
	json.NewEncoder(w).Encode("New User added successfully")
}

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allItems := getAllItems()
	json.NewEncoder(w).Encode(allItems)
}

func InsertOneItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var item model.Inventory
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Fatal(err)
	}

	insertOneItem(item)
	json.NewEncoder(w).Encode(item)
}

func UpdateOneItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-from-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	var item model.Inventory
	json.NewDecoder(r.Body).Decode(&item)
	updateOneItem(params["id"], item.ItemQuantity)

	json.NewEncoder(w).Encode(params["id"])
}

func DeleteOneItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-from-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	params := mux.Vars(r)
	deleteOneItem(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-from-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	count := deleteAllitems()
	json.NewEncoder(w).Encode(count)
}

func UserCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-from-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")
	params := mux.Vars(r)
	var item model.Cart
	json.NewDecoder(r.Body).Decode(&item)
	addToCart(params["id"], item)
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-from-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	params := mux.Vars(r)
	data := getcart(params["username"])
	json.NewEncoder(w).Encode(data)
}
