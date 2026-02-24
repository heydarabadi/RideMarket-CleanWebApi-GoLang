package Dtos

type CreateUpdateCountryRequest struct {
	Name string `json:"name" binding:"required,alpha,min=6"`
}

type CountryResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
