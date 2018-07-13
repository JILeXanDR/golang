package repositories

import "github.com/JILeXanDR/golang/app/db"

func GetUserBalance(userId int) (res float64, err error) {
	var user = &db.User{}
	err = db.Connection.Where(&db.User{Identifier: userId}).First(user).Error
	if err != nil {
		return 0, err
	}

	return user.Balance, nil
}
