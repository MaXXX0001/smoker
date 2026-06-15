// Command chronos — provider-сервіс "Час і числа".
package main

import (
	"smoker/internal/chronos/infra"
	conditionv1 "smoker/pkg/proto/condition/v1"
	"smoker/pkg/env"
	"smoker/pkg/grpcx"
	"smoker/pkg/httpx"
	"smoker/pkg/provider"

	"google.golang.org/grpc"
)

func main() {
	log := env.Logger("chronos")
	addr := env.String("GRPC_ADDR", ":9102")
	country := env.String("COUNTRY_CODE", "UA")
	wikiLang := env.String("WIKI_LANG", "uk")

	hc := httpx.New()
	srv := provider.NewServer(log,
		infra.NewLocalSource(),
		infra.NewHolidaySource(hc, country),
		infra.NewOnThisDaySource(hc, wikiLang),
	)

	if err := grpcx.Serve(addr, log, func(g *grpc.Server) {
		conditionv1.RegisterConditionProviderServer(g, srv)
	}); err != nil {
		log.Error("fatal", "err", err)
	}
}
