// Package main содержит точку входа для клиентского приложения GophKeeper.
// Клиент реализует CLI интерфейс для управления приватными данными пользователя.
package main

import (
	"github.com/fylgushev/go-diplom-final/internal/interfaces/cli"
)

// main основная функция клиентского приложения.
func main() {
	cli.Execute()
}
