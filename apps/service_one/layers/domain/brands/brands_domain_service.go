package brands_domain_service

import (
	api_models "go_boilerplate_project/models/api"
	common_models "go_boilerplate_project/models/commons"
	mysql_models "go_boilerplate_project/models/databases/mysql"
	"net/http"

	"gorm.io/gorm"
)

func (s *service) EditBrand(name string, id string) *api_models.ApiResponse {
	s.Input.Logger.Infow("EditBrand started", "name", name, "id", id)

	response := &api_models.ApiResponse{
		StatusCode: http.StatusOK,
		MessageKey: "brand_edit_success",
		Data:       nil,
	}

	ctx, cancel := s.Input.Services.Context.GetContext()

	err := s.Input.Services.Transactions.RunMySQLTransaction(&mysql_models.TransactionInput{
		Callback: func(tx *gorm.DB) error {
			err := s.Input.Data.Brands.EditBrandMySQL(name, id, common_models.IsContextPresentInput{
				Context:    ctx,
				CancelFunc: cancel,
			})
			if err != nil {
				return err
			}

			err = s.Input.Data.Products.EditProductBrandNameMySQL(name, id, common_models.IsContextPresentInput{
				Context:    ctx,
				CancelFunc: cancel,
			})
			if err != nil {
				return err
			}
			return nil
		},
	})

	if err != nil {
		s.Input.Logger.Errorw("EditBrand failed", "error", err, "name", name, "id", id)
		response.StatusCode = http.StatusInternalServerError
		response.MessageKey = "brand_edit_failed"
		response.ErrorKey = "error_detail"
		response.ErrorKeyParams = map[string]string{"detail": err.Error()}
		return response
	}

	s.Input.Logger.Infow("EditBrand completed successfully", "name", name, "id", id)
	return response
}
