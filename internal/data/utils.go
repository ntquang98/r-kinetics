package data

import "go.mongodb.org/mongo-driver/bson"

func ToOrderedBSON(obj any) (bson.D, error) {
	if obj == nil {
		return bson.D{}, nil
	}

	sel, err := bson.Marshal(obj)
	if err != nil {
		return bson.D{}, err
	}

	b := bson.D{}
	err = bson.Unmarshal(sel, &b)

	return b, err
}

func ToBSON(obj any) (bson.M, error) {
	if obj == nil {
		return bson.M{}, nil
	}

	sel, err := bson.Marshal(obj)
	if err != nil {
		return bson.M{}, err
	}

	b := bson.M{}
	err = bson.Unmarshal(sel, &b)

	return b, err
}
