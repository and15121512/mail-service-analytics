package application

import (
	"context"
	"fmt"

	"github.com/TheZeroSlave/zapsentry"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/adapters/auth_grpc"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/adapters/grpc"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/adapters/http"
	postgresdb "gitlab.com/sukharnikov.aa/mail-service-analytics/internal/adapters/postgres"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/domain/analytics"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

var (
	hs     *http.Server
	gs     *grpc.Server
	logger *zap.Logger
)

func modifyToSentryLogger(log *zap.Logger, DSN string) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, //when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
		Tags: map[string]string{
			"component": "system",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(DSN))

	// to use breadcrumbs feature - create new scope explicitly
	log = log.With(zapsentry.NewScope())

	//in case of err it will return noop core. so we can safely attach it
	if err != nil {
		log.Warn("failed to init zap", zap.Error(err))
	}
	return zapsentry.AttachCoreToLogger(core, log)
}

func Start(ctx context.Context) {
	logger, _ = zap.NewProduction()
	//sentryClient, err := sentry.NewClient(sentry.ClientOptions{
	//	Dsn: "http://b7dd7b3ce3df4f2b81f5af622512658c@localhost:9000/2",
	//})
	//if err != nil {
	//	logger.Sugar().Fatalf("http server creating failed: %s", err)
	//}
	//defer sentryClient.Flush(2 * time.Second)
	logger = modifyToSentryLogger(logger, "http://b7dd7b3ce3df4f2b81f5af622512658c@localhost:9000/2")

	connectionString := "user=postgres password=secret host=postgres-db-analytics port=5432 dbname=analytics sslmode=disable pool_max_conns=10"
	db, err := postgresdb.New(ctx, connectionString, logger.Sugar())
	if err != nil {
		logger.Sugar().Fatalf("postgres init failed: %s", err)
	}

	ac := auth_grpc.New(logger.Sugar())
	analyticsS := analytics.New(db, ac, logger.Sugar())

	hs, err = http.New(analyticsS, logger.Sugar())
	if err != nil {
		logger.Sugar().Fatalf("http server creating failed: %s", err)
	}
	gs, err = grpc.New(logger.Sugar(), analyticsS)
	if err != nil {
		logger.Sugar().Fatalf("grpc server creating failed: %s", err)
	}

	var g errgroup.Group
	g.Go(func() error {
		return hs.Start()
	})
	g.Go(func() error {
		return gs.Start()
	})

	logger.Sugar().Info(fmt.Sprintf("app is started on ports: %d (http) and %d (grpc)", hs.Port(), gs.Port()))
	err = g.Wait()
	if err != nil {
		logger.Sugar().Fatalw("http server start failed", zap.Error(err))
	}
}

func Stop() {
	_ = hs.Stop(context.Background())
	_ = gs.Stop(context.Background())
	logger.Sugar().Info("app has stopped")
}
