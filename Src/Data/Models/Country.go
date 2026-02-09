package Models

type Country struct {
	BaseModel
	Name   string `gorm:"size:100;not null;"`
	Cities []City `gorm:"foreignKey:CountryID"`
}
