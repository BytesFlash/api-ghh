package mail

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/imedcl/manager-api/pkg/config"
	"github.com/imedcl/manager-api/pkg/data"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

type SensorEmail struct {
	Code  string `json:"code"`
	Glosa string `json:"glosa"`
}

func getTemplate(name string) string {
	dir, _ := os.Getwd()
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/pkg/mail/templates/%s.html", dir, name))
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(string(content))
	return string(content)
}

func SendUserRegister(user *data.User) (bool, error) {
	cfg := config.New()
	from := mail.NewEmail("Autentia Manager", "no-reply@autentia.io")
	subject := "Bienvenido a Autentia Manager"
	to := mail.NewEmail(user.Name, user.Email)
	plainTextContent := fmt.Sprintf("Bienvenido/a al sistema de Autentia Manager, para confirmar tu correo e ingresar tu contraseña, por favor ingresa al siguiente link: %s%s/%s", cfg.AppUrl(), "validate/account", user.Token)
	html := getTemplate("confirmation")
	html = strings.Replace(html, "{{username}}", user.Name, 1)
	html = strings.Replace(html, "{{email}}", user.Email, 1)
	html = strings.Replace(html, "{{nick_name}}", user.NickName, 1)
	html = strings.Replace(html, "{{link}}", fmt.Sprintf("%s%s/%s", cfg.AppUrl(), "validate/account", user.Token), 1)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, html)
	client := sendgrid.NewSendClient(cfg.SendgridKey())
	response, err := client.Send(message)
	if err != nil {
		return false, err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return true, nil
	}
}

func SendRecoveryPassword(user *data.User) (bool, error) {
	cfg := config.New()
	from := mail.NewEmail("Autentia Manager", "no-reply@autentia.io")
	subject := "Recuperar contraseña - Autentia Manager"
	to := mail.NewEmail(user.Name, user.Email)
	plainTextContent := fmt.Sprintf("Hola! se ha solicitado la recuperación de contraseña de tu cuenta, para realizar el cambio de contraseña ingresa en el siguiente link: %s%s/%s, Si no solicitaste esta recuperación puedes ignorar este correo", cfg.AppUrl(), "recovery", user.Token)
	html := getTemplate("recovery")
	html = strings.Replace(html, "{{username}}", user.Name, 1)
	html = strings.Replace(html, "{{email}}", user.Email, 1)
	html = strings.Replace(html, "{{link}}", fmt.Sprintf("%s%s/%s", cfg.AppUrl(), "recovery", user.Token), 1)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, html)
	client := sendgrid.NewSendClient(cfg.SendgridKey())
	response, err := client.Send(message)
	if err != nil {
		return false, err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return true, nil
	}
}

func SendFeedback(user *data.User, url string, browser string, system string, comment string) (bool, error) {
	cfg := config.New()
	from := mail.NewEmail("Autentia Manager", "no-reply@autentia.io")
	subject := "Feedback - Autentia Manager"
	to := mail.NewEmail("Autentia Manager", cfg.ContactEmail())
	plainTextContent := fmt.Sprintf("Hola!\n\n El usuario %s (%s) ha enviado feedback de Autentia Manager, su comentario es: \n\n %s\n\n El comentario fue enviado desde: \n\n  URL: %s\n Navegador: %s\n Sistema Operativo: %s", user.Name, user.Email, comment, url, browser, system)
	html := getTemplate("feedback")
	html = strings.Replace(html, "{{user}}", user.Name, 1)
	html = strings.Replace(html, "{{url}}", url, 1)
	html = strings.Replace(html, "{{navigator}}", browser, 1)
	html = strings.Replace(html, "{{system}}", system, 1)
	html = strings.Replace(html, "{{message}}", comment, 1)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, html)
	client := sendgrid.NewSendClient(cfg.SendgridKey())
	response, err := client.Send(message)
	if err != nil {
		return false, err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return true, nil
	}
}

func SendActivateUser(user *data.User) (bool, error) {
	cfg := config.New()
	from := mail.NewEmail("Autentia Manager", "no-reply@autentia.io")
	subject := "Activar cuenta - Autentia Manager"
	to := mail.NewEmail(user.Name, user.Email)
	plainTextContent := fmt.Sprintf("Hola! se ha solicitado la activación de tu cuenta, para realizar la activación ingresa en el siguiente link: %s%s/%s", cfg.AppUrl(), "activate", user.Token)
	html := getTemplate("activate")
	html = strings.Replace(html, "{{username}}", user.Name, 1)
	html = strings.Replace(html, "{{email}}", user.Email, 1)
	html = strings.Replace(html, "{{linkActivate}}", fmt.Sprintf("%s%s/%s", cfg.AppUrl(), "activate", user.Token), 1)
	html = strings.Replace(html, "{{linkPassword}}", fmt.Sprintf("%s%s/%s", cfg.AppUrl(), "recovery", user.Token), 1)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, html)
	client := sendgrid.NewSendClient(cfg.SendgridKey())
	response, err := client.Send(message)
	if err != nil {
		return false, err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return true, nil
	}
}

func SendSensorNotRegistered(user *data.User, sensors []SensorEmail) (bool, error) {
	var listSensor string
	for _, sensor := range sensors {
		listSensor = listSensor + sensor.Code + "-" + sensor.Glosa + "<br>"
	}
	cfg := config.New()
	from := mail.NewEmail("Autentia Manager", "no-reply@autentia.io")
	subject := "Registro de lectores Autentia Manager"
	to := mail.NewEmail(user.Name, user.Email)
	plainTextContent := fmt.Sprintf("Termino el proceso de registro de lectores, Lectores no registrados: %s", sensors)
	html := getTemplate("sensor")
	html = strings.Replace(html, "{{username}}", user.Name, 1)
	html = strings.Replace(html, "{{email}}", user.Email, 1)
	html = strings.Replace(html, "{{title}}", "Registro de lectores", 1)
	html = strings.Replace(html, "{{register}}", "ha terminado el proceso de registro de lectores en", 1)
	html = strings.Replace(html, "{{sensor}}", "Detalle de los lectores: ", 1)
	html = strings.Replace(html, "{{sensors}}", listSensor, 1)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, html)
	client := sendgrid.NewSendClient(cfg.SendgridKey())
	response, err := client.Send(message)
	if err != nil {
		return false, err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return true, nil
	}
}

func SendSensorNotUpdated(user *data.User, sensors []SensorEmail) (bool, error) {
	var listSensor string
	for _, sensor := range sensors {
		listSensor = listSensor + sensor.Code + "-" + sensor.Glosa + "<br>"
	}
	cfg := config.New()
	from := mail.NewEmail("Autentia Manager", "no-reply@autentia.io")
	subject := "Actualizar lectores Autentia Manager"
	to := mail.NewEmail(user.Name, user.Email)
	plainTextContent := fmt.Sprintf("Termino el proceso de actualizar lectores, Lectores no registrados: %s", sensors)
	html := getTemplate("sensor")
	html = strings.Replace(html, "{{username}}", user.Name, 1)
	html = strings.Replace(html, "{{email}}", user.Email, 1)
	html = strings.Replace(html, "{{title}}", "Actualización de lectores", 1)
	html = strings.Replace(html, "{{register}}", "ha terminado el proceso de actualización de lectores en", 1)
	html = strings.Replace(html, "{{sensor}}", "Detalle de los lectores: ", 1)
	html = strings.Replace(html, "{{sensors}}", listSensor, 1)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, html)
	client := sendgrid.NewSendClient(cfg.SendgridKey())
	response, err := client.Send(message)
	if err != nil {
		return false, err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return true, nil
	}
}
