// Package version предоставляет информацию о версии сборки GophKeeper.
// Этот пакет содержит функции для получения версии приложения и даты сборки,
// которые устанавливаются во время компиляции через ldflags.
package version

import (
	"fmt"
	"runtime"
)

// Переменные устанавливаются во время сборки через -ldflags
var (
	// Version версия приложения.
	// Устанавливается через -ldflags "-X github.com/rfruffer/gophkeeper/pkg/version.Version=1.0.0"
	Version = "dev"

	// BuildDate дата и время сборки.
	// Устанавливается через -ldflags "-X github.com/rfruffer/gophkeeper/pkg/version.BuildDate=2024-01-01T00:00:00Z"
	BuildDate = "unknown"
)

// Info содержит полную информацию о версии сборки.
// Используется для структурированного представления данных о версии
type Info struct {
	// Version версия приложения.
	Version string `json:"version"`

	// BuildDate дата и время сборки.
	BuildDate string `json:"build_date"`

	// GoVersion версия Go, использованная для сборки.
	GoVersion string `json:"go_version"`

	// Platform целевая платформа (OS/Arch).
	// Для поддержки множественных платформ (Windows, Linux, macOS).
	Platform string `json:"platform"`
}

// Get возвращает полную информацию о версии сборки.
// Включает версию, дату сборки, версию Go и платформу.
// Основная функция для получения всей информации о сборке.
func Get() Info {
	return Info{
		Version:   Version,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetVersion возвращает версию приложения.
func GetVersion() string {
	return Version
}

// GetBuildDate возвращает дату сборки приложения.
func GetBuildDate() string {
	return BuildDate
}

// GetShortVersion возвращает краткую версию (только номер версии).
// Убирает префикс "v" если он присутствует для более чистого отображения.
func GetShortVersion() string {
	if len(Version) > 0 && Version[0] == 'v' {
		return Version[1:]
	}
	return Version
}

// GetFullVersion возвращает полную версию с дополнительной информацией.
// Включает версию и дату сборки в удобочитаемом формате для CLI.
// Удобно для команды --version в CLI приложении.
func GetFullVersion() string {
	info := Get()
	return fmt.Sprintf("%s (built: %s)", info.Version, info.BuildDate)
}

// String возвращает строковое представление информации о версии.
// Реализует интерфейс fmt.Stringer для структуры Info.
// Используется для вывода полной информации в CLI.
func (i Info) String() string {
	return fmt.Sprintf("Version: %s\nBuild Date: %s\nGo Version: %s\nPlatform: %s",
		i.Version, i.BuildDate, i.GoVersion, i.Platform)
}

// IsDevBuild проверяет, является ли сборка разработческой.
// Возвращает true, если версия установлена как "dev".
// Для включения дополнительного логирования или функций отладки.
func IsDevBuild() bool {
	return Version == "dev"
}

// HasValidBuildInfo проверяет, что информация о сборке корректно установлена.
// Возвращает true, если все основные поля версии установлены.
// Помогает определить, была ли сборка выполнена правильно с ldflags.
func HasValidBuildInfo() bool {
	return Version != "dev" && Version != "" && BuildDate != "unknown" && BuildDate != ""
}

// FormatForCLI возвращает отформатированную строку для вывода в CLI.
// Выводит версию и дату сборки в понятном пользователю формате.
func FormatForCLI() string {
	if IsDevBuild() {
		return fmt.Sprintf("GophKeeper %s (development build)\nBuilt with %s for %s",
			Version, runtime.Version(), runtime.GOOS+"/"+runtime.GOARCH)
	}
	return fmt.Sprintf("GophKeeper %s\nBuild Date: %s\nBuilt with %s for %s",
		Version, BuildDate, runtime.Version(), runtime.GOOS+"/"+runtime.GOARCH)
}
