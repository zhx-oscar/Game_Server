package DB

import (
	_ "Cinder/Base/ServerConfig"
	"context"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	if err := initMongoDB(); err != nil {
		log.Error("Init mongo err ", err)
	}
}

///////////////////////////////////////////////////////////////////////////////

// MongoDB相关
var MongoDB *mongo.Database

func initMongoDB() error {
	addr := viper.GetString("MongoDB.Addr")
	database := viper.GetString("MongoDB.DataBase")
	username := viper.GetString("MongoDB.User")
	password := viper.GetString("MongoDB.Password")

	clientInfo := options.Client()
	clientInfo.ApplyURI(fmt.Sprintf("mongodb://%s%s/",
		func() string {
			if username == "" {
				return ""
			} else {
				return username + ":" + password + "@"
			}
		}(),
		addr))

	client, err := mongo.Connect(context.Background(), clientInfo)
	if err != nil {
		return err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return err
	}

	MongoDB = client.Database(database)
	if MongoDB == nil {
		return errors.New("no database")
	}

	return nil
}
