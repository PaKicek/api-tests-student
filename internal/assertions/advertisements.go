package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	titles := make([]string, 0, len(items))
	for _, item := range items {
		if item.Title == title {
			return
		}
		titles = append(titles, item.Title)
	}

	require.Failf(t, "объявление не найдено",
		"Среди найденных объявлений нет объявления с названием %q. Найденные названия: %v", title, titles)
}
