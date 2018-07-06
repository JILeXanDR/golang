package common

import (
	"errors"
	"github.com/JILeXanDR/golang/db"
	"github.com/jinzhu/gorm"
)

// transfer money from one user to another user
func TransferMoney(fromId int, toId int, amount float64) (internalErr error, myErr error) {

	if fromId == toId {
		return nil, errors.New("Could not transfer money beetwen same accounts")
	}

	// get sender user
	sender := &db.User{}
	err := db.Connection.Where(&db.User{Identifier: fromId}).First(sender).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("Sender user does not exist")
	} else if err != nil {
		return err, nil
	}

	// check balance
	if sender.Balance < amount {
		return nil, errors.New("Not enough money")
	}

	// get recipient user
	recipient := &db.User{}
	err = db.Connection.Where(&db.User{Identifier: toId}).First(recipient).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("Recipient user does not exist")
	} else if err != nil {
		return err, nil
	}

	transaction := db.Connection.Begin()

	sender.Balance -= amount
	err = db.Connection.Save(sender).Error
	if err != nil {
		transaction.Rollback()
		return err, nil
	}

	recipient.Balance += amount
	err = db.Connection.Save(recipient).Error
	if err != nil {
		transaction.Rollback()
		return err, nil
	}

	transaction.Commit()

	return nil, nil
}
