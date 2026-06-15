// Command orchestrator — core domain сервіс: агрегує умови провайдерів і видає
// вердикт через контракт SmokeAdvisor.
package main

import (
	"smoker/internal/orchestrator/app"
	"smoker/internal/orchestrator/infra"
	advisorv1 "smoker/pkg/proto/advisor/v1"
	"smoker/pkg/env"
	"smoker/pkg/grpcx"

	"google.golang.org/grpc"
)

func main() {
	log := env.Logger("orchestrator")
	addr := env.String("GRPC_ADDR", ":9100")

	// Адреси провайдерів. Кожен — окремий bounded context/сервіс.
	targets := map[string]string{
		"cosmos":  env.String("COSMOS_ADDR", "localhost:9101"),
		"chronos": env.String("CHRONOS_ADDR", "localhost:9102"),
		"chaos":   env.String("CHAOS_ADDR", "localhost:9103"),
	}

	var providers []app.Provider
	for name, target := range targets {
		pc, err := infra.DialProvider(name, target)
		if err != nil {
			log.Error("не вдалось підключити провайдер", "name", name, "target", target, "err", err)
			continue
		}
		defer pc.Close()
		providers = append(providers, pc)
		log.Info("провайдер під'єднано", "name", name, "target", target)
	}

	svc := app.NewService(log, providers...)
	srv := infra.NewAdvisorServer(svc)

	if err := grpcx.Serve(addr, log, func(g *grpc.Server) {
		advisorv1.RegisterSmokeAdvisorServer(g, srv)
	}); err != nil {
		log.Error("fatal", "err", err)
	}
}
