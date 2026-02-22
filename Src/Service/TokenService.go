package Service

import (
	"RideMarket-CleanWebApi-GoLang/Api/Dtos"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"RideMarket-CleanWebApi-GoLang/pkg/ServiceErrors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	logger Log.ILogger
	cfg    *Config.Config
}

type tokenDto struct {
	UserId       int
	FullName     string
	UserName     string
	MobileNumber string
	Email        string
	Role         []string
}

func NewTokenService(cfg *Config.Config) *TokenService {
	logger := Log.NewLogger(cfg)
	return &TokenService{logger: logger, cfg: cfg}
}

func (s *TokenService) GenerateToken(token *tokenDto) (*Dtos.TokenDetail, error) {
	tokenDetail := &Dtos.TokenDetail{}

	expireTimeAccessToken := time.Now().Add(time.Duration(s.cfg.Jwt.AccessTokenExpireDuration) * time.Minute)
	tokenDetail.AccessTokenExpireTime = expireTimeAccessToken.Unix()

	expireTimeRefreshToken := time.Now().Add(time.Duration(s.cfg.Jwt.RefreshTokenExpireDuration) * time.Minute)
	tokenDetail.RefreshTokenExpireTime = expireTimeRefreshToken.Unix()

	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["user_id"] = token.UserId
	accessTokenClaims["fullname"] = token.FullName
	accessTokenClaims["username"] = token.UserName
	accessTokenClaims["email"] = token.Email
	accessTokenClaims["role"] = strings.Join(token.Role, ",")
	accessTokenClaims["exp"] = tokenDetail.AccessTokenExpireTime
	accessTokenClaims["mobile_number"] = token.MobileNumber

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	var err error
	tokenDetail.AccessToken, err = accessToken.SignedString([]byte(s.cfg.Jwt.Secret))

	if err != nil {
		return nil, err
	}

	refreshTokenClaim := jwt.MapClaims{}
	refreshTokenClaim["user_id"] = token.UserId
	refreshTokenClaim["exp"] = tokenDetail.RefreshTokenExpireTime
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaim)

	tokenDetail.RefreshToken, err = refreshToken.SignedString([]byte(s.cfg.Jwt.RefreshTokenSecret))
	if err != nil {
		return nil, err
	}

	return tokenDetail, nil

}

func (s *TokenService) VerifyToken(accessToken string) (*jwt.Token, error) {
	accessTokenVar, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.UnExpectedError}
		}
		return []byte(s.cfg.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return accessTokenVar, nil
}

func (s *TokenService) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}
	verifyToken, err := s.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for key, value := range claims {
			claimMap[key] = value
		}
		return claimMap, nil
	}
	return nil, &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.ClaimNotFound}
}
