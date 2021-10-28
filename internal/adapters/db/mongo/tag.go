package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *St) TagSet(ctx context.Context, value string) error {
	var err error

	collection := d.Db.Collection("tag")

	upsert := true

	res := collection.FindOneAndUpdate(ctx, bson.M{"value": value}, bson.M{
		"$set": bson.M{"value": value},
	}, &options.FindOneAndUpdateOptions{Upsert: &upsert})
	if err = res.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			return d.handleErr(err)
		}
	}

	return nil
}

func (d *St) TagList(ctx context.Context) ([]string, error) {
	var err error

	collection := d.Db.Collection("tag")

	cur, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Sort: bson.M{"value": 1},
	})
	if err != nil {
		return nil, d.handleErr(err)
	}
	defer cur.Close(ctx)

	result := make([]string, 0)

	obj := map[string]string{}

	for cur.Next(ctx) {
		err = cur.Decode(&obj)
		if err != nil {
			return nil, d.handleErr(err)
		}

		result = append(result, obj["value"])
	}
	if err = cur.Err(); err != nil {
		return nil, d.handleErr(err)
	}

	return result, nil
}

func (d *St) TagRemove(ctx context.Context, value string) error {
	collection := d.Db.Collection("tag")

	_, err := collection.DeleteOne(ctx, bson.M{"value": value})
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return d.handleErr(err)
		}
	}

	return nil
}
