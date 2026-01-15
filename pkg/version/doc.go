// Package version предоставляет информацию о версии приложения GophKeeper.
//
// Этот пакет содержит константы и функции для управления версией приложения,
// датой сборки и другой информацией о релизе.
//
// Версия и дата сборки устанавливаются во время компиляции через ldflags:
//
//	go build -ldflags "-X github.com/fylgushev/go-diplom-final/pkg/version.Version=1.0.0 \
//	                   -X github.com/fylgushev/go-diplom-final/pkg/version.BuildDate=$(date +%Y-%m-%d)"
//
// Пример использования:
//
//	fmt.Printf("Version: %s\n", version.Version)
//	fmt.Printf("Built at: %s\n", version.BuildDate)
//	version.PrintVersion()
package version
