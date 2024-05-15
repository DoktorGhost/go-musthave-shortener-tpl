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
)

func StartServer(conf *config.Config) error {
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
	sugar.Infow("server started", "addr", conf.Host+":"+conf.Port)

	//чтение конфигурациооного файла бд
	cons, err := osfile.NewConsumer(conf.FileStoragePath)

	if err != nil {
		log.Println("ошибка чтения конфигурациооного файла", err)
	} else {
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
	}

	//инициализация роутов
	r := handlers.InitRoutes(*shortURLUseCase, conf)

	//запускаем сервер
	err = http.ListenAndServe(":"+conf.Port, r)

	if err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
		return err
	}
	return nil
}
