package main // Текущий пакет

// Импортируем пакеты
import (
	"fmt"
	"os" // Пакет для работы с ОС из стандартной библиотеки Go
	// Пакет "searcher" нашего модуля rsslook
	"rsslook/searcher"
)

// Константа уровня пакета
const keyword = "sport"

// «Точка входа» в программу
func main() {
	key := keyword
	count := len(os.Args)
stop:
	for i := 0; i < count; i++ {
		switch {
		case count > 1 && i >= count-1:
			break stop
		case i != 0 || count > 1:
			key = os.Args[i+1]
		}
		getNews(key)
	}
}

func getNews(key string) {
	fmt.Printf("--- keyword for searching: %s ---\n", key)
	// Вызываем функцию "Search" из пакета "searcher"
	searcher.Search(feeds, key, os.Stdout)
	fmt.Println("")
}
