// Package main содержит точку входа для серверного приложения GophKeeper.
// Сервер реализует gRPC API для управления пользователями и их данными.
package main

import (
	"fmt"

	"github.com/rfruffer/go-diplom-final/pkg/version"
)

// main основная функция серверного приложения.
// todo Временная реализация для тестирования сборки и версионирования.
func main() {
	fmt.Println("GophKeeper Server")
	fmt.Println("==================")
	fmt.Println(version.FormatForCLI())
	fmt.Println("\nСервер в разработке...")
}
