package server

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/osfile"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage/postgres"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

func StartServer(conf *config.Config) error {
	//логирование
	logg, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logg.Sync()
	logger.InitLogger(logg)
	sugar := *logg.Sugar()
	sugar.Infow("server started", "addr", conf.Host+":"+conf.Port)

	var shortURLUseCase *usecase.ShortURLUseCase

	//подключение к БД
	if conf.DatabaseDSN != "" {

		db, err := postgres.NewPostgresStorage(conf.DatabaseDSN)
		if err != nil {
			sugar.Fatalw("Ошибка при подключении к БД", "error", err)
		}
		sugar.Infow("Успешное подключение к БД")
		shortURLUseCase = usecase.NewShortURLUseCase(db)

	} else {
		db := maps.NewMapStorage()
		shortURLUseCase = usecase.NewShortURLUseCase(db)
		sugar.Infow("Использование оперативной памяти вместо БД")
	}

	//чтение конфигурациооного файла бд
	cons, err := osfile.NewConsumer(conf.FileStoragePath)

	if err != nil {
		log.Println("ошибка чтения конфигурационного файла", err)
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
			err = shortURLUseCase.Write(event.Short, event.ShortURL, event.OriginalURL, event.UserID)
			if err != nil {
				log.Println(err)
				continue
			}
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
