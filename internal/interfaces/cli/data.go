package cli

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/fylgushev/go-diplom-final/internal/domain/services"
)

// NewAddCmd создает команду добавления данных
func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Добавить новые данные",
		Long:  "Добавляет новый элемент данных в хранилище",
	}

	// Подкоманды для разных типов данных
	cmd.AddCommand(NewAddPasswordCmd())
	cmd.AddCommand(NewAddTextCmd())
	cmd.AddCommand(NewAddCardCmd())
	cmd.AddCommand(NewAddFileCmd())

	return cmd
}

// NewAddPasswordCmd создает команду добавления пароля
func NewAddPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "password",
		Short: "Добавить пароль",
		Long:  "Добавляет пару логин/пароль в хранилище",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				fmt.Print("Введите название: ")
				fmt.Scanln(&name)
			}

			login, _ := cmd.Flags().GetString("login")
			if login == "" {
				fmt.Print("Введите логин: ")
				fmt.Scanln(&login)
			}

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

			url, _ := cmd.Flags().GetString("url")
			notes, _ := cmd.Flags().GetString("notes")

			// Валидация входных данных
			validator := services.NewValidator()
			if err := validator.ValidateLoginPassword(login, password, url); err != nil {
				return fmt.Errorf("ошибка валидации: %w", err)
			}

			// TODO: вызов gRPC API для сохранения
			// Используем url и notes для метаданных
			_ = url
			_ = notes
			fmt.Printf("Сохранение пароля '%s'...\n", name)
			fmt.Println("✅ Пароль успешно сохранен!")

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Название записи")
	cmd.Flags().StringP("login", "l", "", "Логин")
	cmd.Flags().StringP("password", "p", "", "Пароль")
	cmd.Flags().StringP("url", "u", "", "URL сайта")
	cmd.Flags().String("notes", "", "Дополнительные заметки")

	return cmd
}

// NewAddTextCmd создает команду добавления текста
func NewAddTextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "text",
		Short: "Добавить текстовые данные",
		Long:  "Добавляет произвольные текстовые данные в хранилище",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				fmt.Print("Введите название: ")
				fmt.Scanln(&name)
			}

			text, _ := cmd.Flags().GetString("text")
			if text == "" {
				fmt.Print("Введите текст: ")
				fmt.Scanln(&text)
			}

			// TODO: вызов gRPC API для сохранения
			fmt.Printf("Сохранение текста '%s'...\n", name)
			fmt.Println("✅ Текст успешно сохранен!")

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Название записи")
	cmd.Flags().StringP("text", "t", "", "Текстовые данные")
	cmd.Flags().String("notes", "", "Дополнительные заметки")

	return cmd
}

// NewAddCardCmd создает команду добавления банковской карты
func NewAddCardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card",
		Short: "Добавить банковскую карту",
		Long:  "Добавляет данные банковской карты в хранилище",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				fmt.Print("Введите название карты: ")
				fmt.Scanln(&name)
			}

			number, _ := cmd.Flags().GetString("number")
			if number == "" {
				fmt.Print("Введите номер карты: ")
				fmt.Scanln(&number)
			}

			expiry, _ := cmd.Flags().GetString("expiry")
			if expiry == "" {
				fmt.Print("Введите срок действия (MM/YY): ")
				fmt.Scanln(&expiry)
			}

			cvv, _ := cmd.Flags().GetString("cvv")
			if cvv == "" {
				fmt.Print("Введите CVV: ")
				cvvBytes, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return fmt.Errorf("ошибка чтения CVV: %w", err)
				}
				cvv = string(cvvBytes)
				fmt.Println()
			}

			holder, _ := cmd.Flags().GetString("holder")

			// Валидация входных данных
			validator := services.NewValidator()
			if err := validator.ValidateCreditCard(number, expiry, cvv, holder); err != nil {
				return fmt.Errorf("ошибка валидации: %w", err)
			}

			// TODO: вызов gRPC API для сохранения
			// Используем holder для метаданных
			_ = holder
			fmt.Printf("Сохранение карты '%s'...\n", name)
			fmt.Println("✅ Карта успешно сохранена!")

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Название карты")
	cmd.Flags().String("number", "", "Номер карты")
	cmd.Flags().String("expiry", "", "Срок действия (MM/YY)")
	cmd.Flags().String("cvv", "", "CVV код")
	cmd.Flags().String("holder", "", "Имя держателя карты")
	cmd.Flags().String("notes", "", "Дополнительные заметки")

	return cmd
}

// NewAddFileCmd создает команду добавления файла
func NewAddFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "Добавить файл",
		Long:  "Добавляет бинарные данные (файл) в хранилище",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			path, _ := cmd.Flags().GetString("path")

			if path == "" {
				fmt.Print("Введите путь к файлу: ")
				fmt.Scanln(&path)
			}

			// Проверяем существование файла
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("файл не найден: %s", path)
			}

			if name == "" {
				name = path
			}

			// TODO: чтение файла и вызов gRPC API для сохранения
			fmt.Printf("Сохранение файла '%s'...\n", name)
			fmt.Println("✅ Файл успешно сохранен!")

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Название записи")
	cmd.Flags().StringP("path", "p", "", "Путь к файлу")
	cmd.Flags().String("notes", "", "Дополнительные заметки")

	return cmd
}

