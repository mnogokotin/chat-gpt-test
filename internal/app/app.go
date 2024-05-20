package app

import (
	"encoding/csv"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io"
	"log"
	"os"
	"sort"
)

type Event struct {
	ID          string `db:"id"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
}

func openOrCreateFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func createTable(db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		description TEXT,
		created_at TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func getEventsNotInCSV(db *sqlx.DB, idsInCSV map[string]struct{}) ([]Event, error) {
	var events []Event
	query := "SELECT id, description, created_at FROM events"
	err := db.Select(&events, query)
	if err != nil {
		return nil, err
	}

	var eventsNotInCSV []Event
	for _, event := range events {
		if _, exists := idsInCSV[event.ID]; !exists {
			eventsNotInCSV = append(eventsNotInCSV, event)
		}
	}

	return eventsNotInCSV, nil
}

func getCSVIDs(file *os.File) (map[string]struct{}, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.Comma = '|'
	ids := make(map[string]struct{})
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) > 0 {
			ids[record[0]] = struct{}{}
		}
	}

	return ids, nil
}

func writeEventsToCSV(file *os.File, events []Event) error {
	// Переместить курсор в конец файла
	_, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	writer.Comma = '|'
	for _, event := range events {
		record := []string{event.ID, event.Description, event.CreatedAt}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()

	return writer.Error()
}

func Main() {
	// Имя файла
	filename := "example.csv"

	// Подключение к базе данных PostgreSQL
	db, err := sqlx.Connect("postgres", "host=localhost port=54354 user=youruser password=yourpassword dbname=yourdb sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Создать таблицу events
	err = createTable(db)
	if err != nil {
		log.Fatalf("Не удалось создать таблицу: %v", err)
	}

	// Открыть или создать файл
	file, err := openOrCreateFile(filename)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
	defer file.Close()

	log.Printf("Файл %s успешно открыт или создан", filename)

	// Прочитать все ID из CSV файла
	idsInCSV, err := getCSVIDs(file)
	if err != nil {
		log.Fatalf("Ошибка при чтении ID из CSV: %v", err)
	}

	// Найти записи в базе данных, которых нет в CSV файле
	eventsNotInCSV, err := getEventsNotInCSV(db, idsInCSV)
	if err != nil {
		log.Fatalf("Ошибка при получении записей из базы данных: %v", err)
	}

	// Если есть записи, которых нет в CSV файле, добавить их в CSV
	if len(eventsNotInCSV) > 0 {
		log.Printf("Добавление %d записей в CSV файл", len(eventsNotInCSV))

		// Сортировать записи по возрастанию ID
		sort.Slice(eventsNotInCSV, func(i, j int) bool {
			return eventsNotInCSV[i].ID < eventsNotInCSV[j].ID
		})

		err = writeEventsToCSV(file, eventsNotInCSV)
		if err != nil {
			log.Fatalf("Ошибка при записи в CSV файл: %v", err)
		}

		log.Println("Записи успешно добавлены в CSV файл")
	} else {
		log.Println("Все записи из базы данных уже присутствуют в CSV файле")
	}
}
