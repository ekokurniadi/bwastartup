package handler

import (
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Index(c *gin.Context) {
	session := sessions.Default(c)
	data := session.Get("message")
	session.Set("message", "")
	session.Save()
	c.HTML(http.StatusOK, "user_index.html", gin.H{"data": data})
}

func (h *userHandler) New(c *gin.Context) {
	c.HTML(http.StatusOK, "user_create.html", nil)
}

func (h *userHandler) Create(c *gin.Context) {
	var input user.FormCreateUserInput

	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusOK, "user_create.html", input)
		return
	}

	registerInput := user.RegisterUserInput{}
	registerInput.Name = input.Name
	registerInput.Email = input.Email
	registerInput.Occupation = input.Occupation
	registerInput.Password = input.Password

	_, err = h.userService.RegisterUser(registerInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Create User Success")
	session.Save()
	c.Redirect(http.StatusFound, "/users")

}

func (h *userHandler) Edit(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	registeredUser, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := user.FormUpdateUserInput{}
	input.ID = registeredUser.ID
	input.Name = registeredUser.Name
	input.Email = registeredUser.Email
	input.Occupation = registeredUser.Occupation

	c.HTML(http.StatusOK, "user_edit.html", input)
}

func (h *userHandler) Update(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	var input user.FormUpdateUserInput
	err := c.ShouldBind(&input)

	if err != nil {
		input.Error = err
		c.HTML(http.StatusOK, "user_edit.html", input)
		return
	}

	input.ID = id

	_, err = h.userService.UpdateUser(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Update User Success")
	session.Save()
	c.Redirect(http.StatusFound, "/users")
}

func (h *userHandler) NewAvatar(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	c.HTML(http.StatusOK, "user_upload.html", gin.H{"ID": id})
}

func (h *userHandler) CreateAvatar(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	file, err := c.FormFile("avatar")

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	userID := id
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Upload User Avatars Success")
	session.Save()
	c.Redirect(http.StatusFound, "/users")
}
