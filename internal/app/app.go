package app

import (
	"context"
	"errors"
	"github.com/kardianos/service"
	"github.com/mnogokotin/chat-gpt-test/internal/config"
	"github.com/mnogokotin/chat-gpt-test/internal/service/csv"
	"github.com/mnogokotin/chat-gpt-test/internal/service/database"
	"github.com/mnogokotin/chat-gpt-test/internal/service/file"
	"github.com/mnogokotin/golang-packages/database/postgres"
	ufile "github.com/mnogokotin/golang-packages/utils/file"
	"log"
	"os"
	"time"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go task(ctx, cfg.Service.UpdateInterval, cfg.Service.OutputFilePath, cfg.Service.CsvSeparator, pg)

	time.Sleep(cfg.Service.CancelInterval)
	cancel()
}

func (p *program) Stop(s service.Service) error {
	<-time.After(1 * time.Second)
	return nil
}

func Run() {
	svcConfig := &service.Config{
		Name:        "GolangService",
		DisplayName: "Golang service",
	}

	p := &program{}
	s, err := service.New(p, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func task(ctx context.Context, updateInterval time.Duration, outputFilePath string, csvSeparator string, pg *postgres.Postgres) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			func() {
				f, err := file.OpenOrCreateFileOnRead(outputFilePath)
				defer f.Close()
				if err != nil {
					log.Fatal(err)
				}

				lastLine, err := ufile.GetLastLine(f)
				lastLineId := ""
				if err != nil {
					var emptyFileError *ufile.EmptyFileError
					if errors.As(err, &emptyFileError) {
						lastLineId = "0"
					} else {
						log.Fatal(err)
					}
				} else {
					lastLineId = csv.GetIdFromLine(lastLine, csvSeparator)
				}

				eventModels := database.GetEventModelsWithGreaterId(pg, lastLineId)
				if len(eventModels) > 0 {
					f, err := file.OpenFileOnWriteAtTheEnd(outputFilePath)
					defer f.Close()
					if err != nil {
						log.Fatal(err)
					}

					err = file.WriteModelsToFile(f, csvSeparator, eventModels)
					if err != nil {
						log.Fatal(err)
					}
				}
			}()
		}
		time.Sleep(updateInterval)
	}
}
