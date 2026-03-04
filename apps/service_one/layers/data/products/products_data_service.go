package products_data_service

import (
	common_models "go_boilerplate_project/models/commons"
	mongodb_models "go_boilerplate_project/models/databases/mongodb"
	mysql_models "go_boilerplate_project/models/databases/mysql"
	products_models "go_boilerplate_project/models/products"

	"go.mongodb.org/mongo-driver/bson"
)

var productCollection = "products"
var productsTable = "products"

func (s *service) CreateProductMySQL(product *products_models.Product) error {
	if err := s.Helpers.MySQL.Create(&mysql_models.CreateInput{
		Value: product,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) CreateProductMongoDB(product *products_models.Product) error {
	_, err := s.Helpers.MongoDB.InsertOne(&mongodb_models.InsertOneInput{
		Collection: productCollection,
		Document:   product,
	})
	if err != nil {
		return err
	}
	// product.ID = result.ID.(primitive.ObjectID)
	return nil
}

func (s *service) GetCollectionName() string {
	return productCollection
}

func (s *service) GetTableName() string {
	return productsTable
}

func (s *service) EditProductBrandNameMySQL(brandName string, productId string, ctx ...common_models.IsContextPresentInput) error {
	updateInput := &mysql_models.UpdateInput{
		Model: &products_models.Product{
			BrandName: brandName,
		},
		Conditions: map[string]any{
			"id": productId,
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

func (s *service) EditProductBrandNameMongoDB(brandName string, productId string, ctx ...common_models.IsContextPresentInput) error {
	updateInput := &mongodb_models.UpdateOneInput{
		Collection: productCollection,
		Filter: bson.M{
			"id": productId,
		},
		Update: bson.M{
			"brandName": brandName,
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
