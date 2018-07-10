package tests

import (
	"testing"
	"github.com/JILeXanDR/golang/common"
	"github.com/JILeXanDR/golang/db"
	"github.com/JILeXanDR/golang/app"
)

func TestTransferMoneyWithoutBalance(t *testing.T) {

	app.CreateTest()

	var sender = 1
	var recipient = 2
	var amount = float64(1001)

	_, myErr := common.TransferMoney(sender, recipient, amount)
	if myErr == nil || myErr != common.ErrNotEnoughMoney {
		t.Error("No errors even if money is not enough")
	}
}

func TestTransferMoneyBetweenSameAccounts(t *testing.T) {

	app.CreateTest()

	var sender = 1
	var recipient = 1
	var amount = float64(20)

	_, myErr := common.TransferMoney(sender, recipient, amount)
	if myErr == nil || myErr != common.ErrSameAccounts {
		t.Error("No errors for operation between same accounts")
	}
}

func TestSuccessfulMoneyTransfer(t *testing.T) {

	app.CreateTest()

	var (
		senderId                  = 1
		recipientId               = 2
		amount                    = 100.0
		senderInitialAmount, _    = db.GetUserBalance(senderId)
		recipientInitialAmount, _ = db.GetUserBalance(recipientId)
	)

	common.TransferMoney(senderId, recipientId, amount)

	var (
		senderFinalAmount, _    = db.GetUserBalance(senderId)
		recipientFinalAmount, _ = db.GetUserBalance(recipientId)
	)

	if ((senderInitialAmount - amount) != senderFinalAmount) || ((recipientInitialAmount + amount) != recipientFinalAmount) {
		t.Error("Bad calculation")
		return
	}
}
