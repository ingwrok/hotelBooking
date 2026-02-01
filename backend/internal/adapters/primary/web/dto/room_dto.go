package dto

import "time"

type RoomRequest struct {
	RoomNumber string `json:"roomNumber"`
}

type RoomResponse struct {
	RoomID       int    `json:"roomId"`
	RoomTypeID   int    `json:"roomTypeId"`
	RoomTypeName string `json:"roomTypeName"`
	RoomNumber   string `json:"roomNumber"`
	Status       string `json:"status"`
}

type RoomStatusRequest struct {
	Status string `json:"status"`
}

type RoomTypeRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	SizeSQM     float64  `json:"sizeSqm"`
	BedType     string   `json:"bedType"`
	Capacity    int      `json:"capacity"`
	PictureURL  []string `json:"pictureUrl"`
	AmenityIDs  []int    `json:"amenityIds"`
}
type RoomTypeResponse struct {
	RoomTypeID  int      `json:"roomTypeId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	SizeSQM     float64  `json:"sizeSqm"`
	BedType     string   `json:"bedType"`
	Capacity    int      `json:"capacity"`
	PictureURL  []string `json:"pictureUrl"`
	TotalRooms  int      `json:"totalRooms"`
}

type RoomTypeDetailResponse struct {
	RoomTypeID  int      `json:"roomTypeId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	SizeSQM     float64  `json:"sizeSqm"`
	BedType     string   `json:"bedType"`
	Capacity    int      `json:"capacity"`
	PictureURL  []string `json:"pictureUrl"`
	Amenities   []string `json:"amenities"`
}

type RoomBlockRequest struct {
	RoomID    int    `json:"roomId"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Reason    string `json:"reason"`
}

type RoomBlockResponse struct {
	RoomBlockID int       `json:"roomBlockId"`
	RoomID      int       `json:"roomId"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Reason      string    `json:"reason"`
}

type AvailabilityRequest struct {
	CheckIn  string `json:"checkIn"`
	CheckOut string `json:"checkOut"`
}

type FindRoomRequest struct {
	RoomTypeID int    `json:"roomTypeId"`
	CheckIn    string `json:"checkIn"`
	CheckOut   string `json:"checkOut"`
}
