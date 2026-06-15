package infra

import (
	"context"
	"time"

	"smoker/internal/orchestrator/app"
	advisorv1 "smoker/pkg/proto/advisor/v1"
	"smoker/pkg/smoke"
)

// AdvisorServer реалізує контракт SmokeAdvisor поверх застосункового сервісу.
type AdvisorServer struct {
	advisorv1.UnimplementedSmokeAdvisorServer
	svc *app.Service
}

func NewAdvisorServer(svc *app.Service) *AdvisorServer {
	return &AdvisorServer{svc: svc}
}

// Recommend — реалізація RPC: делегує в app і мапить результат у proto.
func (s *AdvisorServer) Recommend(ctx context.Context, req *advisorv1.RecommendRequest) (*advisorv1.RecommendResponse, error) {
	loc := smoke.LocationFromProto(req.GetLocation())
	t := time.Unix(req.GetUnixTs(), 0)
	if req.GetUnixTs() == 0 {
		t = time.Now()
	}

	res := s.svc.Recommend(ctx, loc, t)
	rec := res.Recommendation
	return &advisorv1.RecommendResponse{
		Decision:   smoke.DecisionToProto(rec.Decision),
		TotalScore: int32(rec.TotalScore),
		Confidence: rec.Confidence,
		Reasons:    smoke.ConditionsToProto(rec.Reasons),
		Message:    res.Message,
	}, nil
}
