package main

import (
	"context"
	"log/slog"
	"notion_saver/src/routes"
	"notion_saver/src/utils"
	blocksWorker "notion_saver/src/workers/blocks"
	pagesWorker "notion_saver/src/workers/pages"
	savesWorker "notion_saver/src/workers/saves"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Load environment variables
	env := utils.LoadEnv()

	traceProvider := utils.NewTracerProvider(ctx, env)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	logger := utils.NewLogger(ctx)
	slog.SetDefault(logger)

	// Create the server
	server := utils.Server{
		Logger:  logger,
		Env:     env,
		Context: ctx,
	}

	server.ConnectDatabase()
	server.ConnectQueue()

	saveWorker := savesWorker.NewSaveWorker(&server)
	go saveWorker.Run(&server)

	pageWorker := pagesWorker.NewPageWorker(&server)
	go pageWorker.Run(&server)

	blockWorker := blocksWorker.NewBlockWorker(&server)
	go blockWorker.Run(&server)

	engine := gin.Default()
	routes.Router(&server, engine)
	engine.Run(":8080")
}
