package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

// ContainsAdvertisementByTitle проверяет, что среди объявлений есть объявление с указанным названием
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	FindAdvertisementByTitle(t, items, title)
}

// NotContainsAdvertisementByTitle проверяет, что среди объявлений нет объявления с указанным названием
func NotContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	for _, item := range items {
		if item.Title == title {
			require.Failf(t, "объявление найдено",
				"Среди найденных объявлений есть объявление с названием %q, хотя его там быть не должно", title)
		}
	}
}

// FindAdvertisementByTitle возвращает объявление с указанным названием (если такого объявления нет - падает со списком найденных названий)
func FindAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) models.Advertisement {
	t.Helper()

	titles := make([]string, 0, len(items))
	for _, item := range items {
		if item.Title == title {
			return item
		}
		titles = append(titles, item.Title)
	}

	require.Failf(t, "объявление не найдено",
		"Среди найденных объявлений нет объявления с названием %q. Найденные названия: %v", title, titles)

	return models.Advertisement{}
}

// AdvertisementEquals проверяет совпадение ожидаемых и фактических параметров объявления
func AdvertisementEquals(t *testing.T, expected, actual models.Advertisement) {
	t.Helper()

	require.Equal(t, expected.Title, actual.Title, "Название объявления не совпадает")
	require.Equal(t, expected.Description, actual.Description, "Описание объявления не совпадает")
	require.Equal(t, expected.Price, actual.Price, "Цена объявления не совпадает")
	require.Equal(t, expected.Quantity, actual.Quantity, "Количество объявления не совпадает")
	require.Equal(t, expected.UserId, actual.UserId, "Владелец объявления не совпадает")
}

// PhotoIds возвращает идентификаторы фотографий объявления в порядке ответа API
func PhotoIds(photos []models.Photo) []string {
	ids := make([]string, 0, len(photos))
	for _, photo := range photos {
		ids = append(ids, photo.Id)
	}

	return ids
}
