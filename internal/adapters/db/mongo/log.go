package mongo

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *St) LogCreate(ctx context.Context, obj map[string]interface{}) error {
	var err error

	collection := d.Db.Collection("log")

	_, err = collection.InsertOne(ctx, obj)
	if err != nil {
		return d.handleErr(err)
	}

	return nil
}

func (d *St) LogList(ctx context.Context, pars *entities.LogListParsSt) ([]map[string]interface{}, error) {
	collection := d.Db.Collection("log")

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, d.handleErr(err)
	}
	defer cur.Close(ctx)

	result := make([]map[string]interface{}, 0)

	for cur.Next(ctx) {
		obj := map[string]interface{}{}

		err = cur.Decode(obj)
		if err != nil {
			return nil, d.handleErr(err)
		}

		result = append(result, obj)
	}
	if err = cur.Err(); err != nil {
		return nil, d.handleErr(err)
	}

	return result, nil
}

// func (d *St) LogGet(ctx context.Context, id string, pars *entities.LogGetParsSt) (*entities.LogSt, error) {
// 	var err error
//
// 	ops := &options.FindOneOptions{}
//
// 	if pars != nil {
// 		if pars.Projection != nil {
// 			ops.Projection = *pars.Projection
// 		}
// 	}
//
// 	collection := d.Db.Collection("log")
//
// 	res := collection.FindOne(ctx, bson.M{"id": id}, ops)
// 	if err = res.Err(); err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, nil
// 		}
// 		return nil, d.handleErr(err)
// 	}
//
// 	result := &entities.LogSt{}
//
// 	err = res.Decode(result)
// 	if err != nil {
// 		return nil, d.handleErr(err)
// 	}
//
// 	return result, nil
// }
