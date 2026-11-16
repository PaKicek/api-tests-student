package photos

import (
	"path"

	"api-tests-template/internal/utils"
)

// Имена тестовых фотографий, которые лежат в каталоге photosDir
const (
	FirstPhoto  = "1.jpg"
	SecondPhoto = "2.jpg"
	ThirdPhoto  = "3.jpg"
)

// photosDir каталог с тестовыми фотографиями относительно корня проекта
const photosDir = "testdata/photos"

// Path возвращает абсолютный путь к тестовой фотографии по её имени
func Path(photoFileName string) string {
	return utils.ProjectPath(path.Join(photosDir, photoFileName))
}
