package Dtos

type GetOtpRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,len=11"`
}

type TokenDetail struct {
	AccessToken            string `json:"accessToken"`
	RefreshToken           string `json:"refreshToken"`
	AccessTokenExpireTime  int64
	RefreshTokenExpireTime int64
}
