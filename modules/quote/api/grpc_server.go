//go:build grpcgen

package api

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/revenu-tech/exchangeos/modules/quote/application"
	"github.com/revenu-tech/exchangeos/modules/quote/domain"
	pb "github.com/revenu-tech/exchangeos/proto/gen/exchangeos/v1"
)

// GRPCServer adapts pb.QuoteServiceServer to the application service.
type GRPCServer struct {
	svc *application.Service
}

func NewGRPCServer(svc *application.Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

// GetQuote — pb.QuoteServiceServer
func (s *GRPCServer) GetQuote(ctx context.Context, req *pb.GetQuoteRequest) (*pb.GetQuoteResponse, error) {
	tid, err := parseTenant(req.GetTenant())
	if err != nil {
		return nil, err
	}
	if req.GetNotional() == nil {
		return nil, status.Error(codes.InvalidArgument, "notional required")
	}
	notional, err := decimal.NewFromString(req.GetNotional().GetAmount())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "notional.amount: %v", err)
	}

	q, err := s.svc.GetQuote(ctx, application.GetQuoteRequest{
		TenantID:    tid,
		BaseCCY:     req.GetBaseCcy(),
		QuoteCCY:    req.GetQuoteCcy(),
		Notional:    notional,
		NotionalCCY: req.GetNotional().GetCurrency(),
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.GetQuoteResponse{Quote: toPBQuote(q)}, nil
}

// AcceptQuote — pb.QuoteServiceServer
// Returns the quote_id back as the response trade_id. Real trade creation is
// driven downstream by a worker reacting to the `quote.accepted.v1` event.
func (s *GRPCServer) AcceptQuote(ctx context.Context, req *pb.AcceptQuoteRequest) (*pb.AcceptQuoteResponse, error) {
	if _, err := parseTenant(req.GetTenant()); err != nil {
		return nil, err
	}
	qid, err := uuid.Parse(req.GetQuoteId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "quote_id: %v", err)
	}
	q, err := s.svc.AcceptQuote(ctx, application.AcceptQuoteRequest{
		QuoteID: qid,
		Actor:   req.GetTenant().GetActorId(),
	})
	if err != nil {
		return nil, mapErr(err)
	}
	// Placeholder: until trade-creation worker is implemented, return the quote
	// id as the trade id. Caller correlates via event log.
	return &pb.AcceptQuoteResponse{TradeId: q.ID().String()}, nil
}

// StreamQuotes — pb.QuoteServiceServer streaming RPC (TODO).
func (s *GRPCServer) StreamQuotes(_ *pb.StreamQuotesRequest, _ pb.QuoteService_StreamQuotesServer) error {
	return status.Error(codes.Unimplemented, "streaming quotes not yet implemented")
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

func toPBQuote(q *domain.Quote) *pb.Quote {
	return &pb.Quote{
		QuoteId:   q.ID().String(),
		TenantId:  q.TenantID().String(),
		BaseCcy:   q.BaseCCY(),
		QuoteCcy:  q.QuoteCCY(),
		Notional:  &pb.Money{Amount: q.Notional().String(), Currency: q.NotionalCCY()},
		Bid:       &pb.FxRate{Rate: q.Bid().String(), BaseCcy: q.BaseCCY(), QuoteCcy: q.QuoteCCY()},
		Ask:       &pb.FxRate{Rate: q.Ask().String(), BaseCcy: q.BaseCCY(), QuoteCcy: q.QuoteCCY()},
		ValidFrom: timestamppb.New(q.ValidFrom()),
		ValidTo:   timestamppb.New(q.ValidTo()),
	}
}

func mapErr(err error) error {
	switch {
	case errors.Is(err, application.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, application.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrQuoteExpired):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrInvalidTransition):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
