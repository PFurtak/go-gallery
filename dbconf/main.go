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
	us.DestructiveReset()
	user := models.User{
		Name:  "Michael Scott",
		Email: "michael@dundermifflin.com",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}

	if err := us.Delete(user.ID); err != nil {
		panic(err)
	}

	userByID, err := us.ByID(user.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(userByID)
}
