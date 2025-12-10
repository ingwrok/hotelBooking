package domain

import "time"

type Room struct {
	RoomID					int
	RoomTypeID		int
	RoomNumber		string
	Status 			string
}

type RoomType struct {
	RoomTypeID		int
	Name				string
	Description	string
	SizeSQM				float64
	BedType			string
	Capacity		int
	PictureURL		[]string
	AmenityIDs	 []int
}

type RoomTypeDetails struct {
	RoomTypeID		int
	Name				string
	Description	string
	SizeSQM				float64
	BedType			string
	Capacity		int
	PictureURL		[]string
	Amenities	 []string
}

type RoomBlock struct {
	RoomBlockID	int
	RoomID			int
	StartDate	time.Time
	EndDate		time.Time
	Reason			string
}