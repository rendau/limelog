package mongo

import (
	"context"

	"github.com/rendau/limelog/internal/cns"
	"github.com/rendau/limelog/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *St) LogCreate(ctx context.Context, obj interface{}) error {
	var err error

	collection := d.Db.Collection("log")

	_, err = collection.InsertOne(ctx, obj)
	if err != nil {
		return d.handleErr(err)
	}

	return nil
}

func (d *St) LogCreateMany(ctx context.Context, objs []interface{}) error {
	var err error

	collection := d.Db.Collection("log")

	_, err = collection.InsertMany(ctx, objs)
	if err != nil {
		return d.handleErr(err)
	}

	return nil
}

func (d *St) LogList(ctx context.Context, pars *entities.LogListParsSt) ([]map[string]interface{}, int64, error) {
	var err error

	collection := d.Db.Collection("log")

	filter := bson.M{}

	if pars.FilterObj != nil {
		for k, v := range pars.FilterObj {
			filter[k] = v
		}
	}

	var totalCnt int64

	if pars.PageSize > 0 {
		totalCnt, err = collection.CountDocuments(ctx, filter)
		if err != nil {
			return nil, 0, d.handleErr(err)
		}
	}

	skip := pars.Page * pars.PageSize

	cur, err := collection.Find(ctx, filter, &options.FindOptions{
		Sort:  bson.M{"sf_ts": -1},
		Skip:  &skip,
		Limit: &pars.PageSize,
	})
	if err != nil {
		return nil, 0, d.handleErr(err)
	}
	defer cur.Close(ctx)

	result := make([]map[string]interface{}, 0)

	for cur.Next(ctx) {
		obj := map[string]interface{}{}

		err = cur.Decode(obj)
		if err != nil {
			return nil, 0, d.handleErr(err)
		}

		result = append(result, obj)
	}
	if err = cur.Err(); err != nil {
		return nil, 0, d.handleErr(err)
	}

	return result, totalCnt, nil
}

func (d *St) LogRemove(ctx context.Context, pars *entities.LogRemoveParsSt) error {
	collection := d.Db.Collection("log")

	filter := bson.M{}

	if pars.TsLt != nil {
		filter[cns.SfTsFieldName] = bson.M{
			"$lte": pars.TsLt.UnixMilli(),
		}
	}

	if len(filter) == 0 {
		return nil
	}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return d.handleErr(err)
	}

	return nil
}

func (d *St) LogListDistinctTag(ctx context.Context) ([]string, error) {
	collection := d.Db.Collection("log")

	dbResult, err := collection.Distinct(ctx, cns.SfTagFieldName, bson.D{})
	if err != nil {
		return nil, d.handleErr(err)
	}

	result := make([]string, 0, len(dbResult))

	for _, dbV := range dbResult {
		strV, ok := dbV.(string)
		if !ok {
			d.lg.Warnw("Strange type of value", "value", dbV)
			continue
		}

		result = append(result, strV)
	}

	return result, nil
}
