package mongo

import (
	"context"
	"net/url"
	"time"

	"github.com/mechta-market/limelog/internal/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type St struct {
	debug bool
	lg    interfaces.Logger

	Client *mongo.Client
	Db     *mongo.Database
}

func New(
	lg interfaces.Logger,
	username string,
	password string,
	host string,
	dbName string,
	replicaSet string,
	debug bool,
) (*St, error) {
	uri := url.URL{
		Scheme: "mongodb",
		Host:   host,
	}

	if username != "" {
		uri.User = url.UserPassword(username, password)
	}

	uri.Path = "/" + dbName

	uvs := uri.Query()
	if replicaSet != "" {
		uvs.Set("replicaSet", replicaSet)
	}
	uri.RawQuery = uvs.Encode()

	ops := options.Client().ApplyURI(uri.String())

	ops.SetConnectTimeout(15 * time.Second)

	client, err := mongo.Connect(context.Background(), ops)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	res := &St{
		debug: debug,
		lg:    lg,

		Client: client,
		Db:     db,
	}

	return res, nil
}

func (d *St) handleErr(err error) error {
	if err == nil {
		return err
	}

	d.lg.Errorw("Mongo error", err)

	return err
}
