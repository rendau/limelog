package mongo

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ConfigId = 1

func (d *St) ConfigSet(ctx context.Context, config *entities.ConfigSt) error {
	var err error

	collection := d.Db.Collection("config")

	upsert := true

	res := collection.FindOneAndUpdate(ctx, bson.M{"id": ConfigId}, bson.M{
		"$set": config,
	}, &options.FindOneAndUpdateOptions{Upsert: &upsert})
	if err = res.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			return d.handleErr(err)
		}
	}

	return nil
}

func (d *St) ConfigGet(ctx context.Context) (*entities.ConfigSt, error) {
	var err error

	collection := d.Db.Collection("config")

	res := collection.FindOne(ctx, bson.M{"id": ConfigId})
	if err = res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, d.handleErr(err)
	}

	result := &entities.ConfigSt{}

	err = res.Decode(result)
	if err != nil {
		return nil, d.handleErr(err)
	}

	return result, nil
}
