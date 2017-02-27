package entity

type (
	//Brand shows a Phone brand
	BrandType string
)

const (
	BrandTypeApple    BrandType = "Apple"
	BrandTypeSamsung  BrandType = "Samsung"
	BrandTypeLG       BrandType = "LG"
	BrandTypeMotorlla BrandType = "Motorlla"
)

// OS is the os
type OS struct {
	Valid  bool
	ID     int64
	Name   string
	Mobile bool
}
