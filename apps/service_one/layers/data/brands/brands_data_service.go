package brands_data_service

import (
	brands_models "go_boilerplate_project/models/brands"
	common_models "go_boilerplate_project/models/commons"
	mongodb_models "go_boilerplate_project/models/databases/mongodb"
	mysql_models "go_boilerplate_project/models/databases/mysql"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var brandCollection = "brands"
var brandsTable = "brands"

func (s *service) CreateBrandMySQL(brand *brands_models.Brand) error {
	if err := s.Helpers.MySQL.Create(&mysql_models.CreateInput{
		Value: brand,
	}); err != nil {
		return err
	}
	return nil
}

func (s *service) CreateBrandMongoDB(brand *brands_models.Brand) error {
	result, err := s.Helpers.MongoDB.InsertOne(&mongodb_models.InsertOneInput{
		Collection: brandCollection,
		Document:   brand,
	})
	if err != nil {
		return err
	}
	brand.ID = result.ID.(primitive.ObjectID).Hex()
	return nil
}

func (s *service) GetBrandMySQL(brand *brands_models.Brand) error {
	err := s.Helpers.MySQL.FindOne(&mysql_models.FindOneInput{
		Result: brand,
		Conditions: map[string]any{
			"id": brand.ID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetBrandMongoDB(brand *brands_models.Brand) error {
	err := s.Helpers.MongoDB.FindOne(&mongodb_models.FindOneInput{
		Result: brand,
		Filter: bson.M{
			"id": brand.ID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) EditBrandMySQL(brandName string, brandId string, ctx ...common_models.IsContextPresentInput) error {
	updateInput := &mysql_models.UpdateInput{
		Model:      &brands_models.Brand{},
		Conditions: map[string]any{
			"id": brandId,
		},
		Values: map[string]any{
			"name": brandName,
		},
	}

	if s.Helpers.Custom.IsContextPresent(ctx...) {
		updateInput.Context = ctx[0].Context
		updateInput.CancelFunc = ctx[0].CancelFunc
	}

	_, err := s.Helpers.MySQL.Update(updateInput)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) EditBrandMongoDB(brandName string, brandId string, ctx ...common_models.IsContextPresentInput) error {
	updateInput := &mongodb_models.UpdateOneInput{
		Collection: brandCollection,
		Filter: bson.M{
			"id": brandId,
		},
		Update: bson.M{
			"$set": bson.M{"name": brandName},
		},
	}

	if s.Helpers.Custom.IsContextPresent(ctx...) {
		updateInput.Context = ctx[0].Context
		updateInput.CancelFunc = ctx[0].CancelFunc
	}

	_, err := s.Helpers.MongoDB.UpdateOne(updateInput)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetCollectionName() string {
	return brandCollection
}

func (s *service) GetTableName() string {
	return brandsTable
}
