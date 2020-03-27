package model

import (
	"context"
    "fmt"
    "net/http"
    "encoding/json"
    "bytes"
    "io/ioutil"
    // "strings"

	"github.com/mongodb/mongo-go-driver/mongo"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)


type InfraDB struct {
    Host    string  `json:"host"`
    Id      string  `json:"id"`
    Name    string  `json:"name"`
    Port    string  `json:"port"`
    Status  string  `json:"status"`
    Type    string  `json:"type"`
}


type InfraDBReq struct {
    Type  string    `json:"type"`
    Host  string    `json:"host"`
    Port  string    `json:"port"`
    Name  string    `json:"name"`
}


type InfraDBResp struct {
    Data InfraDB    `json:"data"`
    Desc string     `json:"desc"`
}


// DB is a mongodb database
var DB *mongo.Database

func touchDB(baseURL string, dbHost string) string {
    databaseName := fmt.Sprintf("face_recognition_backend_db_%s", dbHost)
    queryURL := fmt.Sprintf("%s?name=%s", baseURL, databaseName)
	resp, err := http.Get(queryURL)
    log.Info("start touchDB: %v %v", databaseName, queryURL)

    if err != nil {
		panic(err)
    } else if resp.StatusCode == http.StatusNotFound {
        // create database if not exist
        createURL := baseURL

        pc := &InfraDBReq{
            Type: "mongo",
            Host: dbHost,
            Port: "17017",
            Name: databaseName,
        }

        jsonValue, _ := json.Marshal(pc)
	    cresp, err := http.Post(createURL, "application/json",  bytes.NewBuffer(jsonValue))
        if err != nil {
            panic(err)
        }
        if cresp.StatusCode == http.StatusOK {
            return fmt.Sprintf("mongodb://%s:17017", dbHost)
        } else {
            panic(fmt.Errorf("check infrastructure for database create"))
        }

    } else if resp.StatusCode == http.StatusOK {
	    b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            panic(err)
        }
        var p InfraDBResp
        err = json.Unmarshal(b, &p)
        if err != nil {
            panic(err)
        }
        return fmt.Sprintf("mongodb://%s:%s", p.Data.Host, p.Data.Port)

    } else {
        panic(fmt.Errorf("check infrastructure for database"))
    }

    return ""
}

// InitDB initalizes mongo db
func InitDB() {
	cfg := config.Config()
	dbAddr := cfg.GetString("db-addr")

    // create database?
    // TODO(test)
    appid := cfg.GetString("appid")
    if appid != "" {
        // checkout if database exists and addr
        gatewayAddr := cfg.GetString("gateway-addr")
        baseURL := fmt.Sprintf("%s/api/v1/infrastructure/database", gatewayAddr)
        // host := strings.Trim(strings.Split(gatewayAddr, ":")[1], "/")
        host := cfg.GetString("hostip")
        dbAddr = touchDB(baseURL, host)
    }

	client, err := mongo.NewClient(dbAddr)
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	DB = client.Database(cfg.GetString("app-name"))
}
