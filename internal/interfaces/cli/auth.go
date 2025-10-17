package cli

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewRegisterCmd создает команду регистрации
func NewRegisterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Регистрация нового пользователя",
		Long:  "Создает новую учетную запись в системе GophKeeper",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Получаем логин
			login, _ := cmd.Flags().GetString("login")
			if login == "" {
				fmt.Print("Введите логин: ")
				fmt.Scanln(&login)
			}

			// Получаем пароль
			password, _ := cmd.Flags().GetString("password")
			if password == "" {
				fmt.Print("Введите пароль: ")
				passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return fmt.Errorf("ошибка чтения пароля: %w", err)
				}
				password = string(passwordBytes)
				fmt.Println()
			}

			// TODO: вызов gRPC API регистрации
			fmt.Printf("Регистрация пользователя %s...\n", login)
			fmt.Println("✅ Пользователь успешно зарегистрирован!")
			
			return nil
		},
	}

	cmd.Flags().StringP("login", "l", "", "Логин пользователя")
	cmd.Flags().StringP("password", "p", "", "Пароль пользователя")

	return cmd
}

// NewLoginCmd создает команду входа
func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Вход в систему",
		Long:  "Аутентификация пользователя в системе GophKeeper",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Получаем логин
			login, _ := cmd.Flags().GetString("login")
			if login == "" {
				fmt.Print("Введите логин: ")
				fmt.Scanln(&login)
			}

			// Получаем пароль
			password, _ := cmd.Flags().GetString("password")
			if password == "" {
				fmt.Print("Введите пароль: ")
				passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return fmt.Errorf("ошибка чтения пароля: %w", err)
				}
				password = string(passwordBytes)
				fmt.Println()
			}

			// TODO: вызов gRPC API аутентификации
			fmt.Printf("Вход пользователя %s...\n", login)
			fmt.Println("✅ Успешный вход в систему!")
			
			return nil
		},
	}

	cmd.Flags().StringP("login", "l", "", "Логин пользователя")
	cmd.Flags().StringP("password", "p", "", "Пароль пользователя")

	return cmd
}

// NewLogoutCmd создает команду выхода
func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Выход из системы",
		Long:  "Завершение сессии и удаление локальных токенов",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: удаление токенов из конфигурации
			fmt.Println("✅ Вы успешно вышли из системы")
			return nil
		},
	}

	return cmd
}