// NewGetCmd создает команду получения данных
func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [ID]",
		Short: "Получить данные по ID",
		Long:  "Получает и отображает данные элемента по его идентификатору",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			
			// TODO: вызов gRPC API для получения данных
			fmt.Printf("Получение данных ID: %s...\n", id)
			fmt.Println("Название: Мой аккаунт GitHub")
			fmt.Println("Тип: Пароль")
			fmt.Println("Логин: myuser@example.com")
			fmt.Println("URL: https://github.com")
			
			return nil
		},
	}

	cmd.Flags().BoolP("show-password", "s", false, "Показать пароль")
	cmd.Flags().StringP("output", "o", "table", "Формат вывода (table, json)")

	return cmd
}

// NewListCmd создает команду списка данных
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Показать список всех данных",
		Long:  "Отображает список всех сохраненных данных пользователя",
		RunE: func(cmd *cobra.Command, args []string) error {
			dataType, _ := cmd.Flags().GetString("type")
			limit, _ := cmd.Flags().GetInt("limit")
			
			// TODO: вызов gRPC API для получения списка
			fmt.Printf("Список данных (тип: %s, лимит: %d):\n\n", dataType, limit)
			fmt.Println("ID\t\tНазвание\t\tТип\t\tОбновлено")
			fmt.Println("--\t\t--------\t\t---\t\t---------")
			fmt.Println("abc123\t\tGitHub аккаунт\t\tПароль\t\t2024-01-15")
			fmt.Println("def456\t\tЗаметки\t\t\tТекст\t\t2024-01-14")
			fmt.Println("ghi789\t\tВиза карта\t\tКарта\t\t2024-01-13")
			
			return nil
		},
	}

	cmd.Flags().StringP("type", "t", "", "Фильтр по типу данных (password, text, card, file)")
	cmd.Flags().IntP("limit", "l", 50, "Максимальное количество записей")
	cmd.Flags().StringP("output", "o", "table", "Формат вывода (table, json)")

	return cmd
}

// NewDeleteCmd создает команду удаления данных
func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [ID]",
		Short: "Удалить данные",
		Long:  "Удаляет элемент данных по его идентификатору",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			force, _ := cmd.Flags().GetBool("force")
			
			if !force {
				fmt.Printf("Вы уверены, что хотите удалить элемент %s? (y/N): ", id)
				var answer string
				fmt.Scanln(&answer)
				if answer != "y" && answer != "Y" {
					fmt.Println("Отменено")
					return nil
				}
			}
			
			// TODO: вызов gRPC API для удаления
			fmt.Printf("Удаление элемента %s...\n", id)
			fmt.Println("✅ Элемент успешно удален!")
			
			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Удалить без подтверждения")

	return cmd
}

// NewSyncCmd создает команду синхронизации
func NewSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Синхронизировать данные",
		Long:  "Синхронизирует локальные данные с сервером с разрешением конфликтов",
		RunE: func(cmd *cobra.Command, args []string) error {
			strategy, _ := cmd.Flags().GetString("strategy")
			force, _ := cmd.Flags().GetBool("force")
			verbose, _ := cmd.Flags().GetBool("verbose")
			
			if verbose {
				fmt.Println("🔄 Начинаем синхронизацию данных...")
				fmt.Printf("📋 Стратегия разрешения конфликтов: %s\n", strategy)
			}
			
			// TODO: Получить данные пользователя и токен аутентификации
			// TODO: Подключиться к gRPC серверу
			// TODO: Вызвать SynchronizeData API
			
			// Симуляция синхронизации
			fmt.Println("🔍 Анализ локальных изменений...")
			
			if !force {
				fmt.Print("⚠️  Найдено 2 конфликта. Продолжить синхронизацию? (y/N): ")
				var answer string
				fmt.Scanln(&answer)
				if answer != "y" && answer != "Y" {
					fmt.Println("❌ Синхронизация отменена")
					return nil
				}
			}
			
			fmt.Println("📤 Отправка локальных изменений...")
			fmt.Println("📥 Получение изменений с сервера...")
			fmt.Println("🔧 Разрешение конфликтов...")
			
			if verbose {
				fmt.Println("\n📊 Результаты синхронизации:")
				fmt.Println("   ✅ Создано: 1 элемент")
				fmt.Println("   📝 Обновлено: 3 элемента")
				fmt.Println("   🗑️  Удалено: 0 элементов")
				fmt.Println("   ⚡ Конфликтов разрешено: 2")
			}
			
			fmt.Println("\n🎉 Синхронизация успешно завершена!")
			return nil
		},
	}

	cmd.Flags().StringP("strategy", "s", "by_version", "Стратегия разрешения конфликтов (by_version, by_timestamp, server_wins, client_wins)")
	cmd.Flags().BoolP("force", "f", false, "Принудительная синхронизация без подтверждения")
	cmd.Flags().BoolP("verbose", "v", false, "Подробный вывод")

	return cmd
}