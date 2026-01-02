package dto

import "time"


type RoomRequest struct {
	RoomNumber		string `json:"room_number"`
}

type RoomResponse struct {
	RoomID					int `json:"room_id"`
	RoomTypeID		int `json:"room_type_id"`
	RoomTypeName	string `json:"room_type_name"`
	RoomNumber		string `json:"room_number"`
	Status 			string `json:"status"`
}

type RoomStatusRequest struct {
	Status string `json:"status"`
}

type RoomTypeRequest struct {
	Name				string `json:"name"`
	Description	string `json:"description"`
	SizeSQM				float64 `json:"size_sqm"`
	BedType			string `json:"bed_type"`
	Capacity		int `json:"capacity"`
	PictureURL		[]string `json:"picture_url"`
	AmenityIDs	 []int `json:"amenityIDs"`
}
type RoomTypeResponse struct {
	RoomTypeID		int `json:"room_type_id"`
	Name				string `json:"name"`
	Description	string `json:"description"`
	SizeSQM				float64 `json:"size_sqm"`
	BedType			string `json:"bed_type"`
	Capacity		int `json:"capacity"`
	PictureURL		[]string `json:"picture_url"`
}

type RoomTypeDetailResponse struct {
	RoomTypeID		int `json:"room_type_id"`
	Name				string `json:"name"`
	Description	string `json:"description"`
	SizeSQM				float64 `json:"size_sqm"`
	BedType			string `json:"bed_type"`
	Capacity		int `json:"capacity"`
	PictureURL		[]string `json:"picture_url"`
	Amenities	 []string `json:"amenities"`
}


type RoomBlockRequest struct {
	RoomID    int       `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Reason    string    `json:"reason"`
}

type RoomBlockResponse struct {
	RoomBlockID int       `json:"room_block_id"`
	RoomID      int       `json:"room_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Reason      string    `json:"reason"`
}


type AvailabilityRequest struct {
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
}

type FindRoomRequest struct {
	RoomTypeID int       `json:"room_type_id"`
	CheckIn    string `json:"check_in"`
	CheckOut   string `json:"check_out"`
}