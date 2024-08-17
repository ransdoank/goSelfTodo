# # Author
## @[Amran Tirtana](https://github.com/ransdoank)

# # Preparing started
```
fork from repository && git pull from fork
install MySql in local/server
make file '.env'
add connection database
DB_USERNAME=yourusername
DB_PASSWORD=yourpassword
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=yourdatabase
PORT=8012 #This optional, as default 8080
```

# # packages for this project:
```
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get github.com/joho/godotenv
go get github.com/stretchr/testify

```


# # running apps
```
go run main.go
```

# # test endpoint
### CLI Base
- Create a Task:
```
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"Title":"Task 1","Description":"Description 1","Status":"Waiting List","DueDate":"2024-08-20T00:00:00Z"}' 
```
- Read All Tasks:
```
curl 'http://localhost:8080/tasks?status=Waiting%20List&due_date=2024-08-21'
```
- Update a Task:
```
curl -X PUT http://localhost:8080/tasks/1 -H "Content-Type: application/json" -d '{"Status":"In Progress"}'
```
- Delete a Task:
```
curl -X DELETE http://localhost:8080/tasks/1
```

### Postman 
Link: [Documentation API](https://documenter.getpostman.com/view/3057077/2sA3s9BnaS)

Or Import manual from
```
goSelfTodo.postman_collection.json
```

### Automation Test
```
go test -v
```
if you wan to check feedback invalid test, you can uncheck in main_test.go section
```
TestInvalidCreateTask
TestUpdateNonexistentTask
TestDeleteNonexistentTask
```