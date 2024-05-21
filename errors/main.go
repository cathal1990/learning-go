package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
)

var ErrorUserNotFound error = errors.New("Could not find user")

type User struct {
	ID int
	Username string
}

func fakeGetUserFromDb(ctx context.Context, id int64) (User, error) {
	if id == 1 {
		return User{
			ID: 1,
			Username: "Cathal",
		}, nil
	} else {
		return User{}, sql.ErrNoRows
	}

}

func fakeFetchUserProfile(ctx context.Context, id int64) (User, error) {
	user, err := fakeGetUserFromDb(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errors.Join(err, ErrorUserNotFound)
		}
	}
	return user,nil
}

func main() {
	user, err := fakeFetchUserProfile(context.Background(), 2)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	fmt.Printf("%+v\n", user)
}