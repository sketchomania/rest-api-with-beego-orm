package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	Id       int    `orm:"column(id);pk;auto_add"`
	Name     string `orm:"column(name);size(50)"`
	Email    string `orm:"column(email);size(100)"`
	Password string `orm:"column(password);size(100)"`
	// AddressId   []*UserAddress    `orm:"reverse(many)"`
}

type UserAddress struct {
	Id     int   `orm:"column(id);pk;"`
	UserId *User `orm:"column(user_id);rel(fk);"`
	Line1  string
}

func init() {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "my-secret-pw"
	dbName := "go_demo"

	orm.RegisterModel(new(User))
	orm.RegisterDriver(dbDriver, orm.DRMySQL)
	orm.RegisterDataBase("default", dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?charset=utf8")
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All Users Endpoint Hit")
	o := orm.NewOrm()

	var users []User
	u, err := o.QueryTable("user").All(&users)
	if err == nil {
		fmt.Printf("Result Nums: %d\n", u)
		for _, user := range users {
			fmt.Println(user.Id, user.Name)
		}
	}

	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Create User Endpoint Hit")
	o := orm.NewOrm()

	var user User
	vars := mux.Vars(r)
	user.Name = vars["name"]
	user.Email = vars["email"]
	user.Password = vars["password"]

	id, err := o.Insert(&user)
	// fmt.Println("ID: %d, ERR: %v\n", id, err)
	fmt.Println(id, err)

	// json.NewEncoder(w).Encode(user)
	fmt.Fprintf(w, "New User Sucessfully created")
}

func getSingleUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Single User Endpoint Hit")
	o := orm.NewOrm()

	var user User
	vars := mux.Vars(r)
	user.Name = vars["name"]

	// user := User{Id: 1}
	err := o.Read(&user)
	fmt.Println("ERR: \n", err, user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	// user := User{Name: "Name Updated"}
	// u, err := o.Update(&user)
	// fmt.Println("u:", u, "err: ", err)

	user := User{Id: 1}
	if o.Read(&user) == nil {
		user.Name = "MyName"
		if num, err := o.Update(&user); err == nil {
			fmt.Println(num, user)
		}
	}
	fmt.Fprintf(w, "User updated sucessfully")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	var user User
	vars := mux.Vars(r)
	user.Name = vars["name"]

	num, err := o.Delete(&user)
	fmt.Println("num:", num, "err: ", err)

	fmt.Fprintf(w, "User Sucessfully Deleted")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/", hello).Methods("GET")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/user/{name}", getSingleUser).Methods("GET")
	router.HandleFunc("/user/{name}/{email}/{password}", createUser).Methods("POST")
	router.HandleFunc("/user/{name}", deleteUser).Methods("DELETE")
	router.HandleFunc("/user/{name}/{email}", updateUser).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func main() {
	fmt.Println("GO REST API WITH BEEGO ORM")

	handleRequests()
}
