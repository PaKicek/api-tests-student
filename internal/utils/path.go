package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const goModFileName = "go.mod"

// ProjectPath возвращает абсолютный путь к файлу или каталогу внутри проекта.
// Корень проекта (каталог с go.mod) ищется вверх от текущей рабочей директории теста
func ProjectPath(relativePath string) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Не удалось определить текущую директорию: %s", err.Error()))
	}

	for {
		if _, errStat := os.Stat(filepath.Join(dir, goModFileName)); errStat == nil {
			return filepath.Join(dir, relativePath)
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			panic(fmt.Sprintf("Не удалось найти корень проекта (файл %s) для пути %s", goModFileName, relativePath))
		}

		dir = parentDir
	}
}
