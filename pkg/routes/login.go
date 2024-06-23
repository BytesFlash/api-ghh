package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/imedcl/manager-api/pkg/config"
	"github.com/imedcl/manager-api/pkg/data"
	"github.com/imedcl/manager-api/pkg/events"
	"github.com/imedcl/manager-api/pkg/mail"
)

type loginParams struct {
	User struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
}

type oauthParams struct {
	Token string `form:"token" binding:"required"`
}

type responseValidation struct {
	Status bool `json:"status"`
}
type googleUserInfo struct {
	Id            string `form:"id" json:"id"`
	Email         string `form:"email" json:"email"`
	VerifiedEmail bool   `form:"verified_email" json:"verified_email"`
	Name          string `form:"name" json:"name"`
	Given         string `form:"given_name" json:"given_name"`
	FamilyName    string `form:"family_name" json:"family_name"`
	Picture       string `form:"picture" json:"picture"`
	Locale        string `form:"locale" json:"locale"`
	Hd            string `form:"hd" json:"hd"`
	Error         struct {
		Code    int    `form:"code" json:"code"`
		Message string `form:"message" json:"message"`
		Status  string `form:"status" json:"status"`
	} `form:"error" json:"error"`
}

type recoveryParams struct {
	Email string `form:"email" binding:"required"`
}

// @Summary login
// @Description login
// @Tags login
// @Accept json
// @Produce json
// @Param login body loginParams true "login"
// @Success 200 "Inició sesión correctamente"
// @failure 404 {object} MessageResponse
// @Router /login [post]
func LoginRoute(router *gin.Engine, db *data.DB, client *unleash.Client) {

	router.POST("/login", func(c *gin.Context) {
		var params loginParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Usuario - Contraseña no es una combinación válida",
				Code:    http.StatusBadRequest,
			})
			return
		}
		emailName := strings.ToLower(params.User.Email)
		getNickName, errEmail := db.GetUserEmail(emailName)
		user, err := db.GetUser(getNickName.NickName)

		if errEmail != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: "Usuario - Contraseña no es una combinación válida",
				Code:    http.StatusNotFound,
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: "Usuario - Contraseña no es una combinación válida",
				Code:    http.StatusNotFound,
			})
			return
		}

		if user.StatusUser == config.STATUS_INACTIVE || user.StatusUser == config.STATUS_LOCKED {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Su cuenta está desactivada, contáctese con el administrador",
				Code:    http.StatusBadRequest,
			})
			return
		}

		userLastPasswordDate, err := db.GetLastPasswordDate(user.ID)
		if err == nil {
			if time.Now().Sub(userLastPasswordDate.CreatedAt).Hours()/24/30 >= 3 {
				c.JSON(http.StatusBadRequest, MessageResponse{
					Details: "Su contraseña ha caducado. Por favor, vaya a la página de inicio y seleccione la opción '¿Olvidó su contraseña?' para restablecerla.",
					Code:    http.StatusBadRequest,
				})
				return
			}
		}
		PrtyParams, _ := events.PrettyParams(params)
		event := &events.EventLog{
			UserNickname: user.NickName,
			Resource:     "Login",
			Params:       PrtyParams,
		}
		jwt, err := db.Login(user, params.User.Password)
		if err != nil {
			event.Event = "Intento de login fallido"
			c.JSON(http.StatusUnauthorized, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusUnauthorized,
			})
		} else {
			event.Event = "Login realizado correctamente"
			c.Header("Authorization", fmt.Sprintf("Bearer %s", jwt))
			c.JSON(http.StatusOK, MessageResponse{
				Details: "Inició sesión correctamente",
				Code:    http.StatusOK,
			})
		}
		event.Write()
	})
}

// @Summary login guest
// @Description login guest
// @Tags login
// @Success 200 "Inició sesión correctamente"
// @failure 401 {object} MessageResponse
// @Router /users/guest [post]
func LoginGuestRoute(router *gin.Engine, db *data.DB, client *unleash.Client) {

	router.POST("/users/guest", func(c *gin.Context) {
		user, err := db.GetUser("guest")

		event := &events.EventLog{
			UserNickname: user.NickName,
			Resource:     "Login",
		}
		jwt, err := db.Login(user, "guest")
		if err != nil {
			event.Event = "Intento de login fallido"
			c.JSON(http.StatusUnauthorized, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusUnauthorized,
			})
		} else {
			event.Event = "Login realizado correctamente"
			c.Header("Authorization", fmt.Sprintf("Bearer %s", jwt))
			c.JSON(http.StatusOK, MessageResponse{
				Details: "Inició sesión correctamente",
				Code:    http.StatusOK,
			})
		}
		event.Write()
	})

}

