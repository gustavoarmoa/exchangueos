package container_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/revenu-tech/exchangeos/internal/config"
	"github.com/revenu-tech/exchangeos/internal/container"

	quoteapp "github.com/revenu-tech/exchangeos/modules/quote/application"
)

func newMemCfg() *config.Config {
	return &config.Config{
		Env:   "dev",
		HTTP:  config.HTTPConfig{Port: 18094},
		GRPC:  config.GRPCConfig{Port: 19094, MaxRecvBytes: 1 << 20, MaxSendBytes: 1 << 20},
		Repos: config.ReposConfig{Backend: "memory"},
	}
}

func TestContainer_New_MemoryBackend(t *testing.T) {
	c, err := container.New(context.Background(), newMemCfg())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()
	if c.RefData == nil || c.Quote == nil || c.Trade == nil {
		t.Fatal("services not wired")
	}
	if c.SpotBook == nil {
		t.Fatal("spot book not wired")
	}
	if c.EventBus == nil {
		t.Fatal("event bus not wired")
	}
}

// End-to-end: GetQuote → AcceptQuote should publish quote.accepted.v1, which the
// container's eventbus handler routes to trade.BookTrade. Result: a new FXTrade
// exists in the trade repo with the correct economics.
func TestContainer_QuoteAccepted_BooksTrade(t *testing.T) {
	c, err := container.New(context.Background(), newMemCfg())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()
	ctx := context.Background()
	tenant := uuid.New()

	q, err := c.Quote.GetQuote(ctx, quoteapp.GetQuoteRequest{
		TenantID:    tenant,
		BaseCCY:     "EUR",
		QuoteCCY:    "USD",
		Notional:    decimal.NewFromInt(1_000_000),
		NotionalCCY: "EUR",
		Venue:       "INTERNAL",
		TTL:         5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("GetQuote: %v", err)
	}

	if _, err := c.Quote.AcceptQuote(ctx, quoteapp.AcceptQuoteRequest{
		QuoteID: q.ID(),
		Actor:   "trader-a",
	}); err != nil {
		t.Fatalf("AcceptQuote: %v", err)
	}

	// No trade must be booked. The Quote aggregate carries no counterparty BICs,
	// no side and no venue, and the handler used to fill those in with literals
	// ("DEUTDEFF" / "CHASUS33" / "CLS", rate = mid). This assertion previously
	// expected exactly one trade, which meant it was pinning the fabricated
	// counterparties as correct behaviour.
	//
	// Auto-booking returns here once the RFQ captures buyer/seller/side and the
	// Quote carries them; at that point this test flips back to asserting the
	// booked trade AND its counterparties.
	trades, err := c.Trade.ListTrades(ctx, tenant, "", time.Time{}, time.Time{}, 100)
	if err != nil {
		t.Fatalf("ListTrades: %v", err)
	}
	if len(trades) != 0 {
		t.Fatalf("expected no trade to be booked from a quote without counterparties, got %d", len(trades))
	}
}
