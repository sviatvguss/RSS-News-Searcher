package searcher

import (
	"encoding/xml" // Работа с XML
	"fmt"          // Работа с форматированным вводом-выводом
	"io"           // Работа с вводом-выводом
	"log"
	"net/http" // Работа с HTTP
	"sync"     // Работа с примитивами синхронизации
)

// Search выполняет конкурентный поиск новостей со словом keyword среди рассылок feeds.
// Результат записывается в out
func Search(feeds []string, keyword string, out io.Writer) {
	outCh := make(chan string)

	// Механизм синхронизации группы горутин
	wg := sync.WaitGroup{}
	wg.Add(len(feeds))

	// Обрабатываем каждый URL в отдельной горутине
	for _, f := range feeds {
		go func(url string) {
			defer wg.Done()
			outCh <- fmt.Sprintf("Process %s...\n", url)
			err := getFeedAndSearch(url, keyword, outCh)
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
			}
		}(f)
	}

	go func() {
		// Ждём, пока не будут обработаны все RSS-ленты
		wg.Wait()
		close(outCh)
	}()

	// Читаем из канала, пока его не закроют
	for msg := range outCh {
		_, _ = out.Write([]byte(msg))
	}
}

// getFeedAndSearch получает XML по url, парсит RSS ленту и ищет в ней новости с keyword
func getFeedAndSearch(url, keyword string, outCh chan<- string) error {
	// Выполняем HTTP GET-запрос
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Декодируем XML в структуру
	var feed rss
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return err
	}

	// Проходимся по новостям
	for _, item := range feed.Channel.Items {
		if item.HasKeyword(keyword) {
			outCh <- fmt.Sprintf("\n%s\n(%s)\n", item.Title, item.Link)
		}
	}
	return nil
}
