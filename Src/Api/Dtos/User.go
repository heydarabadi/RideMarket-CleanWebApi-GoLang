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

type RegisterUserByUsernameRequest struct {
	FullName string `json:"fullName" binding:"required,min=6"`
	UserName string `json:"userName" binding:"required,min=5"`
	Email    string `json:"email" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=6,password"`
}

type RegisterLoginByMobileRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,len=11"`
	Otp          string `json:"otp" binding:"required,min=6,max=6"`
}

type LoginByUsernameRequest struct {
	Username string `json:"username" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=6"`
}
