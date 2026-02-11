package Dtos

type GetOtpRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,len=11"`
}
