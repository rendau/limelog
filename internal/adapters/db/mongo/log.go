package mongo

// func (d *St) LogCreate(ctx context.Context, obj *entities.LogCUSt, upsert bool) error {
// var err error
//
// collection := d.Db.Collection("log")
//
// res := collection.FindOneAndUpdate(ctx, bson.M{"id": id}, bson.M{
// "$set": obj,
// }, &options.FindOneAndUpdateOptions{Upsert: &upsert})
// if err = res.Err(); err != nil {
// if err != mongo.ErrNoDocuments {
// return d.handleErr(err)
// }
// }
//
// return nil
// }
//
// func (d *St) LogList(ctx context.Context, pars *entities.LogListParsSt) ([]*entities.LogSt, error) {
// 	collection := d.Db.Collection("log")
//
// 	ops := &options.FindOptions{}
//
// 	if pars != nil {
// 		if pars.Projection != nil {
// 			ops.Projection = *pars.Projection
// 		}
// 	}
//
// 	cur, err := collection.Find(ctx, bson.M{}, ops)
// 	if err != nil {
// 		return nil, d.handleErr(err)
// 	}
// 	defer cur.Close(ctx)
//
// 	result := make([]*entities.LogSt, 0)
//
// 	for cur.Next(ctx) {
// 		obj := &entities.LogSt{}
//
// 		err = cur.Decode(obj)
// 		if err != nil {
// 			return nil, d.handleErr(err)
// 		}
//
// 		result = append(result, obj)
// 	}
// 	if err = cur.Err(); err != nil {
// 		return nil, d.handleErr(err)
// 	}
//
// 	return result, nil
// }
//
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
