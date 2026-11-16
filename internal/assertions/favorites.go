package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites/models"
)

// ContainsFavoriteAdvertisement проверяет, что в избранном есть объявление с указанным идентификатором
func ContainsFavoriteAdvertisement(t *testing.T, items []models.FavoriteItem, advertisementId string) models.FavoriteItem {
	t.Helper()

	for _, item := range items {
		if item.Advertisement.Id == advertisementId {
			return item
		}
	}

	require.Failf(t, "объявление не найдено в избранном",
		"В избранном нет объявления с идентификатором %s. Идентификаторы найденных объявлений: %v",
		advertisementId, favoriteAdvertisementIds(items))

	return models.FavoriteItem{}
}

// NotContainsFavoriteAdvertisement проверяет, что в избранном нет объявления с указанным идентификатором
func NotContainsFavoriteAdvertisement(t *testing.T, items []models.FavoriteItem, advertisementId string) {
	t.Helper()

	for _, item := range items {
		if item.Advertisement.Id == advertisementId {
			require.Failf(t, "объявление осталось в избранном",
				"В избранном всё ещё есть объявление с идентификатором %s и названием %q",
				advertisementId, item.Advertisement.Title)
		}
	}
}

// favoriteAdvertisementIds возвращает идентификаторы объявлений из списка избранного
func favoriteAdvertisementIds(items []models.FavoriteItem) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.Advertisement.Id)
	}

	return ids
}
