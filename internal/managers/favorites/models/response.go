package models

import advModels "api-tests-template/internal/managers/advertisements/models"

type FavoriteItem struct {
	FavoriteId    string                  `json:"favorite_id,omitempty"`
	AddedAt       string                  `json:"added_at,omitempty"`
	Advertisement advModels.Advertisement `json:"advertisement"`
}

type FavoritesList struct {
	Items      []FavoriteItem `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}
