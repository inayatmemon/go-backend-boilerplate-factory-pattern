package brands_models

import "time"

type Brand struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"index,not null,column:name"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime,column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime,column:updated_at"`
}

type CreateBrandRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type EditBrandRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
	ID   string `json:"id" validate:"required,uuid"`
}