// @Summary Recovery
// @Description recovery pass
// @Tags login
// @Accept json
// @Produce json
// @Param recovery body recoveryParams true "recovery"
// @Success 201 "Se enviará correo para recuperar su contraseña a la dirección example@mail.com"
// @failure 404 {object} MessageResponse
// @Router /users/recovery [post]
func RecoveryRoute(router *gin.Engine, db *data.DB, client *unleash.Client) {

	// Recovery Password
	router.POST("/users/recovery", func(c *gin.Context) {
		var params recoveryParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		email := strings.ToLower(params.Email)
		user, err := db.GetUserEmail(email)
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: "Usuario no existe, o no se encuentra en Base de datos",
				Code:    http.StatusNotFound,
			})
			return
		}
		_, errEmail := db.GetUserReccovery(email)
		if errEmail != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: "La dirección de correo no coincide con la registrada para este usuario",
				Code:    http.StatusNotFound,
			})
			return
		}

		if user.StatusUser != config.STATUS_ACTIVE {
			c.JSON(http.StatusUnauthorized, MessageResponse{

				Details: "Usuario inactivo",
				Code:    http.StatusUnauthorized,
			})
			return
		}
		user.Token = uuid.New().String()
		user, err = db.UpdateUser(email, user, "")
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		mail.SendRecoveryPassword(user)
		PrtyParams, _ := events.PrettyParams(params)
		event := &events.EventLog{
			UserNickname: user.NickName,
			Resource:     "Recuperar contraseña",
			Event:        "Solicitud de recuperación de contraseña exitosa",
			Params:       PrtyParams,
		}
		event.Write()
		message := fmt.Sprintf("Se enviará correo para recuperar su contraseña a la dirección %s", params.Email)
		c.JSON(http.StatusCreated, message)
	})

}

// @Summary Reconfirmation
// @Description user reconfirmation
// @Tags login
// @Accept json
// @Produce json
// @Success 201 {object} response
// @failure 404 {object} MessageResponse
// @Router /users/{nickName}/reconfirmation [post]
// @param nickName path string true "nickName"
func ReconfirmationRoute(router *gin.Engine, db *data.DB, client *unleash.Client) {

	// Resend Confirmation Mail
	router.POST("/users/:nickName/reconfirmation", func(c *gin.Context) {
		nickName := strings.ToLower(c.Param("nickName"))
		user, err := db.GetUser(nickName)
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		user.Token = uuid.New().String()
		user, err = db.UpdateUser(nickName, user, "")
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		mail.SendUserRegister(user)
		PrtyParams, _ := events.PrettyParams(c.Params)
		event := &events.EventLog{
			UserNickname: nickName,
			Resource:     "Reconfirmación",
			Event:        "Reconfirmación de usuario exitosa",
			Params:       PrtyParams,
		}
		event.Write()
		c.JSON(http.StatusCreated, response{Data: user})
	})
}

// @Summary Confirm
// @Description confirm user
// @Tags login
// @Accept json
// @Produce json
// @Param recovery body validateParams true "confirm"
// @Success 200 "Confirmación de usuario exitosa"
// @failure 400 {object} MessageResponse
// @Router /users/confirm [post]
func ConfirmRoute(router *gin.Engine, db *data.DB, client *unleash.Client) {
	router.POST("/users/confirm", func(c *gin.Context) {
		var params validateParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		if params.Token == "" || params.Password == "" {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Parametros incorrectos",
				Code:    http.StatusBadRequest,
			})
			return
		}
		user, err := db.GetUserByToken(params.Token)
		if err != nil || (user.Validated && !params.Recovery) {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Token no encontrado o el usuario ya se encuentra validado",
				Code:    http.StatusBadRequest,
			})
			return
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
		user.Validated = true
		user.Token = uuid.New().String()
		nickName := strings.ToLower(user.NickName)
		user, err = db.UpdateUser(nickName, user, params.Password)
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		PrtyParams, _ := events.PrettyParams(params)
		event := &events.EventLog{
			UserNickname: nickName,
			Resource:     "Confirmación",
			Event:        "Confirmación de usuario exitosa",
			Params:       PrtyParams,
		}
		event.Write()
		c.JSON(http.StatusOK, response{Data: user})
	})
}

// @Summary Validate
// @Description validate user
// @Tags login
// @Accept json
// @Produce json
// @Param validate body userValidateParams true "validate"
// @Success 200 {object} responseValidation
// @failure 400 {object} MessageResponse
// @Router /users/validate [post]
func ValidateRoute(router *gin.Engine, db *data.DB, client *unleash.Client) {
	router.POST("/users/validate", func(c *gin.Context) {
		var params userValidateParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		if params.Token == "" {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Parametros incorrectos",
				Code:    http.StatusBadRequest,
			})
			return
		}
		_, err := db.GetUserByToken(params.Token)
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusOK, responseValidation{Status: true})
	})
}

// @Summary activate
// @Description activate user
// @Tags login
// @Accept json
// @Produce json
// @Param activate body userValidateParams true "activate"
// @Success 200 {object} responseValidation
// @failure 400 {object} MessageResponse
// @Router /users/activate [post]
func ActivateRoute(router *gin.Engine, db *data.DB, client *unleash.Client) {
	router.POST("/users/activate", func(c *gin.Context) {
		var params userValidateParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		if params.Token == "" {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Parametros incorrectos",
				Code:    http.StatusBadRequest,
			})
			return
		}
		user, err := db.GetUserByToken(params.Token)
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}
		user.StatusUser = config.STATUS_ACTIVE
		nickName := strings.ToLower(user.NickName)
		user, err = db.UpdateUser(nickName, user, "")
		c.JSON(http.StatusOK, responseValidation{Status: true})
	})
}
