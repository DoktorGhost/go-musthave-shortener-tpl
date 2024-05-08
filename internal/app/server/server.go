package server

import (
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"go.uber.org/zap"
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
