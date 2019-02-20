package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

var (
	Db          *sqlx.DB
	RedisClient *redis.Client
)

type User struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
}

type Category struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Task struct {
	ID         int        `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	User       User       `db:"user" json:"user"`
	Category   Category   `db:"category" json:"category"`
	StatusID   int        `db:"status_id" json:"status_id"`
	EndDate    *time.Time `db:"end_date" json:"end_date"`
	PublishFlg bool       `db:"publish_flg" json:"publish_flg"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

func main() {
	e := echo.New()

	e.GET("/tasks", tasksHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":8080"))
}

func tasksHandler(c echo.Context) error {
	// Redis接続
	RedisClient := GetRedisClient()
	defer RedisClient.Close()

	const redisKey = "tasks"
	tasks := []Task{}

	if cache, _ := RedisClient.Get(redisKey).Result(); cache == "" {
		fmt.Printf("*** no cache key:%s ***\n", redisKey)

		// MySQL接続
		Db = InitDB()
		defer Db.Close()

		sql := `
		SELECT
			t.id, t.name, t.status_id, t.end_date, t.publish_flg,
			u.id AS "user.id", u.name AS "user.name", u.email AS "user.email",
			c.id AS "category.id", c.name AS "category.name"
		FROM tasks t 
		JOIN users u ON t.user_id = u.id
		JOIN categories c ON t.category_id = c.id
		`
		err := Db.Select(&tasks, sql)
		checkErr(err)

		// キャッシュ登録
		data, err := json.Marshal(tasks)
		checkErr(err)
		err = RedisClient.Set(redisKey, data, time.Hour*1).Err()
		checkErr(err)

	} else {
		fmt.Printf("*** cache hit key:%s ***\n", redisKey)

		err := json.Unmarshal([]byte(cache), &tasks)
		checkErr(err)

	}

	return c.JSON(http.StatusOK, &Tasks{
		Tasks: tasks,
	})
}

func InitDB() *sqlx.DB {
	db, err := sqlx.Connect("mysql", "root:root@tcp(mysql:3306)/testdb?charset=utf8&parseTime=true")
	checkErr(err)
	return db
}

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
