package path

const (
	AuthLoginPath        = "/auth/login"
	MyAdvertisementsPath = "/my/advertisements"
	AdvertisementsPath   = "/advertisements"
	AdvertisementPath    = "/advertisement"
	FavoritesPath        = "/favorites"
)

// AdvertisementPhotosPath шаблон пути для получения фотографий объявления
const AdvertisementPhotosPath = "/advertisements/%s/photos"

// FavoritePath шаблон пути для работы с объявлением в избранном пользователя
const FavoritePath = "/favorites/%s"
