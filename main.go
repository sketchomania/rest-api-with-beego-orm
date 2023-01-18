package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	Id        int       `orm:"column(Id);pk;auto"`
	UserName  string    `orm:"column(UserName);size(50);unique"`
	Email     string    `orm:"column(Email);size(100);unique"`
	Password  string    `orm:"column(Password);size(100)"`
	CreatedAt time.Time `orm:"column(CreatedAt)"`
	UpdatedAt time.Time `orm:"column(UpdatedAt)"`
	// AddressId   []*UserAddress    `orm:"reverse(many)"`
}

func (u *User) TableName() string {
	return "user"
}

type Post struct {
	Id        int       `orm:"column(Id);pk;auto"`
	Title     string    `orm:"column(Title);size(100)"`
	Content   string    `orm:"column(Content);size(255)"`
	AuthorId  *User     `orm:"column(AuthorId);rel(fk)"`
	CreatedAt time.Time `orm:"column(CreatedAt)"`
	UpdatedAt time.Time `orm:"column(UpdatedAt)"`
}

func (p *Post) TableName() string {
	return "post"
}

// func (p *Post) Validate() error {

// 	if p.Title == "" {
// 		return errors.New("Required Title")
// 	}
// 	if p.Content == "" {
// 		return errors.New("Required Content")
// 	}
// 	if p.AuthorId < 1 {
// 		return errors.New("Required Author")
// 	}
// 	return nil
// }

// 1t1
// 1tm
// type UserPosts struct {}
// type UserAddress struct {
// 	Id     int   `orm:"column(id);pk;"`
// 	UserId *User `orm:"column(user_id);rel(fk);"`
// 	Line1  string
// }

func init() {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "my-secret-pw"
	dbName := "go_demo"

	orm.RegisterModel(new(User), new(Post))
	orm.RegisterDriver(dbDriver, orm.DRMySQL)
	orm.RegisterDataBase("default", dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?charset=utf8")
}

// fetch all
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All Users Endpoint Hit \n")
	o := orm.NewOrm()

	var users []User
	num, err := o.QueryTable(new(User)).All(&users)
	if err == nil {
		fmt.Printf("Result Nums: %d\n", num)
		// for _, user := range users {
		// 	fmt.Println(user.Id, user.CreatedAt, user.UserName, user.Email, user.Password, user.UpdatedAt)
		// }
	}

	// w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	json.NewEncoder(w).Encode(users)
}
func getAllposts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All Posts Endpoint Hit \n")
	o := orm.NewOrm()

	var posts []Post
	num, err := o.QueryTable(new(Post)).All(&posts)
	if err == nil {
		fmt.Printf("Result Nums: %d\n", num)
		// for _, post := range posts {
		// 	fmt.Println(post.Id, post.CreatedAt, post.Title, post.Content, post.AuthorId, post.UpdatedAt)
		// }
	}

	json.NewEncoder(w).Encode(posts)
}

// Create
func createUser(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	var user User
	// vars := mux.Vars(r)
	// user.UserName = vars["name"]
	// user.Email = vars["email"]
	// user.Password = vars["password"]

	// fmt.Println("Body: ", r.Body)
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	// fmt.Println("user: ", user)

	id, err := o.Insert(&user)
	fmt.Printf("ID: %d, ERR: %v\n", id, err)

	fmt.Fprintf(w, "New User Sucessfully created 🎉")
	json.NewEncoder(w).Encode(user)
}
func createPost(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	var user User
	orm.NewOrm().QueryTable("user").Filter("id", 2).One(&user)

	var post Post
	_ = json.NewDecoder(r.Body).Decode(&post)
	// user type or user id only
	post.AuthorId = &user
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	// fmt.Println("post: ", post)

	id, err := o.Insert(&post)
	fmt.Printf("ID: %d, ERR: %v\n", id, err)

	fmt.Fprintf(w, "New Post Sucessfully created 🎉")
	json.NewEncoder(w).Encode(post)
}

// Single user and post
func getSingleUser(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	var user User
	vars := mux.Vars(r)
	userId := vars["id"]
	// user.UserName = vars["name"]

	if err := o.QueryTable("user").Filter("Id", userId).One(&user); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// err := o.Read(&user)
	// fmt.Println("ERR: \n", err)
	fmt.Fprintf(w, "Single User Endpoint Hit \n")
	json.NewEncoder(w).Encode(user)
}
func getSinglePost(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	var post Post
	params := mux.Vars(r)
	postId := params["id"]
	// fmt.Println("post: \n", post)

	err := o.QueryTable("post").Filter("Id", postId).One(&post)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Single Post Endpoint Hit \n")
	json.NewEncoder(w).Encode(post)
}

// Update
func updateUser(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	var user User
	params := mux.Vars(r)
	userId := params["id"]

	if err := o.QueryTable("user").Filter("Id", userId).One(&user); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if o.Read(&user) == nil {
		_ = json.NewDecoder(r.Body).Decode(&user)
		user.UpdatedAt = time.Now()
		if num, err := o.Update(&user); err == nil {
			fmt.Println(num, user)
		}
	}
	fmt.Fprintf(w, "User updated sucessfully \n")
	json.NewEncoder(w).Encode(user)
}
func updatePost(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	var post Post
	params := mux.Vars(r)
	postId := params["id"]

	if err := o.QueryTable("post").Filter("Id", postId).One(&post); err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if o.Read(&post) == nil {
		_ = json.NewDecoder(r.Body).Decode(&post)
		post.UpdatedAt = time.Now()
		if num, err := o.Update(&post); err == nil {
			fmt.Println(num, post)
		}
	}
	fmt.Fprintf(w, "Post updated sucessfully \n")
	json.NewEncoder(w).Encode(post)
}

// Delete
func deleteUser(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	var user User
	params := mux.Vars(r)
	userId := params["id"]

	if err := o.QueryTable("user").Filter("Id", userId).One(&user); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	num, err := o.Delete(&user)
	fmt.Println("Affected rows: ", num, "err: ", err)

	fmt.Fprintf(w, "User Sucessfully Deleted \n")
	json.NewEncoder(w).Encode(user)
}
func deletePost(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()

	var post Post
	params := mux.Vars(r)
	postId := params["id"]

	if err := o.QueryTable("post").Filter("Id", postId).One(&post); err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	num, err := o.Delete(&post)
	fmt.Println("Affected rows: ", num, "err: ", err)

	fmt.Fprintf(w, "Post Sucessfully Deleted \n")
	json.NewEncoder(w).Encode(post)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/", hello).Methods("GET")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/posts", getAllposts).Methods("GET")
	router.HandleFunc("/user/{id}", getSingleUser).Methods("GET")
	router.HandleFunc("/post/{id}", getSinglePost).Methods("GET")

	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/post", createPost).Methods("POST")

	router.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/post/{id}", deletePost).Methods("DELETE")

	router.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/post/{id}", updatePost).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8081", router))
}

func main() {
	fmt.Println("GO REST API WITH BEEGO ORM")
	fmt.Println("🚀 Listening on port http://localhost:8081/")
	handleRequests()
}
