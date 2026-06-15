// Command cosmos — provider-сервіс "Природа та космос". Збирає погодні/космічні
// джерела та віддає їх через спільний контракт ConditionProvider.
package main

import (
	"smoker/internal/cosmos/infra"
	conditionv1 "smoker/pkg/proto/condition/v1"
	"smoker/pkg/env"
	"smoker/pkg/grpcx"
	"smoker/pkg/httpx"
	"smoker/pkg/provider"

	"google.golang.org/grpc"
)

func main() {
	log := env.Logger("cosmos")
	addr := env.String("GRPC_ADDR", ":9101")

	hc := httpx.New()
	srv := provider.NewServer(log,
		infra.NewMoonSource(),
		infra.NewWeatherSource(hc),
		infra.NewPollenSource(hc),
		infra.NewKpSource(hc),
		infra.NewISSSource(hc),
	)

	if err := grpcx.Serve(addr, log, func(g *grpc.Server) {
		conditionv1.RegisterConditionProviderServer(g, srv)
	}); err != nil {
		log.Error("fatal", "err", err)
	}
}
