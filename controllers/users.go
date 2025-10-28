package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/amghazanfari/pryx/context"
	"github.com/amghazanfari/pryx/views"
	"github.com/gin-gonic/gin"

	"github.com/amghazanfari/pryx/models"
)

type Users struct {
	Templates struct {
		Signup         Template
		Signin         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
		ModelList      Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
	EndpointService      *models.EndpointService
}

func (u Users) SignUp(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.Signup.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrng", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setUserSessionCookie(w, session.Token)

	fmt.Println("session set for user for signup: %w", session)

	http.Redirect(w, r, "/user/me", http.StatusFound)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.Signin.Execute(w, r, data)
}

func (u Users) Authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setUserSessionCookie(w, session.Token)

	fmt.Printf("session set for user for signin: %s\n", session.Token)

	http.Redirect(w, r, "/user/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("you opened current user")
	ctx := r.Context()
	user := context.User(ctx)

	fmt.Fprintf(w, "Logged in as: %s\n", user.Email)
}

func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("session") // Read from the new cookie
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.Delete(tokenCookie.Value)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	deleteUserSessionCookie(w)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := c.Request.Cookie("session")
		if err != nil {
			c.Next()
			return
		}
		fmt.Printf("session found for user: %s\n", tokenCookie.Value)

		user, err := umw.SessionService.User(tokenCookie.Value)
		if err != nil {
			c.Next()
			return
		}

		// Store user in Gin context
		c.Set("user", user)

		// Continue to next handler
		c.Next()
	}
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("you came to middleware for require user")
		user := context.User(r.Context())
		fmt.Println(user)
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	vals := url.Values{
		"token": {pwReset.Token},
	}

	resetURL := "http://localhost:8080/reset-pw?" + vals.Encode()

	err = u.EmailService.ForgotPassword(data.Email, resetURL)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	err = u.UserService.Update(user.ID, data.Password)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)

		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setUserSessionCookie(w, session.Token)
	http.Redirect(w, r, "/user/me", http.StatusFound)
}

func (u Users) ModelList(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Endpoints *[]models.Endpoint
		CsrfField template.HTML
	}
	var err error
	data.Endpoints, err = u.EndpointService.List()

	csrfField := r.Context().Value(views.CsrfFieldKey).(string)
	data.CsrfField = template.HTML(csrfField)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	u.Templates.ModelList.Execute(w, r, data)
}
