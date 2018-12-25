package controller

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
    "errors"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/ticket"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/user"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// LoginPOST handle camera list
func LoginPOST(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	var iu schema.User
	err = json.Unmarshal(b, &iu)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

    log.Infof("test: recv %v", iu)

	var du *schema.User
	du, err = user.FindUser(iu.UserName)
	if err == user.ErrUserNotFound {
		Error(w, err, http.StatusNotFound)
		return
	} else if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

    log.Infof("test: find %v", du)

	if du.Password != iu.Password {
		Error(w, err, http.StatusForbidden)
		return
	}

	// create cookie with ticket
	var t string
	t, err = ticket.CreateTicket(du.UserName)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

    log.Infof("test: create %v", t)

	http.SetCookie(w, &http.Cookie{
		Name:    "ticket",
		Value:   t,
		Expires: time.Now().Add(86400 * 30 * time.Second), // TODO configuration?
	})

    log.Infof("return: create %v", t)

	w.Header().Set("Content-Type", "application/text")
	io.WriteString(w, "hello world")
}

// LogoutPOST set empty login cookie
func LogoutPOST(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "ticket",
		Value:   "",
		Expires: time.Now().Add(86400 * 30 * time.Second), // TODO configuration?
	})
	w.Header().Set("Content-Type", "application/text")
	io.WriteString(w, "Goodbye")
}

// CheckLoginSession check every api call with cookie
func CheckLoginSession(w http.ResponseWriter, r *http.Request) error {
	ticketCookie, err := r.Cookie("ticket")
    if ticketCookie != nil {
        if ticketCookie.Value == "" {
            err = errors.New("Cookies not set")
        } else {
            t, err := ticket.DecodeTicket(ticketCookie.Value)
            if err == nil {
                err = ticket.FindTicket(t.UserName, t.NonceStr)
            }
        }
    } else {
        err = errors.New("Cookies not set")
    }

    if err != nil {
        Error(w, err, 499)
    }

	// TODO check login session life-time
	return err
}
