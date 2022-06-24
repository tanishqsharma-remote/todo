package database_dir

import (
	"database/sql"
	"time"
	"todo/model_dir"
)

func InsertUser(db *sql.DB, item model_dir.User) (sql.Result, error) {
	query := "Insert into users(username,password) values($1,$2)"

	res, er := db.Exec(query, item.Username, item.Password)
	return res, er

}
func InsertSession(db *sql.DB, sessionToken string, authorized model_dir.User, Expires time.Time) (sql.Result, error) {
	query := "insert into sessions(sessiontoken, username, expiry) VALUES ($1,$2,$3)"
	res, er := db.Exec(query, sessionToken, authorized.Username, Expires)
	return res, er

}
func InsertRefreshedSession(db *sql.DB, sessionToken string, authorized model_dir.Session, Expires time.Time) (sql.Result, error) {
	query := "insert into sessions(sessiontoken, username, expiry) VALUES ($1,$2,$3)"
	res, er := db.Exec(query, sessionToken, authorized.Username, Expires)
	return res, er

}
func InsertTask(db *sql.DB, todoTask model_dir.Todolist) (sql.Result, error) {
	query := "Insert into todolist(user_id, task, completed,archived) values($1,$2,$3,$4)"
	res, er := db.Exec(query, todoTask.UserId, todoTask.Task, todoTask.Completed, todoTask.Archived)
	return res, er
}

func GetTaskRows(db *sql.DB, userid string, pageNum string, pageSize string) (*sql.Rows, error) {
	rows, err := db.Query("with pagingCTE as(SELECT user_id,task,completed,archived, row_number() over (order by task) as rowNumber FROM todolist)select user_id,task,completed,archived from pagingCTE where user_id=$1 and rowNumber between ($2-1)*$3+1 and $2*$3", userid, pageNum, pageSize)
	return rows, err
}
func GetSession(db *sql.DB, sessionToken string) (*sql.Rows, error) {
	rows, err := db.Query("select username,expiry from sessions where sessiontoken=$1", sessionToken)
	return rows, err
}
func GetUser(db *sql.DB, credentials model_dir.Credentials) (*sql.Rows, error) {
	rows, er := db.Query("Select * from users where username=$1", credentials.Username)
	return rows, er
}

func DoneTaskQuery(db *sql.DB, Task model_dir.TodoTask) (sql.Result, error) {
	query := "update todolist set completed=true where task=$1"
	res, er := db.Exec(query, Task.Task)
	return res, er
}
func ArchiveTaskQuery(db *sql.DB, Task model_dir.TodoTask) (sql.Result, error) {
	query := "update todolist set archived=true where task=$1"
	res, er := db.Exec(query, Task.Task)
	return res, er
}

func DelSession(db *sql.DB, sessionToken string) (sql.Result, error) {
	query := "delete from sessions where sessiontoken=$1"
	res, err := db.Exec(query, sessionToken)
	return res, err
}
