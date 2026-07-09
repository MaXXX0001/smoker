// Command chaos — provider-сервіс "Випадковість заради сміху".
package main

import (
	"smoker/internal/chaos/infra"
	conditionv1 "smoker/pkg/proto/condition/v1"
	"smoker/pkg/env"
	"smoker/pkg/grpcx"
	"smoker/pkg/httpx"
	"smoker/pkg/provider"

	"google.golang.org/grpc"
)

func main() {
	log := env.Logger("chaos")
	addr := env.String("GRPC_ADDR", ":9103")

	hc := httpx.New()
	srv := provider.NewServer(log,
		infra.NewDiceSource(),
		infra.NewJokeSource(),
		infra.NewCatFactSource(),
		infra.NewOracleSource(hc),
	)

	if err := grpcx.Serve(addr, log, func(g *grpc.Server) {
		conditionv1.RegisterConditionProviderServer(g, srv)
	}); err != nil {
		log.Error("fatal", "err", err)
	}
}
