package handler

import (
	"context"
	"time"

	"github.com/SwanHtetAungPhyo/closure/closure"
	"github.com/SwanHtetAungPhyo/user-service/models"
	"github.com/SwanHtetAungPhyo/user-service/services"
	"github.com/valyala/fasthttp"
)

func SignUpHandler(ctx *fasthttp.RequestCtx) interface{} {
	var user models.User

	if err := closure.Binder(ctx, &user); err != nil {
		closure.JSONfiy(ctx, fasthttp.StatusBadRequest, "Binding Error", err.Error())
		return nil
	}

	err := service.SignUp(&user)
	if err != nil {
		closure.JSONfiy(ctx, fasthttp.StatusInternalServerError, "SignUp Error", err.Error())
		return nil
	}
	user.Password = ""
	closure.JSONfiy(ctx, fasthttp.StatusCreated, "User Created", user)
	return nil
}



func SignInHandler(ctx *fasthttp.RequestCtx) interface{} {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User


	if err := closure.Binder(ctx, &user); err != nil {
		closure.JSONfiy(ctx, fasthttp.StatusBadRequest, "Binding Error", err.Error())
		return nil
	}
	authToken, err := service.SignIn(&user, timeoutCtx)
	if err != nil {
		if err == context.DeadlineExceeded {
			closure.JSONfiy(ctx, fasthttp.StatusRequestTimeout, "Request Timeout", "The sign-in operation took too long")
		} else {
			closure.JSONfiy(ctx, fasthttp.StatusUnauthorized, "SignIn Error", err.Error())
		}
		return nil
	}


	closure.JSONfiy(ctx, fasthttp.StatusOK, "Signed In", map[string]string{
		"auth_token": authToken,
	})
	return nil
}
