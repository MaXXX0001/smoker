package infra

import (
	"context"
	"time"

	advisorv1 "smoker/pkg/proto/advisor/v1"
	"smoker/pkg/grpcx"
	"smoker/pkg/smoke"

	"google.golang.org/grpc"
)

// AdvisorClient — gRPC-клієнт до orchestrator (SmokeAdvisor).
type AdvisorClient struct {
	conn   *grpc.ClientConn
	client advisorv1.SmokeAdvisorClient
}

// DialAdvisor відкриває з'єднання до orchestrator.
func DialAdvisor(target string) (*AdvisorClient, error) {
	conn, err := grpcx.Dial(target)
	if err != nil {
		return nil, err
	}
	return &AdvisorClient{conn: conn, client: advisorv1.NewSmokeAdvisorClient(conn)}, nil
}

func (a *AdvisorClient) Close() error { return a.conn.Close() }

// Advice — стислий результат для gateway.
type Advice struct {
	GoNow   bool
	Message string
}

// Recommend запитує вердикт у orchestrator для місця/моменту.
func (a *AdvisorClient) Recommend(ctx context.Context, loc smoke.Location, t time.Time) (Advice, error) {
	resp, err := a.client.Recommend(ctx, &advisorv1.RecommendRequest{
		Location: smoke.LocationToProto(loc),
		UnixTs:   t.Unix(),
	})
	if err != nil {
		return Advice{}, err
	}
	return Advice{
		GoNow:   resp.GetDecision() == advisorv1.Decision_DECISION_GO,
		Message: resp.GetMessage(),
	}, nil
}
