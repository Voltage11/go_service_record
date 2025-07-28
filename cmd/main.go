package main

import (
	"fmt"
	"service-record/internal/config"
)

func main() {
	fmt.Println("Загрузка конфигурации...")
	cfg := config.GetConfig()
	fmt.Println("Конфигурация загружена")
	
}