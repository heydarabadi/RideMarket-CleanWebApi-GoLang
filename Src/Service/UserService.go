package Service

import (
	"RideMarket-CleanWebApi-GoLang/Api/Dtos"
	"RideMarket-CleanWebApi-GoLang/Common"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Data/Database/DatabaseConfig"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"

	"gorm.io/gorm"
)

type UserService struct {
	logger     Log.ILogger
	cfg        *Config.Config
	otpService *OtpService
	database   *gorm.DB
}

func NewUserService(cfg *Config.Config) *UserService {
	database := DatabaseConfig.GetDb()
	logger := Log.NewLogger(cfg)
	otpService := NewOtpService(cfg)
	return &UserService{logger: logger, database: database, cfg: cfg, otpService: otpService}

}

func (s *UserService) SendOtp(req *Dtos.GetOtpRequest) error {
	otp := Common.GenerateOtp()
	err := s.otpService.SetOtp(req.MobileNumber, otp)
	if err != nil {
		return err
	} else {
		return nil
	}
}
