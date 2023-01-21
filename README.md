# Rest api design with Beego orm

## Tools
- mux-router
- beego-orm
- go-sql-driver

## Relations
- One-to-one   [ Post - User ]
- One-to-many  [ User - Post ]
- Many-to-many [ Post - Tag  ]

## Routes
- send data body in json form
```go
	router := mux.NewRouter()
	router.HandleFunc("/", hello).Methods("GET")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/posts", getAllposts).Methods("GET")
	router.HandleFunc("/user/{id}", getSingleUser).Methods("GET")
	router.HandleFunc("/post/{id}", getSinglePost).Methods("GET")

	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/postbyuser/{id}", createPost).Methods("POST")

	router.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/post/{id}", deletePost).Methods("DELETE")

	router.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/post/{id}", updatePost).Methods("PUT")
```
