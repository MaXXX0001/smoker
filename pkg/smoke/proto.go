package smoke

import (
	advisorv1 "smoker/pkg/proto/advisor/v1"
	conditionv1 "smoker/pkg/proto/condition/v1"
)

// Цей файл — єдине місце, де доменні типи зустрічаються з gRPC-контрактом.
// Решта домену лишається proto-free.

// LocationFromProto мапить protobuf Location у доменний.
func LocationFromProto(p *conditionv1.Location) Location {
	if p == nil {
		return Location{}
	}
	return Location{Lat: p.GetLat(), Lon: p.GetLon(), Name: p.GetName(), TZ: p.GetTz()}
}

// LocationToProto мапить доменний Location у protobuf.
func LocationToProto(l Location) *conditionv1.Location {
	return &conditionv1.Location{Lat: l.Lat, Lon: l.Lon, Name: l.Name, Tz: l.TZ}
}

func verdictToProto(v Verdict) conditionv1.Verdict {
	switch v {
	case Favorable:
		return conditionv1.Verdict_VERDICT_FAVORABLE
	case Unfavorable:
		return conditionv1.Verdict_VERDICT_UNFAVORABLE
	default:
		return conditionv1.Verdict_VERDICT_NEUTRAL
	}
}

func verdictFromProto(v conditionv1.Verdict) Verdict {
	switch v {
	case conditionv1.Verdict_VERDICT_FAVORABLE:
		return Favorable
	case conditionv1.Verdict_VERDICT_UNFAVORABLE:
		return Unfavorable
	default:
		return Neutral
	}
}

// ConditionToProto мапить доменну умову у protobuf.
func ConditionToProto(c Condition) *conditionv1.Condition {
	return &conditionv1.Condition{
		Code:     c.Code,
		Category: c.Category,
		Verdict:  verdictToProto(c.Verdict),
		Score:    int32(c.Score),
		Headline: c.Headline,
	}
}

// ConditionFromProto мапить protobuf-умову у доменну.
func ConditionFromProto(p *conditionv1.Condition) Condition {
	return Condition{
		Code:     p.GetCode(),
		Category: p.GetCategory(),
		Verdict:  verdictFromProto(p.GetVerdict()),
		Score:    int(p.GetScore()),
		Headline: p.GetHeadline(),
	}
}

// ConditionsToProto / ConditionsFromProto — пакетні версії.
func ConditionsToProto(cs []Condition) []*conditionv1.Condition {
	out := make([]*conditionv1.Condition, 0, len(cs))
	for _, c := range cs {
		out = append(out, ConditionToProto(c))
	}
	return out
}

func ConditionsFromProto(ps []*conditionv1.Condition) []Condition {
	out := make([]Condition, 0, len(ps))
	for _, p := range ps {
		out = append(out, ConditionFromProto(p))
	}
	return out
}

// DecisionToProto мапить доменне рішення у protobuf.
func DecisionToProto(d Decision) advisorv1.Decision {
	if d == Go {
		return advisorv1.Decision_DECISION_GO
	}
	return advisorv1.Decision_DECISION_WAIT
}

// DecisionFromProto мапить protobuf-рішення у доменне.
func DecisionFromProto(d advisorv1.Decision) Decision {
	if d == advisorv1.Decision_DECISION_GO {
		return Go
	}
	return Wait
}
