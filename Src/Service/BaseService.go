package Service

import (
	"RideMarket-CleanWebApi-GoLang/Common"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Constants"
	"RideMarket-CleanWebApi-GoLang/Data/Database/DatabaseConfig"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"RideMarket-CleanWebApi-GoLang/pkg/ServiceErrors"
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type BaseService[T any, Tc any, Tu any, Tr any] struct {
	Database *gorm.DB
	Logger   Log.ILogger
}

func NewBaseService[T any, Tc any, Tu any, Tr any](cfg *Config.Config) *BaseService[T, Tc, Tu, Tr] {
	return &BaseService[T, Tc, Tu, Tr]{
		Database: DatabaseConfig.GetDb(),
		Logger:   Log.NewLogger(cfg),
	}
}

func (s *BaseService[T, Tc, Tu, Tr]) Create(ctx context.Context, req *Tc) (*Tr, error) {
	model, err := Common.TypeConverter[T](req)
	tx := s.Database.WithContext(ctx).Begin()
	err = tx.WithContext(ctx).Create(model).Error
	if err != nil {
		tx.WithContext(ctx).Rollback()
		s.Logger.Error(Log.Postgres, Log.Insert, err.Error(), nil)
	}
	err = tx.WithContext(ctx).Commit().Error
	if err != nil {
		tx.WithContext(ctx).Rollback()
		s.Logger.Error(Log.Postgres, Log.Commit, err.Error(), nil)
	}
	return Common.TypeConverter[Tr](model)

}

func (s *BaseService[T, Tc, Tu, Tr]) Update(ctx context.Context, id int, req *Tc) (*Tr, error) {
	updateMap, _ := Common.TypeConverter[map[string]interface{}](req)
	(*updateMap)["UpdatedBy"] = &sql.NullInt64{Int64: int64(ctx.Value(Constants.UserIdKey).(float64)), Valid: true}
	(*updateMap)["UpdatedAt"] = &sql.NullTime{Time: time.Now().UTC(), Valid: true}

	model := new(T)
	tx := s.Database.WithContext(ctx).Begin()

	if err := tx.Model(model).
		Where("id = ? and DeletedBy is null", id).
		Updates(*updateMap).Error; err != nil {
		tx.Rollback()
		s.Logger.Error(Log.Postgres, Log.Update, err.Error(), nil)
		return nil, err
	}

	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		tx.Rollback()
		s.Logger.Error(Log.Postgres, Log.Commit, err.Error(), nil)
		return nil, err
	}

	result, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *BaseService[T, Tc, Tu, Tr]) Delete(ctx context.Context, id int) error {
	tx := s.Database.WithContext(ctx).Begin()
	model := new(T)
	userId := tx.Statement.Context.Value(Constants.UserIdKey).(float64)

	deleteMap := map[string]interface{}{
		"DeletedBy": &sql.NullInt64{Int64: int64(userId), Valid: true},
		"DeletedAt": &sql.NullTime{time.Now().UTC(), true},
	}

	countAffected := tx.Model(model).
		Where("id = ? and DeletedBy is null", id).
		Updates(deleteMap).
		RowsAffected

	if countAffected == 0 {
		s.Logger.Error(Log.Postgres, Log.Update, ServiceErrors.RecordNotFound, nil)
		tx.Rollback()
		return &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.RecordNotFound}
	}
	tx.Commit()
	return nil
}

func (s *BaseService[T, Tc, Tu, Tr]) GetById(ctx context.Context, id int) (*Tr, error) {
	model := new(T)
	err := s.Database.
		Where("id = ? and DeletedBy is null", id).
		First(model).Error

	if err != nil {
		s.Logger.Error(Log.Postgres, Log.Select, err.Error(), nil)
		return nil, err
	}
	res, err := Common.TypeConverter[Tr](model)
	if err != nil {
		return nil, err
	}
	return res, nil
}
