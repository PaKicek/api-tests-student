package models

type Advertisement struct {
	Id          string  `json:"id,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Quantity    int     `json:"quantity,omitempty"`
	UserId      string  `json:"user_id,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
	Photos      []Photo `json:"photos,omitempty"`
}

type Photo struct {
	Id              string `json:"id,omitempty"`
	AdvertisementId string `json:"advertisement_id,omitempty"`
	Url             string `json:"url,omitempty"`
	SortOrder       int    `json:"sort_order,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
}

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}
