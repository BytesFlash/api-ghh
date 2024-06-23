package routes

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/imedcl/manager-api/pkg/config"
	"github.com/imedcl/manager-api/pkg/data"
	"github.com/imedcl/manager-api/pkg/events"
	"github.com/imedcl/manager-api/pkg/mail"
)

type userStruct struct {
	Active    bool   `json:"active"`
	Email     string `json:"email"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	NickName  string `json:"nick_name"`
	Picture   string `json:"picture"`
	Status    string `json:"status_user"`
	Validated bool   `json:"validated"`
	ActiInst  bool   `json:"acti_inst"`
}

type response struct {
	Data *data.User `json:"data"`
}

type feedbackResponse struct {
	Status bool `json:"status"`
}

type feedBackParams struct {
	Message string `form:"message" json:"message"`
	Url     string `form:"url" json:"url"`
	Browser string `form:"browser" json:"browser"`
	System  string `form:"system" json:"system"`
}
type validateParams struct {
	Token    string `form:"token"`
	Password string `form:"password"`
	Recovery bool   `form:"recovery"`
}

type userParams struct {
	NickName    string `form:"nickname" binding:"required"`
	Name        string `form:"name" binding:"required"`
	Email       string `form:"email" binding:"required"`
	Password    string `form:"password"`
	Description string `form:"description"`
	Country     string `form:"country" binding:"required"`
	Dni         string `form:"dni" binding:"required"`
	ActiInst    bool   `form:"acti_inst" json:"acti_inst"`
}

type updateUserParams struct {
	NickName    string `form:"nickname"`
	Name        string `form:"name"`
	Email       string `form:"email"`
	Password    string `form:"password"`
	Description string `form:"description"`
	Country     string `form:"country"`
	Dni         string `form:"dni"`
	StatusUser  string `form:"status_user" json:"status_user" `
	ActiInst    bool   `form:"acti_inst" json:"acti_inst"`
}

type userValidateParams struct {
	Token string `form:"token"`
}

// @Summary users
// @Description users
// @Tags user
// @security barerToken
// @Accept json
// @Produce json
// @Success 200 {object} []userStruct
// @failure 404 {object} MessageResponse
// @Router /users [get]
func UsersRolesRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	router.GET("/users", func(c *gin.Context) {

		var users []*data.User

		var userDataList []userStruct
		for _, userInfo := range users {
			var user userStruct
			user.Active = *&userInfo.Active
			user.Email = *&userInfo.Email
			user.ID = *&userInfo.ID
			user.Name = *&userInfo.Name
			user.NickName = *&userInfo.NickName
			user.Picture = *&userInfo.Picture
			user.Status = *&userInfo.StatusUser
			user.Validated = *&userInfo.Validated
			user.ActiInst = *&userInfo.ActiInst

			userDataList = append(userDataList, user)
		}

		c.JSON(http.StatusOK, userDataList)
	})
}

// @Summary user
// @Description create user
// @Tags user
// @security barerToken
// @Accept json
// @Produce json
// @Param user body userParams true "user"
// @Success 201 "Usuario Creado"
// @failure 400 {object} MessageResponse
// @Router /users [post]
func UserPostRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	router.POST("/users", func(c *gin.Context) {

		var params userParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		name := strings.Title(params.Name)
		email := strings.ToLower(params.Email)
		dni := strings.ToLower(params.Dni)
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		nickName := strings.ToLower(params.NickName)
		if db.UserAllDBExists(nickName) {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Usuario ya existe en la base de datos",
				Code:    http.StatusBadRequest,
			})
			return
		}
		if !emailRegex.MatchString(email) {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "El email no es válido",
				Code:    http.StatusBadRequest,
			})
			return
		}
		if len(params.Name) <= 3 {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "El Nombre debe tener más de 3 carácteres",
				Code:    http.StatusBadRequest,
			})
			return
		}
		if len(dni) == 0 {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Dni no es válido",
				Code:    http.StatusBadRequest,
			})
			return
		}
		if len(params.Description) > 500 {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "La descripción no debe tener más de 500 carácteres",
				Code:    http.StatusBadRequest,
			})
			return
		}
		if len(nickName) <= 3 {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "El usuario debe tener más de 3 carácteres",
				Code:    http.StatusBadRequest,
			})
			return
		}
		nickNameRegex := regexp.MustCompile(`^[{a-zA-Z-}]+$`)
		if !nickNameRegex.MatchString(nickName) {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Usuario incorrecto, solo puede tener letras y guiones",
				Code:    http.StatusBadRequest,
			})
			return
		}

		if db.UserExistsEmail(email) {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "El email ya se encuentra registrado",
				Code:    http.StatusBadRequest,
			})
			return
		}
		user := &data.User{
			Name:        name,
			Email:       email,
			Password:    params.Password,
			Description: params.Description,
			NickName:    nickName,
			Dni:         dni,
			ActiInst:    params.ActiInst,
		}
		db.CreateUser(user)
		mail.SendUserRegister(user)
		PrtyParams, _ := events.PrettyParams(user)
		event := &events.EventLog{
			UserNickname: currentUser.NickName,
			Resource:     "Usuarios",
			Event:        fmt.Sprintf("Se ha creado el usuario %s", nickName),
			Params:       PrtyParams,
		}
		event.Write()
		c.JSON(http.StatusCreated, "Usuario Creado")
	})
}

// @Summary user instrospection
// @Description instrospection
// @Tags user
// @security barerToken
// @Accept json
// @Produce json
// @Success 200 {object} response
// @Router /users/introspection [get]
func UserInstrospectionRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {

	// Get User By Token
	router.GET("/users/introspection", func(c *gin.Context) {
		c.JSON(http.StatusOK, response{Data: currentUser})
	})
}

// @Summary user get
// @Description user get
// @Tags user
// @security barerToken
// @Accept json
// @Produce json
// @Success 200 {object} response
// @failure 404 {object} MessageResponse
// @Router /users/{userIdentifier} [get]
// @param userIdentifier path string true "userIdentifier"
func UserGetRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	// Get User
	router.GET("/users/:userIdentifier", func(c *gin.Context) {

		userIdentifier := strings.ToLower(c.Param("userIdentifier"))
		_, err := uuid.Parse(userIdentifier)
		var user *data.User
		if err != nil {
			user, err = db.GetUser(userIdentifier)
		} else {
			user, err = db.GetUserByID(userIdentifier)
		}
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusOK, response{Data: user})
	})
}

// @Summary update user
// @Description update user
// @Tags user
// @security barerToken
// @Accept json
// @Produce json
// @Success 200 "El usuario se ha actualizado con éxito"
// @failure 400 {object} MessageResponse
// @Router /users/{nickName} [put]
// @Param user body updateUserParams true "user"
// @param nickName path string true "userIdentifier"
func UserPutRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	// Update User
	router.PUT("/users/:nickName", func(c *gin.Context) {
		nickName := strings.ToLower(c.Param("nickName"))
		if !db.UserExists(nickName) {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Usuario no encontrado",
				Code:    http.StatusBadRequest,
			})
			return
		}
		var params updateUserParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		newNickName := strings.ToLower(params.NickName)
		if newNickName != "" {
			if len(newNickName) <= 3 {
				c.JSON(http.StatusBadRequest, MessageResponse{
					Details: "El usuario debe tener más de 3 carácteres",
					Code:    http.StatusBadRequest,
				})
				return
			}
			nickNameRegex := regexp.MustCompile(`^[{a-zA-Z-}]+$`)
			if !nickNameRegex.MatchString(newNickName) {
				c.JSON(http.StatusBadRequest, MessageResponse{
					Details: "Usuario incorrecto, solo puede tener letras y guiones",
					Code:    http.StatusBadRequest,
				})
				return
			}

		}

		name := strings.Title(params.Name)
		if params.Name != "" {
			if len(name) <= 3 {
				c.JSON(http.StatusBadRequest, MessageResponse{
					Details: "El Nombre debe tener más de 3 carácteres",
					Code:    http.StatusBadRequest,
				})
				return
			}
		}
		email := strings.ToLower(params.Email)
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if params.Email != "" {
			if !emailRegex.MatchString(email) {
				c.JSON(http.StatusBadRequest, MessageResponse{
					Details: "El email no es válido",
					Code:    http.StatusBadRequest,
				})
				return
			}
		}

		userId, _ := db.UserExistsName(nickName)
		userEmail, _ := db.UserExistEmail(email)
		if userEmail.ID != userId.ID {
			if db.UserExistsEmail(email) {
				c.JSON(http.StatusBadRequest, MessageResponse{
					Details: "El email ya se encuentra registrado",
					Code:    http.StatusBadRequest,
				})
				return
			}

		}

		if len(params.Description) > 500 {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "La descripción no debe tener más de 500 carácteres",
				Code:    http.StatusBadRequest,
			})
			return
		}
		dni := strings.ToLower(params.Dni)
		user := &data.User{
			NickName:    newNickName,
			Name:        name,
			Email:       email,
			StatusUser:  params.StatusUser,
			Description: params.Description,
			Dni:         dni,
			ActiInst:    params.ActiInst,
		}
		if params.Password != "" {
			passValidate, passError := config.Password(params.Password)
			if !passValidate {
				c.JSON(http.StatusBadRequest, PasswordMessageResponse{
					Details:      "La contraseña no es válida",
					Code:         http.StatusBadRequest,
					Requirements: passError,
				})
				return
			}
			result := db.PasswordExists(params.Password, user.ID)
			if result == nil {
				c.JSON(http.StatusBadRequest, MessageResponse{
					Details: "No puede registrar una contraseña anteriormente utilizada",
					Code:    http.StatusBadRequest,
				})
				return
			} else {
				db.CreatePassword(user.Password, user)
			}
		}
		if params.StatusUser != "" {
			listArray := contains(config.ALL_STATUS, user.StatusUser)
			if !listArray {
				c.JSON(http.StatusBadRequest, MessageResponse{
					Details: "Código de estado incorrecto",
					Code:    http.StatusBadRequest,
				})
				return
			}
		}
		_, err := db.UpdateUser(nickName, user, params.Password)
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		PrtyParams, _ := events.PrettyParams(params)
		event := &events.EventLog{
			UserNickname: currentUser.NickName,
			Resource:     "Usuarios",
			Event:        fmt.Sprintf("Se ha actualizado el usuario %s", nickName),
			Params:       PrtyParams,
		}
		event.Write()

		c.JSON(http.StatusOK, "El usuario se ha actualizado con éxito")

	})
}

// @Summary user feedback
// @Description user feedback
// @Tags user
// @security barerToken
// @Accept json
// @Produce json
// @Success 201 {object} feedbackResponse
// @failure 400 {object} MessageResponse
// @Router /users/feedback [post]
// @Param feedback body feedBackParams true "feedback"
func UserFeedbackRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	// Send Feedback Mail
	router.POST("/users/feedback", func(c *gin.Context) {
		var params feedBackParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		mail.SendFeedback(currentUser, params.Url, params.Browser, params.System, params.Message)
		c.JSON(http.StatusCreated, feedbackResponse{Status: true})
	})
}

// @Summary user logout
// @Description user logout
// @Tags user
// @security barerToken
// @Accept json
// @Produce json
// @Success 200 {object} MessageResponse
// @failure 400 "Usuario no activo"
// @Router /logout [delete]
func LogoutRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	// Logout User
	router.DELETE("/logout", func(c *gin.Context) {
		if !db.IsActive(currentUser.ID) {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Usuario no activo",
				Code:    http.StatusBadRequest,
			})
			return
		} else {
			db.ExpireToken(currentUser)
			c.JSON(http.StatusOK, MessageResponse{Details: "Sesión expirada", Code: http.StatusOK})
			event := &events.EventLog{
				UserNickname: currentUser.NickName,
				Resource:     "Logout",
				Event:        "Successful Logout",
			}
			event.Write()
		}

	})
}

// @Summary user delete
// @Description user delete
// @Tags user
// @security barerToken
// @Accept json
// @Produce json
// @Success 200 {object} MessageResponse
// @failure 400 {object} MessageResponse
// @Router /users/{nickName} [delete]
// @param nickName path string true "nickName"
func DeleteUserRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	// Delete User
	router.DELETE("/users/:nickName", func(c *gin.Context) {

		nickName := strings.ToLower(c.Param("nickName"))

		user, err := db.GetUser(nickName)
		if err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Usuario no Registrado",
				Code:    http.StatusBadRequest,
			})
			return
		}

		user.StatusUser = config.STATUS_LOCKED

		_, updateErr := db.UpdateUser(nickName, user, "")

		if updateErr != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Error al bloquear usuario, intente nuevamente",
				Code:    http.StatusBadRequest,
			})
			return
		}
		PrtyParams, _ := events.PrettyParams(c.Params)
		event := &events.EventLog{
			UserNickname: currentUser.NickName,
			Resource:     "Usuarios",
			Event:        fmt.Sprintf("Se ha bloqueado al usuario %s", nickName),
			Params:       PrtyParams,
		}
		event.Write()
		c.JSON(http.StatusOK, MessageResponse{Details: "Usuario bloqueado exitosamente", Code: http.StatusOK})
	})

}
