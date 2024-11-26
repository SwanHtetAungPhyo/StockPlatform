package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/SwanHtetAungPhyo/closure/closure"
	service "github.com/SwanHtetAungPhyo/user-service/services"
	"github.com/valyala/fasthttp"
)

func Deposit(ctx *fasthttp.RequestCtx) any {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userEmail := string(ctx.QueryArgs().Peek("email"))
	amountString := string(ctx.QueryArgs().Peek("amount"))
	
	amount, _ := strconv.ParseFloat(amountString, 64)
	newBalance, err :=service.Deposit(userEmail, amount, timeoutCtx)
	if err != nil {
		closure.JSONfiy(ctx, 400, "Deposit failed", err.Error())
		return nil
	}

	closure.JSONfiy(ctx, 200, "Deposited successfully", newBalance)
	return nil 
}

func Read(ctx *fasthttp.RequestCtx) any {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userEmail := string(ctx.QueryArgs().Peek("email"))

	balance, err := service.GetBalance(userEmail, timeoutCtx)
	if err != nil {
		closure.JSONfiy(ctx, 400, "Failed to read balance", err.Error())
		return nil
	}

	closure.JSONfiy(ctx, 200, "User balance is", balance)
	return nil
}