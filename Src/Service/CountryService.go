package Service

import (
	"RideMarket-CleanWebApi-GoLang/Api/Dtos"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Constants"
	"RideMarket-CleanWebApi-GoLang/Data/Database/DatabaseConfig"
	"RideMarket-CleanWebApi-GoLang/Data/Models"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type CountryService struct {
	database *gorm.DB
	logger   Log.ILogger
	base     *BaseService[Models.Country, Dtos.CreateUpdateCountryRequest, Dtos.CreateUpdateCountryRequest, Dtos.CountryResponse]
}

func NewCountryService(cfg *Config.Config) *CountryService {
	return &CountryService{database: DatabaseConfig.GetDb(),
		base: &BaseService[Models.Country, Dtos.CreateUpdateCountryRequest, Dtos.CreateUpdateCountryRequest, Dtos.CountryResponse]{
			Database: DatabaseConfig.GetDb(),
			Logger:   Log.NewLogger(cfg),
		},
		logger: Log.NewLogger(cfg)}

}

func (s *CountryService) Create(ctx context.Context, req *Dtos.CreateUpdateCountryRequest) (*Dtos.CountryResponse, error) {
	country := Models.Country{Name: req.Name}
	country.CreatedBy = uint(ctx.Value(Constants.UserIdKey).(float64))
	country.CreatedAt = time.Now().UTC()
	tx := s.database.WithContext(ctx).Begin()
	err := tx.Create(&country).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(Log.Postgres, Log.Insert, err.Error(), nil)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(Log.Postgres, Log.Commit, err.Error(), nil)
		return nil, err
	}

	dto := &Dtos.CountryResponse{Name: country.Name}
	return dto, nil
}

func (s *CountryService) Update(ctx context.Context, id int, req *Dtos.CreateUpdateCountryRequest) (*Dtos.CountryResponse, error) {
	updateList := map[string]interface{}{
		"Name":      req.Name,
		"UpdatedBy": &sql.NullInt64{Int64: int64(ctx.Value(Constants.UserIdKey).(float64)), Valid: true},
		"UpdatedAt": &sql.NullTime{Time: time.Now().UTC(), Valid: true},
	}
	tx := s.database.WithContext(ctx).Begin()
	err := tx.Model(&Models.Country{}).Where("id=?", id).
		Updates(updateList).Error

	if err != nil {
		tx.Rollback()
		s.logger.Error(Log.Postgres, Log.Update, err.Error(), nil)
		return nil, err
	}

	var country Models.Country
	err = tx.Model(&Models.Country{}).Where("id=? and DeletedBy is null", id).
		First(&country).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(Log.Postgres, Log.Update, err.Error(), nil)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(Log.Postgres, Log.Commit, err.Error(), nil)
		return nil, err
	}

	var countryResponse Dtos.CountryResponse
	countryResponse.Name = country.Name
	countryResponse.Id = int(country.ID)

	return &countryResponse, nil

}

func (s *CountryService) Delete(ctx context.Context, id int) error {
	tx := s.database.WithContext(ctx).Begin()

	deleteMap := map[string]interface{}{
		"DeletedBy": &sql.NullInt64{Int64: int64(ctx.Value(Constants.UserIdKey).(float64)), Valid: true},
		"DeletedAt": &sql.NullTime{time.Now().UTC(), true},
	}
	err := tx.Model(&Models.Country{}).Where("id=?", id).
		Updates(deleteMap).Error

	if err != nil {
		tx.Rollback()
		s.logger.Error(Log.Postgres, Log.Delete, err.Error(), nil)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		s.logger.Error(Log.Postgres, Log.Commit, err.Error(), nil)
		return err
	}

	return nil
}

func (s *CountryService) GetById(ctx context.Context, id int) (*Dtos.CountryResponse, error) {
	var response Models.Country

	err := s.database.WithContext(ctx).Where("id=? and DeletedBy is null", id).First(response).Error
	if err != nil {
		s.logger.Error(Log.Postgres, Log.Select, err.Error(), nil)
		return nil, err
	}

	return &Dtos.CountryResponse{Id: int(response.ID), Name: response.Name}, nil

}
