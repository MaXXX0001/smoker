// Package infra — gRPC-адаптери orchestrator: клієнти до провайдерів і сервер
// SmokeAdvisor.
package infra

import (
	"context"
	"time"

	conditionv1 "smoker/pkg/proto/condition/v1"
	"smoker/pkg/grpcx"
	"smoker/pkg/smoke"

	"google.golang.org/grpc"
)

// ProviderClient — gRPC-клієнт до одного ConditionProvider; реалізує app.Provider.
type ProviderClient struct {
	name   string
	conn   *grpc.ClientConn
	client conditionv1.ConditionProviderClient
}

// DialProvider відкриває (lazy) з'єднання до провайдера за адресою.
func DialProvider(name, target string) (*ProviderClient, error) {
	conn, err := grpcx.Dial(target)
	if err != nil {
		return nil, err
	}
	return &ProviderClient{
		name:   name,
		conn:   conn,
		client: conditionv1.NewConditionProviderClient(conn),
	}, nil
}

func (p *ProviderClient) Name() string { return p.name }

func (p *ProviderClient) Close() error { return p.conn.Close() }

// Evaluate викликає віддалений провайдер і мапить відповідь у домен.
func (p *ProviderClient) Evaluate(ctx context.Context, loc smoke.Location, t time.Time) ([]smoke.Condition, error) {
	resp, err := p.client.Evaluate(ctx, &conditionv1.EvaluateRequest{
		Location: smoke.LocationToProto(loc),
		UnixTs:   t.Unix(),
	})
	if err != nil {
		return nil, err
	}
	return smoke.ConditionsFromProto(resp.GetConditions()), nil
}
