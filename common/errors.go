package common

import "errors"

var ErrSameAccounts = errors.New("Could not transfer money between same accounts")
var ErrNotEnoughMoney = errors.New("Not enough money")
