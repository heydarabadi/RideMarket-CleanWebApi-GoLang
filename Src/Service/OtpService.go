package Service

import (
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Constants"
	"RideMarket-CleanWebApi-GoLang/Data/Cache"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"RideMarket-CleanWebApi-GoLang/pkg/ServiceErrors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type OtpService struct {
	logger Log.ILogger
	cfg    *Config.Config
	redis  *redis.Client
}

type OtpDto struct {
	Value string
	Used  bool
}

func NewOtpService(cfg *Config.Config) *OtpService {
	logger := Log.NewLogger(cfg)
	redisClient := Cache.GetRedisInstance()

	return &OtpService{logger: logger, cfg: cfg, redis: redisClient}
}

func (service *OtpService) SetOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", Constants.RedisOtpDefaultKey, mobileNumber)
	val := &OtpDto{Value: otp, Used: false}

	res, err := Cache.Get[OtpDto](service.redis, key)
	if err == nil && !res.Used {
		return &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.OtpExists}
	} else if err == nil && res.Used {
		return &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.OtpUsed}
	}
	err = Cache.Set(service.redis, key, val, service.cfg.Otp.ExpireTime*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (service *OtpService) Validate(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", Constants.RedisOtpDefaultKey, mobileNumber)

	res, err := Cache.Get[OtpDto](service.redis, key)
	if err == nil && !res.Used {
		return &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.OtpExists}
	} else if err == nil && res.Used {
		return &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.OtpUsed}
	} else if !res.Used && res.Value == otp {
		res.Used = true
		err = Cache.Set(service.redis, key, res.Value, service.cfg.Otp.ExpireTime*time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}
