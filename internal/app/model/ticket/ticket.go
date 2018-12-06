package ticket

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"math/rand"

	"github.com/google/uuid"

	"github.com/mongodb/mongo-go-driver/mongo"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
)

func collection() *mongo.Collection {
	return model.DB.Collection("ticket")
}

// EncodeTicket translate from UserTicket to string
func EncodeTicket(t *schema.UserTicket) (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(t)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// DecodeTicket translate from string to UserTicket
func DecodeTicket(ticketStr string) (*schema.UserTicket, error) {
	by, err := base64.StdEncoding.DecodeString(ticketStr)
	if err != nil {
		return nil, err
	}

	t := &schema.UserTicket{}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// FindTicket return ticket in database
func FindTicket(userName string, nonceStr string) error {
	_, err := collection().Find(context.Background(), map[string]string{
        "username": userName, 
        "nonce_str": nonceStr,
    })
    return err
}

// CreateTicket create ticket after user login with passwd
func CreateTicket(userName string) (string, error) {
	uid, _ := uuid.NewUUID()
	dbTicket := &schema.UserTicket{
		UserName:  userName,
		ID:        uid.String(),
		Timestamp: util.NowMilli(),
		NonceStr:  RandStr(32),
	}

	_, err := collection().InsertOne(context.Background(), dbTicket)
	if err != nil {
		return "", err
	}

	var ticketStr string
	ticketStr, err = EncodeTicket(dbTicket)
	if err != nil {
		return "", err
	}
	return ticketStr, nil
}

// DropTicket drop ticket in db
func DropTicket(t string) error {
	dbTicket, err := DecodeTicket(t)
	if err != nil {
		return err
	}
	_, err = collection().DeleteOne(context.Background(), map[string]string{"_id": dbTicket.ID})
	return err
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStr create random strings
func RandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
