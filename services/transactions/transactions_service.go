package transactions_service

import (
	mongodb_models "go_boilerplate_project/models/databases/mongodb"
	mysql_models "go_boilerplate_project/models/databases/mysql"
)

func (s *service) RunMongoDBTransaction(input *mongodb_models.TransactionInput) error {
	s.Input.Logger.Debugw("RunMongoDBTransaction started")
	err := s.Input.Helpers.MongoDB.RunTransaction(input)
	if err != nil {
		s.Input.Logger.Errorw("RunMongoDBTransaction failed", "error", err)
		return err
	}
	s.Input.Logger.Debugw("RunMongoDBTransaction completed successfully")
	return nil
}

func (s *service) RunMySQLTransaction(input *mysql_models.TransactionInput) error {
	s.Input.Logger.Debugw("RunMySQLTransaction started")
	err := s.Input.Helpers.MySQL.RunTransaction(input)
	if err != nil {
		s.Input.Logger.Errorw("RunMySQLTransaction failed", "error", err)
		return err
	}
	s.Input.Logger.Debugw("RunMySQLTransaction completed successfully")
	return nil
}
