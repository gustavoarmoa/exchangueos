//go:build grpcgen

package api

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/revenu-tech/exchangeos/modules/compliance/application"
	"github.com/revenu-tech/exchangeos/modules/compliance/domain"
	"github.com/revenu-tech/exchangeos/pkg/bacen"
	pb "github.com/revenu-tech/exchangeos/proto/gen/exchangeos/v1"
)

type GRPCServer struct{ svc *application.Service }

func NewGRPCServer(svc *application.Service) *GRPCServer { return &GRPCServer{svc: svc} }

// ClassifyOperation — pb.ComplianceServiceServer
// Proto carries only trade_id; the hint comes from a downstream description lookup
// (placeholder uses trade_id as a hint). When the proto adds a `hint` field this
// adapter switches to it directly.
func (s *GRPCServer) ClassifyOperation(ctx context.Context, req *pb.ClassifyOperationRequest) (*pb.ClassifyOperationResponse, error) {
	tid, err := parseTenant(req.GetTenant())
	if err != nil {
		return nil, err
	}
	trid, err := uuid.Parse(req.GetTradeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "trade_id: %v", err)
	}
	// Placeholder hint: real wiring resolves hint from trade.External_Ref or refdata.
	c, err := s.svc.ClassifyOperation(ctx, tid, trid, "10001")
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.ClassifyOperationResponse{Classification: &pb.OperationClassification{
		ClassificationId: c.ID().String(),
		TradeId:          c.TradeID().String(),
		Code:             c.Code(),
		Description:      c.Description(),
		Nature:           string(c.Nature()),
	}}, nil
}

// ComputeIOF — pb.ComplianceServiceServer
func (s *GRPCServer) ComputeIOF(ctx context.Context, req *pb.ComputeIOFRequest) (*pb.ComputeIOFResponse, error) {
	tid, err := parseTenant(req.GetTenant())
	if err != nil {
		return nil, err
	}
	trid, err := uuid.Parse(req.GetTradeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "trade_id: %v", err)
	}
	// Defaults until proto carries op/notional/ccy; placeholder for smoke wiring.
	iof, err := s.svc.ComputeIOF(ctx, tid, trid, "DEFAULT", decimal.NewFromInt(10000), "USD")
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.ComputeIOFResponse{Iof: &pb.IOFComputation{
		IofId:        iof.ID().String(),
		TradeId:      iof.TradeID().String(),
		IofAmount:    &pb.Money{Amount: iof.IOFAmount().String(), Currency: iof.NotionalCCY()},
		RateApplied:  iof.Rate().String(),
		OperationType: iof.OperationType(),
	}}, nil
}

// SubmitBACENReport — pb.ComplianceServiceServer
func (s *GRPCServer) SubmitBACENReport(ctx context.Context, req *pb.SubmitBACENReportRequest) (*pb.SubmitBACENReportResponse, error) {
	tid, err := parseTenant(req.GetTenant())
	if err != nil {
		return nil, err
	}
	rep := req.GetReport()
	if rep == nil {
		return nil, status.Error(codes.InvalidArgument, "report required")
	}
	r, err := s.svc.SubmitBACENReport(ctx, domain.NewBACENReportInput{
		TenantID:      tid,
		ReportType:    domain.ReportType(rep.GetReportType()),
		ReferenceDate: rep.GetReferenceDate().AsTime(),
		PayloadHash:   payloadHash(rep.GetPayload()),
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.SubmitBACENReportResponse{Report: &pb.BACENReport{
		ReportId:   r.ID().String(),
		TenantId:   r.TenantID().String(),
		ReportType: string(r.Type()),
		Status:     string(r.Status()),
	}}, nil
}

// ScreenCounterparty — pb.ComplianceServiceServer
func (s *GRPCServer) ScreenCounterparty(ctx context.Context, req *pb.ScreenCounterpartyRequest) (*pb.ScreenCounterpartyResponse, error) {
	tid, err := parseTenant(req.GetTenant())
	if err != nil {
		return nil, err
	}
	// Stub: empty hits → RiskLow. Real screening calls list providers downstream.
	res, err := s.svc.ScreenCounterparty(ctx, domain.NewScreeningInput{
		TenantID: tid, CounterpartyBIC: req.GetCounterpartyBic(), LEI: req.GetLei(),
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.ScreenCounterpartyResponse{Result: &pb.ScreeningResult{
		Clear:     res.IsClear(),
		Hits:      res.Hits(),
		RiskLevel: string(res.RiskLevel()),
	}}, nil
}

func parseTenant(t *pb.TenantContext) (uuid.UUID, error) {
	if t == nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "tenant context required")
	}
	tid, err := uuid.Parse(t.GetTenantId())
	if err != nil {
		return uuid.Nil, status.Errorf(codes.InvalidArgument, "tenant_id: %v", err)
	}
	return tid, nil
}

// payloadHash is a placeholder — real implementation: sha256 of canonical payload.
func payloadHash(payload []byte) string {
	if len(payload) == 0 {
		return "empty-payload"
	}
	// Minimal deterministic value so the domain constructor passes; production
	// must use crypto/sha256.
	return string(payload[:min(len(payload), 32)])
}

func min(a, b int) int { if a < b { return a }; return b }

func mapErr(err error) error {
	switch {
	case errors.Is(err, application.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, application.ErrInvalidInput),
		errors.Is(err, domain.ErrInvalidInput),
		errors.Is(err, bacen.ErrUnknown):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrInvalidTransition):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// Force time import (used by upstream marshalling in proto adapters).
var _ = time.Now
