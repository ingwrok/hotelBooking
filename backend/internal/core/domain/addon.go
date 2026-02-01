package domain

type AddonCategory struct {
	CategoryID int
	Name       string
}

type Addon struct {
	AddonID     int
	CategoryID  int
	Name        string
	Description string
	Price       float64
	UnitName    string
	PictureURL  string
}
