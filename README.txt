e.GET("/users", getAllUsers)        Example: localhost:8080/users
e.GET("/users/:id", getUser)        Example: localhost:8080/users/1
e.POST("/users", createUser)        Example: localhost:8080/users?name=Vladislav
e.PUT("/users/:id", updateUser)     Example: localhost:8080/users/1?name=Alexander
e.DELETE("/users/:id", deleteUser)  Example: localhost:8080/users/1