package version

import (
	"runtime"
	"strings"
	"testing"
)

// TestGet проверяет функцию Get для получения полной информации о версии.
func TestGet(t *testing.T) {
	info := Get()

	// Проверяем, что все поля заполнены
	if info.Version == "" {
		t.Error("Version должна быть установлена")
	}
	if info.BuildDate == "" {
		t.Error("BuildDate должна быть установлена")
	}
	if info.GoVersion == "" {
		t.Error("GoVersion должна быть установлена")
	}
	if info.Platform == "" {
		t.Error("Platform должна быть установлена")
	}

	// Проверяем формат Platform
	expectedPlatform := runtime.GOOS + "/" + runtime.GOARCH
	if info.Platform != expectedPlatform {
		t.Errorf("Platform = %s, ожидалось %s", info.Platform, expectedPlatform)
	}

	// Проверяем, что GoVersion начинается с "go"
	if !strings.HasPrefix(info.GoVersion, "go") {
		t.Errorf("GoVersion = %s, должна начинаться с 'go'", info.GoVersion)
	}
}

// TestGetVersion проверяет функцию GetVersion.
func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("GetVersion() не должна возвращать пустую строку")
	}

	// В тестах по умолчанию должна быть "dev"
	if version != Version {
		t.Errorf("GetVersion() = %s, ожидалось %s", version, Version)
	}
}

// TestGetBuildDate проверяет функцию GetBuildDate.
func TestGetBuildDate(t *testing.T) {
	buildDate := GetBuildDate()
	if buildDate == "" {
		t.Error("GetBuildDate() не должна возвращать пустую строку")
	}

	if buildDate != BuildDate {
		t.Errorf("GetBuildDate() = %s, ожидалось %s", buildDate, BuildDate)
	}
}

// TestGetShortVersion проверяет функцию GetShortVersion.
func TestGetShortVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "версия с префиксом v",
			version:  "v1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "версия без префикса v",
			version:  "1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "пустая версия",
			version:  "",
			expected: "",
		},
		{
			name:     "только v",
			version:  "v",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Временно изменяем Version для теста
			oldVersion := Version
			Version = tt.version
			defer func() { Version = oldVersion }()

			result := GetShortVersion()
			if result != tt.expected {
				t.Errorf("GetShortVersion() = %s, ожидалось %s", result, tt.expected)
			}
		})
	}
}

// TestGetFullVersion проверяет функцию GetFullVersion.
func TestGetFullVersion(t *testing.T) {
	fullVersion := GetFullVersion()

	// Проверяем, что версия содержит и Version и BuildDate
	if !strings.Contains(fullVersion, Version) {
		t.Errorf("GetFullVersion() = %s, должна содержать %s", fullVersion, Version)
	}
	if !strings.Contains(fullVersion, BuildDate) {
		t.Errorf("GetFullVersion() = %s, должна содержать %s", fullVersion, BuildDate)
	}
	if !strings.Contains(fullVersion, "built:") {
		t.Errorf("GetFullVersion() = %s, должна содержать 'built:'", fullVersion)
	}
}

// TestInfoString проверяет метод String структуры Info.
func TestInfoString(t *testing.T) {
	info := Get()
	str := info.String()

	// Проверяем, что строка содержит все ключевые поля
	expectedFields := []string{"Version:", "Build Date:", "Go Version:", "Platform:"}
	for _, field := range expectedFields {
		if !strings.Contains(str, field) {
			t.Errorf("Info.String() = %s, должна содержать %s", str, field)
		}
	}
}

// TestIsDevBuild проверяет функцию IsDevBuild.
func TestIsDevBuild(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "dev версия",
			version:  "dev",
			expected: true,
		},
		{
			name:     "релизная версия",
			version:  "1.0.0",
			expected: false,
		},
		{
			name:     "пустая версия",
			version:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Временно изменяем Version для теста
			oldVersion := Version
			Version = tt.version
			defer func() { Version = oldVersion }()

			result := IsDevBuild()
			if result != tt.expected {
				t.Errorf("IsDevBuild() = %v, ожидалось %v", result, tt.expected)
			}
		})
	}
}

// TestHasValidBuildInfo проверяет функцию HasValidBuildInfo.
func TestHasValidBuildInfo(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		buildDate string
		expected  bool
	}{
		{
			name:      "валидная информация",
			version:   "1.0.0",
			buildDate: "2024-01-01T00:00:00Z",
			expected:  true,
		},
		{
			name:      "dev версия",
			version:   "dev",
			buildDate: "2024-01-01T00:00:00Z",
			expected:  false,
		},
		{
			name:      "неизвестная дата сборки",
			version:   "1.0.0",
			buildDate: "unknown",
			expected:  false,
		},
		{
			name:      "пустая версия",
			version:   "",
			buildDate: "2024-01-01T00:00:00Z",
			expected:  false,
		},
		{
			name:      "пустая дата",
			version:   "1.0.0",
			buildDate: "",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Временно изменяем переменные для теста
			oldVersion := Version
			oldBuildDate := BuildDate
			Version = tt.version
			BuildDate = tt.buildDate
			defer func() {
				Version = oldVersion
				BuildDate = oldBuildDate
			}()

			result := HasValidBuildInfo()
			if result != tt.expected {
				t.Errorf("HasValidBuildInfo() = %v, ожидалось %v", result, tt.expected)
			}
		})
	}
}

// TestFormatForCLI проверяет функцию FormatForCLI.
func TestFormatForCLI(t *testing.T) {
	// Тест для dev сборки
	t.Run("dev build", func(t *testing.T) {
		oldVersion := Version
		Version = "dev"
		defer func() { Version = oldVersion }()

		result := FormatForCLI()
		if !strings.Contains(result, "development build") {
			t.Errorf("FormatForCLI() для dev сборки должна содержать 'development build'")
		}
		if !strings.Contains(result, "GophKeeper") {
			t.Errorf("FormatForCLI() должна содержать 'GophKeeper'")
		}
	})

	// Тест для релизной сборки
	t.Run("release build", func(t *testing.T) {
		oldVersion := Version
		oldBuildDate := BuildDate
		Version = "1.0.0"
		BuildDate = "2024-01-01T00:00:00Z"
		defer func() {
			Version = oldVersion
			BuildDate = oldBuildDate
		}()

		result := FormatForCLI()
		if !strings.Contains(result, "Build Date:") {
			t.Errorf("FormatForCLI() для релизной сборки должна содержать 'Build Date:'")
		}
		if !strings.Contains(result, "GophKeeper") {
			t.Errorf("FormatForCLI() должна содержать 'GophKeeper'")
		}
		if !strings.Contains(result, Version) {
			t.Errorf("FormatForCLI() должна содержать версию %s", Version)
		}
		if !strings.Contains(result, BuildDate) {
			t.Errorf("FormatForCLI() должна содержать дату сборки %s", BuildDate)
		}
	})
}
