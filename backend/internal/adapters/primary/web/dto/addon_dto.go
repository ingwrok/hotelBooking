package dto

type AddonCategoryRequest struct {
	Name string `json:"name"`
}

type AddonCategoryResponse struct {
	CategoryID int    `json:"categoryId"`
	Name       string `json:"name"`
}

type AddonRequest struct {
	CategoryID  int     `json:"categoryId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	UnitName    string  `json:"unitName"`
	PictureURL  string  `json:"pictureUrl"`
}

type AddonResponse struct {
	AddonID     int     `json:"addonId"`
	CategoryID  int     `json:"categoryId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	UnitName    string  `json:"unitName"`
	PictureURL  string  `json:"pictureUrl"`
}
