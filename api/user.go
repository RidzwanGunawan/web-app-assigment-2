package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserTaskCategory(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	var recordUser = model.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	recordUser, err := u.userService.Register(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse("register success"))
}

func (u *userAPI) Login(c *gin.Context) {
	// TODO: answer here
	var loginInput model.UserLogin
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid JSON input"))
		return
	}

	if loginInput.Email == "" || loginInput.Password == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("email or password is empty"))
		return
	}

	loginCredentials := &model.User{
		Email:    loginInput.Email,
		Password: loginInput.Password,
	}

	sessionToken, err := u.userService.Login(loginCredentials)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("internal server error"))
		return
	}

	expiration := time.Now().Add(5 * time.Minute)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_token",
		Value:    *sessionToken,
		Expires:  expiration,
	})

	c.JSON(http.StatusOK, model.NewSuccessResponse("login successful"))

}

func (u *userAPI) GetUserTaskCategory(c *gin.Context) {
	// TODO: answer here
	taskCategories, err := u.userService.GetUserTaskCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("internal server error"))
		return
	}

	c.JSON(http.StatusOK, taskCategories)
}
