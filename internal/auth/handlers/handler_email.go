package handlers

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/smtp"

	"github.com/ezep02/rodeo/internal/auth/models"
)

func (h *AuthHandler) SendResetUserPasswordEmailHandler(rw http.ResponseWriter, r *http.Request) {

	var (
		smtpHost string = "smtp.gmail.com"
		// smtpPort := "587"
		sender   string   = "epereyra443@gmail.com"
		password string   = "cubrrxzypaskawzc"
		to       []string = []string{"pereyraezequiel15617866@outlook.es"}
		u        models.UserEmail
	)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Couldn't parse request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err := json.Unmarshal(b, &u); err != nil {
		http.Error(rw, "Error al deserializar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// user, err := h.AuthServ.SearchUserByEmail(h.ctx, u.Email)

	// if err != nil {
	// 	http.Error(rw, "Si el correo es válido, recibirás un email con instrucciones.", http.StatusInternalServerError)
	// 	return
	// }

	// crear un token utilizando los datos de user
	// tokenString, err := jwt.GenerateToken(user.ID, user.Is_admin, user.Name, user.Email, user.Surname, user.Phone_number, user.Is_barber, time.Now().Add(15*time.Minute))

	// if err != nil {
	// 	http.Error(rw, "[Creacion token] Algo salio mal, vuelve a intentarlo mas tarde", http.StatusInternalServerError)
	// 	return
	// }

	// Autenticación con el servidor
	auth := smtp.PlainAuth("", sender, password, smtpHost)

	// Crear conexión segura con TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// Establecer conexión con el servidor SMTP
	conn, err := tls.Dial("tcp", smtpHost+":465", tlsConfig) // Usa puerto 465 para TLS directo
	if err != nil {
		log.Fatal("Error en conexión TLS:", err)
	}
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Fatal("Error creando cliente SMTP:", err)
	}

	// Autenticarse
	if err = client.Auth(auth); err != nil {
		log.Fatal("Error en autenticación:", err)
	}

	// Configurar el remitente y destinatario
	if err = client.Mail(sender); err != nil {
		log.Fatal(err)
	}
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			log.Fatal(err)
		}
	}

	// Escribir el mensaje
	wc, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}

	// msg := fmt.Sprintf("Subject: 🔐 Recupera tu contraseña\r\n"+
	// 	"MIME-Version: 1.0\r\n"+
	// 	"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
	// 	"\r\n"+
	// 	"<html><body>"+
	// 	"<h2>🔐 Recuperación de contraseña</h2>"+
	// 	"<p>Hola,</p>"+
	// 	"<p>Has solicitado restablecer tu contraseña. Haz clic en el botón de abajo:</p>"+
	// 	"<a href='http://localhost:5173/auth/recover/token=%s' "+
	// 	"style='display:inline-block;background-color:#007bff;color:#ffffff;padding:10px 20px;text-decoration:none;border-radius:5px;'>Restablecer contraseña</a>"+
	// 	"<p>Si no solicitaste esto, ignora este mensaje.</p>"+
	// 	"<p>Saludos,<br>Equipo de Soporte</p>"+
	// 	"</body></html>", tokenString)

	// _, err = wc.Write([]byte(msg))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Cerrar conexión
	client.Quit()

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("ok")
}
