//go:build grpcgen

package main

import (
	"google.golang.org/grpc"

	"github.com/revenu-tech/exchangeos/internal/container"
	adminapi "github.com/revenu-tech/exchangeos/modules/admin/api"
	setlapi "github.com/revenu-tech/exchangeos/modules/cls_settlement/api"
	complapi "github.com/revenu-tech/exchangeos/modules/compliance/api"
	posapi "github.com/revenu-tech/exchangeos/modules/position/api"
	quoteapi "github.com/revenu-tech/exchangeos/modules/quote/api"
	refapi "github.com/revenu-tech/exchangeos/modules/refdata/api"
	riskapi "github.com/revenu-tech/exchangeos/modules/risk/api"
	tradeapi "github.com/revenu-tech/exchangeos/modules/trade/api"
	pb "github.com/revenu-tech/exchangeos/proto/gen/exchangeos/v1"
)

// registerGeneratedServices wires the gRPC adapter for each bounded context to the
// shared application container. Only compiled when -tags grpcgen is passed AND
// proto/gen/exchangeos/v1 exists.
func registerGeneratedServices(srv *grpc.Server, c *container.Container) {
	pb.RegisterRefDataServiceServer(srv, refapi.NewGRPCServer(c.RefData))
	pb.RegisterQuoteServiceServer(srv, quoteapi.NewGRPCServer(c.Quote))
	pb.RegisterTradeServiceServer(srv, tradeapi.NewGRPCServer(c.Trade))
	pb.RegisterSettlementServiceServer(srv, setlapi.NewGRPCServer(c.Settlement, c.PayIn, c.NetReport))
	pb.RegisterRiskServiceServer(srv, riskapi.NewGRPCServer(c.Risk))
	pb.RegisterPositionServiceServer(srv, posapi.NewGRPCServer(c.Position))
	pb.RegisterComplianceServiceServer(srv, complapi.NewGRPCServer(c.Compliance))
	pb.RegisterAdminServiceServer(srv, adminapi.NewGRPCServer(c.Admin))
	// All 8 bounded-context services are now bound. CFETS Capture + Confirmation
	// have no public proto service yet (intentionally internal to MS-023d2 flow).
}
