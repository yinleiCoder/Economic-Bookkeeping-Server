package main

import (
	"crypto/sha512"
	"economic_srvs/user_srv/model"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func main() {
	dsn := "root:yl13795950539@tcp(127.0.0.1:3306)/economic_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode("yinlei19980505", options)
	passwordSaltedStr := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	// mock some users into database
	for i := 0; i < 10; i++ {
		user := model.User{
			Mobile:   fmt.Sprintf("1379595053%d", i),
			Password: passwordSaltedStr,
			NickName: fmt.Sprintf("yinlei%d", i),
		}
		db.Save(&user)
	}
}
