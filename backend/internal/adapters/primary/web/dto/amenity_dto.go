package dto

type AmenityRequest struct {
	Name        string `json:"name"`
}

type AmenityResponse struct {
	AmenityID	int    `json:"amenity_id"`
	Name        string `json:"name"`
}