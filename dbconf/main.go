package main

import (
	"fmt"

	"github.com/Users/patrickfurtak/desktop/go-gallery/models"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "patrickfurtak"
	dbname = "gogallery"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	// us.DestructiveReset()
	user, err := us.ByID(2)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}
