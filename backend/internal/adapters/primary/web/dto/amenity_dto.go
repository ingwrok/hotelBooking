package dto

type AmenityRequest struct {
	Name string `json:"name"`
}

type AmenityResponse struct {
	AmenityID int    `json:"amenityId"`
	Name      string `json:"name"`
}
