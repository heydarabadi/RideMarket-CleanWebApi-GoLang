package Service

import (
	"RideMarket-CleanWebApi-GoLang/Api/Dtos"
	"RideMarket-CleanWebApi-GoLang/Common"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Constants"
	"RideMarket-CleanWebApi-GoLang/Data/Database/DatabaseConfig"
	"RideMarket-CleanWebApi-GoLang/Data/Models"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"RideMarket-CleanWebApi-GoLang/pkg/ServiceErrors"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	logger       Log.ILogger
	cfg          *Config.Config
	otpService   *OtpService
	database     *gorm.DB
	tokenService *TokenService
}

const countQuery = "count(*) > 0"

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
	}

	return nil
}

// پیدا کردن یا آماده‌سازی کاربر جدید
func (s *UserService) findOrInitializeUser(mobile string) (*Models.User, bool, error) {
	var user Models.User
	err := s.database.
		Where("phone = ?", mobile).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		First(&user).Error

	if err == nil {
		return &user, false, nil // وجود دارد → لاگین
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	password, err := Common.GeneratePassword(10, true)
	if err != nil {
		return nil, false, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(Log.General, Log.HashPassword, err.Error(), nil)
		return nil, false, err
	}

	newUser := &Models.User{
		Phone:    mobile,
		Username: mobile,
		Password: string(hashed),
	}

	return newUser, true, nil
}

func (s *UserService) createUserWithDefaultRole(user *Models.User) error {
	roleId, err := s.getDefaultRole()
	if err != nil {
		s.logger.Error(Log.Postgres, Log.DefaultRoleNotFound, err.Error(), nil)
		return err
	}

	tx := s.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		s.logger.Error(Log.Postgres, Log.Rolback, err.Error(), nil)
		return err
	}

	userRole := Models.UserRole{
		UserId: int(user.ID),
		RoleId: roleId,
	}

	if err := tx.Create(&userRole).Error; err != nil {
		tx.Rollback()
		s.logger.Error(Log.Postgres, Log.Rolback, err.Error(), nil)
		return err
	}

	return tx.Commit().Error
}

func (s *UserService) buildTokenData(user *Models.User) *tokenDto {
	dto := &tokenDto{
		UserId:       int(user.ID),
		FullName:     user.FullName,
		Email:        user.Email,
		MobileNumber: user.Phone,
	}

	if user.UserRoles != nil {
		for _, ur := range *user.UserRoles {
			if ur.Role != nil {
				dto.Role = append(dto.Role, ur.Role.Name)
			}
		}
	}

	return dto
}

func (s *UserService) existsByUserName(username string) (bool, error) {
	var exists bool
	if err := s.database.Model(&Models.User{}).
		Select(countQuery).
		Where("UserName = ?", username).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(Log.Postgres, Log.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) existsByEmail(email string) (bool, error) {
	var exists bool
	if err := s.database.Model(&Models.User{}).
		Select(countQuery).
		Where("Email = ?", email).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(Log.Postgres, Log.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) existsByMobileNumber(mobileNumber string) (bool, error) {
	var exists bool
	if err := s.database.Model(&Models.User{}).
		Select(countQuery).
		Where("Phone = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(Log.Postgres, Log.Select, err.Error(), nil)
		return false, nil
	}

	return exists, nil
}

func (s *UserService) getDefaultRole() (roleId int, err error) {
	if err = s.database.Model(&Models.Role{}).
		Select("Id").
		Where("name = ?", Constants.DefaultRoleName).
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}

func (s *UserService) RegisterByUserName(req *Dtos.RegisterUserByUsernameRequest) error {
	user := Models.User{Username: req.UserName, FullName: req.FullName, Email: req.Email}

	exists, err := s.existsByEmail(req.Email)
	if err != nil {
		return err
	}
	if !exists {
		return &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.EmailExists}
	}

	exists, err = s.existsByUserName(req.UserName)
	if err != nil {
		return err
	}
	if !exists {
		return &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.UserNameExists}
	}

	bytePassword := []byte(req.Password)
	hashPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(Log.General, Log.HashPassword, err.Error(), nil)
		return err
	}
	user.Password = string(hashPassword)

	roleId, err := s.getDefaultRole()
	if err != nil {
		s.logger.Error(Log.Postgres, Log.DefaultRoleNotFound, err.Error(), nil)
	}

	transaction := s.database.Begin()
	err = transaction.Create(&user).Error

	if err != nil {
		transaction.Rollback()
		s.logger.Error(Log.Postgres, Log.Rolback, err.Error(), nil)
		return err
	}

	err = transaction.Create(&Models.UserRole{RoleId: roleId, UserId: int(user.ID)}).Error
	if err != nil {
		transaction.Rollback()
		s.logger.Error(Log.Postgres, Log.Rolback, err.Error(), nil)
		return err
	}

	err = transaction.Commit().Error
	if err != nil {
		transaction.Rollback()
		s.logger.Error(Log.Postgres, Log.Rolback, err.Error(), nil)
	}
	return nil
}
func (s *UserService) RegisterLoginByMobileNumber(req *Dtos.RegisterLoginByMobileRequest) (*Dtos.TokenDetail, error) {
	if err := s.otpService.Validate(req.MobileNumber, req.Otp); err != nil {
		return nil, err
	}

	user, isNewUser, err := s.findOrInitializeUser(req.MobileNumber)
	if err != nil {
		return nil, err
	}

	if isNewUser {
		if err := s.createUserWithDefaultRole(user); err != nil {
			return nil, err
		}
	}

	tokenData := s.buildTokenData(user)
	token, err := s.tokenService.GenerateToken(tokenData)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *UserService) LoginByUserName(req *Dtos.LoginByUsernameRequest) (*Dtos.TokenDetail, error) {
	var user Models.User

	err := s.database.
		Model(&Models.User{}).
		Where("Username = ?", req.Username).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&user).Error

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}
	tdo := tokenDto{UserId: int(user.ID), FullName: user.FullName, Email: user.Email,
		MobileNumber: user.Phone}

	if len(*user.UserRoles) > 0 {
		for _, ur := range *user.UserRoles {
			tdo.Role = append(tdo.Role, ur.Role.Name)
		}
	}

	token, err := s.tokenService.GenerateToken(&tdo)
	if err != nil {
		return nil, err
	}
	return token, nil

}
