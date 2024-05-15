package server

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/osfile"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"strconv"
)

func StartServer(hostPort *config.Config) error {
	db := maps.NewMapStorage()
	shortURLUseCase := usecase.NewShortURLUseCase(db)

	//логирование
	logg, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logg.Sync()
	logger.InitLogger(logg)
	sugar := *logg.Sugar()
	sugar.Infow("server started", "addr", hostPort.String())

	//чтение конфигурациооного файла бд
	cons, err := osfile.NewConsumer(config.FileStoragePath)

	if err != nil {
		log.Println("ошибка чтения конфигурациооного файла", err)
		return err
	}
	for {
		event, err := cons.ReadEvent()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("ошибка чтения события", err)
			continue
		}
		if event == nil {
			break
		}
		shortURLUseCase.Write(event.OriginalURL, event.ShortURL)
	}

	//инициализация роутов
	r := handlers.InitRoutes(*shortURLUseCase)

	//запускаем сервер
	err = http.ListenAndServe(":"+strconv.Itoa(hostPort.Port), r)

	if err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
		return err
	}
	return nil
}
