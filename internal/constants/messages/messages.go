package messages

// Коды ошибок, которые API возвращает в поле error
const (
	CodeUnauthorized   = "unauthorized"
	CodeInvalidRequest = "invalid_request"
	CodeNotFound       = "not_found"
	CodeForbidden      = "forbidden"
	CodeAlreadyExists  = "already_exists"
)

// Сообщения об ошибках авторизации
const (
	AuthorizationHeaderRequired = "Authorization header required"
	InvalidSession              = "Invalid or expired token"
)

// Сообщения об ошибках валидации объявления
const (
	TitleRequired       = "title is required"
	DescriptionRequired = "description is required"
	PriceOutOfRange     = "price must be between 0 and 1000000 inclusive"
	QuantityOutOfRange  = "quantity must be between 1 and 100 inclusive"
	AtLeastOnePhoto     = "at least 1 photo is required"
	TooManyPhotos       = "maximum 3 photos allowed"
)

// TitleTooLongTemplate шаблон сообщения о превышении максимальной длины названия
const TitleTooLongTemplate = "title must not exceed %d characters provided %d"

// Сообщения об ошибках работы с избранным
const (
	AdvertisementNotFound           = "Advertisement not found"
	CannotAddOwnAdvertisement       = "Cannot add own advertisement to favorites"
	AdvertisementAlreadyInFavorites = "Advertisement is already in favorites"
	FavoriteNotFound                = "Favorite entry not found"
)
