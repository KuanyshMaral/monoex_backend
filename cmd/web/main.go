package main

import (
	"log"
	"monoex_backend/internal/app"
)

func main() {
	// Создаем новое приложение
	myApp := app.New()

	// Инициализируем (конфиг, БД, репозитории, сервисы, маршруты)
	if err := myApp.Initialize(); err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Запускаем сервер
	if err := myApp.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
