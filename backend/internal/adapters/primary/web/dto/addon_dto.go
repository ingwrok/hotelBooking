package dto

type AddonCategoryRequest struct {
	Name       string `json:"name"`
}

type AddonCategoryResponse struct {
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
}

type AddonRequest struct {
	CategoryID int    `json:"category_id"`
	Name      string `json:"name"`
	Description string `json:"description"`
	Price     float64 `json:"price"`
	UnitName	string `json:"unit_name"`
}


type AddonResponse struct {
	AddonID   int    `json:"addon_id"`
	CategoryID int    `json:"category_id"`
	Name      string `json:"name"`
	Description string `json:"description"`
	Price     float64 `json:"price"`
	UnitName	string `json:"unit_name"`
}

