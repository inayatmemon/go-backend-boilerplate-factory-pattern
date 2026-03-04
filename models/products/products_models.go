package products_models

import "time"

type Product struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"index,not null,column:name"`
	BrandId   string    `json:"brandId" gorm:"index,not null,column:brand_id"`
	BrandName string    `json:"brandName" gorm:"index,not null,column:brand_name"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime,column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime,column:updated_at"`
}

type CreateProductRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=100"`
	BrandId string `json:"brandId"`
}
