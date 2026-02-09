package Models

type City struct {
	BaseModel
	Name      string  `gorm:"size:255;not null;"`
	CountryID uint    `gorm:"not null;index"`
	Country   Country `gorm:"foreignKey:CountryID"`
}
