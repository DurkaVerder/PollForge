package storage

import (
	"auth/internal/models"
	"database/sql"
	"fmt"
	"os"
)

var Db *sql.DB
func ConnectToDb()error{
		dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	var err error

	Db, err = sql.Open("postgres", dsn)
	if err != nil{
		return err
	}
	
	return nil
}

func CompareData(user models.User)(){
	
	
}

func InsertData(hashedPassword string,request models.RegisterRequest)error{
	var UserId int
	err := Db.QueryRow(`INSERT INTO users (name,email,password)
	                          VALUE($1,$2,$3) RETURNING id`,
							  request.Username, request.Email, string(hashedPassword)).Scan(&UserId)
	if err != nil{
		
	}
}