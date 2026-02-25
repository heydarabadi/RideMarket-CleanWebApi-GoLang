package Dtos

type CreateUpdateCountryRequest struct {
	Name string `json:"name" binding:"required,alpha,min=6"`
}

type CountryResponse struct {
	Id     int            `json:"id"`
	Name   string         `json:"name"`
	Cities []CityResponse `json:"cities"`
}

type CreateUpdateCityRequest struct {
	Name string `json:"name" binding:"required,alpha,min=6"`
}

type CityResponse struct {
	Id      int             `json:"id"`
	Name    string          `json:"name"`
	Country CountryResponse `json:"country"`
}
