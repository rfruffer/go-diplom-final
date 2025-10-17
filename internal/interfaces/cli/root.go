package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/fylgushev/go-diplom-final/pkg/version"
)

// NewRootCmd создает корневую команду для CLI
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gophkeeper",
		Short: "GophKeeper - безопасное хранение приватных данных",
		Long: `GophKeeper - это клиент-серверная система для хранения приватных данных пользователя:
- Пароли и логины
- Произвольные текстовые данные  
- Произвольные бинарные данные
- Данные банковских карт

Все данные хранятся в зашифрованном виде и синхронизируются между устройствами.`,
		Version: version.Version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Глобальные флаги
	cmd.PersistentFlags().StringP("server", "s", "localhost:8080", "Адрес gRPC сервера")
	cmd.PersistentFlags().StringP("config", "c", "", "Путь к файлу конфигурации")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Подробный вывод")

	// Добавляем подкоманды
	cmd.AddCommand(NewRegisterCmd())
	cmd.AddCommand(NewLoginCmd())
	cmd.AddCommand(NewLogoutCmd())
	cmd.AddCommand(NewAddCmd())
	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewSyncCmd())

	return cmd
}

// Execute выполняет корневую команду
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}