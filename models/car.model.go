package models

type Car struct {
	CarID       int        `json:"car_id"`
	UserID      int        `json:"user_id"`
	CarName     string     `json:"car_name" validate:"required"`
	Tags        string     `json:"tags" validate:"required"`
	Description string     `json:"description" validate:"required"`
	CarType     string     `json:"car_type" validate:"required"`
	CarCompany  string     `json:"car_company" validate:"required"`
	Dealer      string     `json:"dealer" validate:"required"`
	Images      []CarImage `json:"images"`
}
type CarImage struct {
	ImageID  uint32 `json:"image_id"`
	CarID    uint32 `json:"car_id"`
	ImageURL string `json:"image_url"`
}
