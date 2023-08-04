package service

import (
	"LiveChat/db"
	"LiveChat/util"
	"go.mongodb.org/mongo-driver/bson"
)

// loginAuth 登录验证
func loginAuth(conn *DatabaseConn, userName string, password string) (bool, error) {
	data, err := db.Find(&db.DatabaseConn{Conn: conn.Conn, Name: conn.Name}, "users", bson.M{
		"user_name": userName,
	})

	var msg error = nil
	if err != nil {
		msg = err
	}

	return data[0]["password"] == util.MD5(password), msg
}

// registerAuth 注册验证
func registerAuth(conn *DatabaseConn, userName string, password string) (bool, error) {
	data, err := db.Find(&db.DatabaseConn{Conn: conn.Conn, Name: conn.Name}, "users", bson.M{
		"user_name": userName,
	})

	var msg error = nil
	if err != nil {
		msg = err
	}

	if len(data) != 0 {
		return false, msg
	} else {
		db.Add(&db.DatabaseConn{Conn: conn.Conn, Name: conn.Name}, "users",
			bson.D{{Key: "user_name", Value: userName}, {Key: "password", Value: util.MD5(password)}},
		)
		return true, msg
	}
}
