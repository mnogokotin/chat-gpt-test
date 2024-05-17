package file

import (
	"encoding/csv"
	scsv "github.com/mnogokotin/chat-gpt-test/internal/service/csv"
	"os"
	"reflect"
)

func OpenOrCreateFileOnRead(outputFilePath string) (*os.File, error) {
	f, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func OpenFileOnWriteAtTheEnd(outputFilePath string) (*os.File, error) {
	f, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func WriteModelsToFile(file *os.File, csvSeparator string, models interface{}) error {
	w := csv.NewWriter(file)
	w.Comma = []rune(csvSeparator)[0]
	defer w.Flush()

	v := reflect.ValueOf(models)
	for i := 0; i < v.Len(); i++ {
		line := scsv.GetStringSliceFromModel(v.Index(i).Interface())
		if err := w.Write(line); err != nil {
			return err
		}
	}

	return nil
}
