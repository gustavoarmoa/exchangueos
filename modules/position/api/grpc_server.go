//go:build grpcgen

package api

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/revenu-tech/exchangeos/modules/position/application"
	"github.com/revenu-tech/exchangeos/modules/position/domain"
	pb "github.com/revenu-tech/exchangeos/proto/gen/exchangeos/v1"
)

type GRPCServer struct {
	svc *application.Service
}

func NewGRPCServer(svc *application.Service) *GRPCServer { return &GRPCServer{svc: svc} }

// GetPosition — pb.PositionServiceServer
func (s *GRPCServer) GetPosition(ctx context.Context, req *pb.GetPositionRequest) (*pb.GetPositionResponse, error) {
	tid, err := parseTenant(req.GetTenant())
	if err != nil {
		return nil, err
	}
	p, err := s.svc.Get(ctx, tid, req.GetCurrency())
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.GetPositionResponse{Position: toPB(p)}, nil
}

// ListPositions — pb.PositionServiceServer
func (s *GRPCServer) ListPositions(ctx context.Context, req *pb.ListPositionsRequest) (*pb.ListPositionsResponse, error) {
	tid, err := parseTenant(req.GetTenant())
	if err != nil {
		return nil, err
	}
	list, err := s.svc.List(ctx, tid)
	if err != nil {
		return nil, mapErr(err)
	}
	out := make([]*pb.Position, 0, len(list))
	for _, p := range list {
		out = append(out, toPB(p))
	}
	return &pb.ListPositionsResponse{Positions: out}, nil
}

// RecomputePositions — pb.PositionServiceServer
// Stub: returns 0 updated. Real recomputation reads trades within business_date
// and replays via ApplyTradeLeg — TODO MS-023g.
func (s *GRPCServer) RecomputePositions(_ context.Context, req *pb.RecomputePositionsRequest) (*pb.RecomputePositionsResponse, error) {
	if _, err := parseTenant(req.GetTenant()); err != nil {
		return nil, err
	}
	return &pb.RecomputePositionsResponse{
		PositionsUpdated: 0,
		CompletedAt:      timestamppb.New(time.Now().UTC()),
	}, nil
}

func toPB(p *domain.Position) *pb.Position {
	return &pb.Position{
		PositionId: p.ID().String(),
		TenantId:   p.TenantID().String(),
		Currency:   p.Currency(),
		LongAmount: &pb.Money{Amount: p.Long().String(), Currency: p.Currency()},
		ShortAmount: &pb.Money{Amount: p.Short().String(), Currency: p.Currency()},
		NetAmount:   &pb.Money{Amount: p.Net().String(), Currency: p.Currency()},
		AsOf:        timestamppb.New(p.AsOf()),
	}
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

func mapErr(err error) error {
	switch {
	case errors.Is(err, application.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, application.ErrInvalidInput), errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
