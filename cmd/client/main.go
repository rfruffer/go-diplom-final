// Package main содержит точку входа для клиентского приложения GophKeeper.
// Клиент реализует CLI интерфейс для управления приватными данными пользователя.
package main

import (
	"fmt"

	"github.com/rfruffer/gophkeeper/pkg/version"
)

// main основная функция клиентского приложения.
// todo Временная реализация для тестирования сборки и версионирования.
func main() {
	fmt.Println("GophKeeper Client")
	fmt.Println("==================")
	fmt.Println(version.FormatForCLI())
	fmt.Println("\nКлиент в разработке...")
}
