package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// OpenOrCreateFile открывает файл для чтения и записи или создает его, если он не существует,
// и возвращает ссылку на файл и ошибку, если она возникла.
func OpenOrCreateFile(filename string) (*os.File, error) {
	// Открыть файл или создать его, если он не существует
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// ReadLastLine читает последнюю строку из файла
func ReadLastLine(file *os.File) (string, error) {
	var lastLine string

	// Переместить курсор в начало файла
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if lastLine == "" {
		return "", fmt.Errorf("файл пуст или строк нет")
	}

	return lastLine, nil
}

//func OpenFileOnWriteAtTheEnd(outputFilePath string) (*os.File, error) {
//	f, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
//	if err != nil {
//		return nil, err
//	}
//	return f, nil
//}
//
//func WriteModelsToFile(file *os.File, csvSeparator string, models interface{}) error {
//	w := csv.NewWriter(file)
//	w.Comma = []rune(csvSeparator)[0]
//	defer w.Flush()
//
//	v := reflect.ValueOf(models)
//	for i := 0; i < v.Len(); i++ {
//		line := scsv.GetStringSliceFromModel(v.Index(i).Interface())
//		if err := w.Write(line); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
