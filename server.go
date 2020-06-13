package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
}

func getFilterTodosHandler(c *gin.Context) {
	// fmt.Println("start #getFilterTodosHandler")

	status := c.Query("status")
	// items := []*Todo{}
	// for _, item := range todos {
	// 	if status != "" {
	// 		if item.Status == status {
	// 			items = append(items, item)
	// 		}
	// 	} else {
	// 		items = append(items, item)
	// 	}
	// }

	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatal("Connect to database error", err)
	// }
	// defer db.Close()

	stmt, err := db.Prepare("select * from todos where status=$1 ")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmt.Query(status)
	if err != nil {
		log.Fatal(err)
	}

	items := []*Todo{}
	for rows.Next() {
		var id int
		var title, status string

		err := rows.Scan(&id, &title, &status)
		if err != nil {
			log.Fatal(err)
		}
		item := &Todo{strconv.Itoa(id), title, status}
		items = append(items, item)
		// fmt.Println("one row", id, title, status)
	}

	c.JSON(http.StatusOK, items)

	// fmt.Println("end #getFilterTodosHandler")
}

func getTodosHandler(c *gin.Context) {
	// items := []*Todo{}
	// for _,item := range todos {
	// 	items = append(items, item)
	// }
	// c.JSON(http.StatusOK, items)

	rqid := c.Param("id")
	// t, ok := todos[id]
	// if !ok {
	// 	c.JSON(http.StatusOK, gin.H{})
	// 	return
	// }

	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatal("Connect to database error", err)
	// }
	// defer db.Close()

	stmt, err := db.Prepare("select * from todos where id=$1")
	if err != nil {
		log.Fatal(err)
	}

	row := stmt.QueryRow(rqid)

	var id int
	var title, status string

	err = row.Scan(&id, &title, &status)
	if err != nil {
		log.Fatal(err)
	}
	t := &Todo{strconv.Itoa(id), title, status}

	c.JSON(http.StatusOK, t)
}

func createTodosHandler(c *gin.Context) {
	t := Todo{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// i := len(todos)
	// i++
	// id := strconv.Itoa(i)
	// t.ID = id
	// todos[id] = &t

	var err error
	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatal("Connect to database error", err)
	// }
	// defer db.Close()

	row := db.QueryRow("INSERT INTO todos (title,status) values ($1,$2) RETURNING id", t.Title, t.Status)

	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Fatal("can't scan id ", err)
		return
	}

	t.ID = strconv.Itoa(id)
	// fmt.Println("insert success id: ",id)
	c.JSON(http.StatusCreated, t)
}

func updateTodosHandler(c *gin.Context) {
	id := c.Param("id")
	t := &Todo{ID: id}
	// t := todos[id]
	if err := c.ShouldBindJSON(t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatal("Connect to database error", err)
	// }
	// defer db.Close()

	stmt, err := db.Prepare("UPDATE todos set title=$1,status=$2 where id=$3 ")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(t.Title, t.Status, id)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, t)
}

func deleteTodosHandler(c *gin.Context) {
	id := c.Param("id")
	// delete(todos, id)

	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatal("Connect to database error", err)
	// }
	// defer db.Close()

	stmt, err := db.Prepare("DELETE from todos where id=$1 ")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, "deleted todo.")
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello",
	})
}

func middleware(c *gin.Context) {
	fmt.Println("start #middleware")
	// authKey := c.GetHeader("Authorization")
	// fmt.Println(authKey)
	// if authKey != "Bearer token123" {
	// 	c.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	// 	c.Abort()
	// 	return
	// }

	c.Next()
	fmt.Println("end #middleware")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")

	//define middleware
	//apiV1.Use(middleware,middleware2)
	apiV1.Use(middleware)

	apiV1.GET("/hello", helloHandler)
	apiV1.GET("/todos", getFilterTodosHandler)
	apiV1.GET("/todos/:id", getTodosHandler)
	apiV1.POST("/todos", createTodosHandler)
	apiV1.PUT("/todos/:id", updateTodosHandler)
	apiV1.DELETE("/todos/:id", deleteTodosHandler)

	return r
}

func main() {
	r := setupRouter()
	r.Run(":1234") // listen and serve on 127.0.0.0:8080
}
