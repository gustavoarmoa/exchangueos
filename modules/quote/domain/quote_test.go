package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/revenu-tech/exchangeos/modules/quote/domain"
)

func dec(s string) decimal.Decimal { return decimal.RequireFromString(s) }

func validQuoteInput(t *testing.T) domain.NewQuoteInput {
	t.Helper()
	now := time.Now().UTC()
	return domain.NewQuoteInput{
		TenantID:    uuid.New(),
		BaseCCY:     "eur",
		QuoteCCY:    "usd",
		Notional:    dec("1000000"),
		NotionalCCY: "EUR",
		Bid:         dec("1.0798"),
		Ask:         dec("1.0802"),
		ValidFrom:   now,
		ValidTo:     now.Add(10 * time.Second),
		Venue:       "INTERNAL",
	}
}

func TestQuote_Valid(t *testing.T) {
	q, err := domain.NewQuote(validQuoteInput(t))
	if err != nil {
		t.Fatalf("NewQuote: %v", err)
	}
	if q.BaseCCY() != "EUR" || q.QuoteCCY() != "USD" {
		t.Errorf("ccy normalised: base=%s quote=%s", q.BaseCCY(), q.QuoteCCY())
	}
	if !q.Mid().Equal(dec("1.0800")) {
		t.Errorf("mid: got %s want 1.0800", q.Mid())
	}
	if got := q.Version(); got != 1 {
		t.Errorf("version: got %d want 1", got)
	}
	if got := len(q.PendingEvents()); got != 1 {
		t.Errorf("pending events: got %d want 1", got)
	}
}

func TestQuote_BidExceedsAsk_Rejected(t *testing.T) {
	in := validQuoteInput(t)
	in.Bid = dec("1.10")
	in.Ask = dec("1.09")
	_, err := domain.NewQuote(in)
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}
}

func TestQuote_NotionalCCYMustMatchPair(t *testing.T) {
	in := validQuoteInput(t)
	in.NotionalCCY = "GBP"
	_, err := domain.NewQuote(in)
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}
}

func TestQuote_AcceptWithinWindow(t *testing.T) {
	q, _ := domain.NewQuote(validQuoteInput(t))
	if err := q.Accept(time.Now().UTC(), "trader-a"); err != nil {
		t.Fatalf("Accept: %v", err)
	}
	events := q.PendingEvents()
	if len(events) != 2 {
		t.Fatalf("events: got %d want 2 (created + accepted)", len(events))
	}
	if events[1].EventName() != "quote.accepted.v1" {
		t.Fatalf("event[1]: got %s want quote.accepted.v1", events[1].EventName())
	}
}

func TestQuote_AcceptExpired(t *testing.T) {
	in := validQuoteInput(t)
	in.ValidFrom = time.Now().UTC().Add(-1 * time.Hour)
	in.ValidTo = time.Now().UTC().Add(-1 * time.Minute)
	q, _ := domain.NewQuote(in)
	if err := q.Accept(time.Now().UTC(), "trader-a"); !errors.Is(err, domain.ErrQuoteExpired) {
		t.Fatalf("want ErrQuoteExpired, got %v", err)
	}
}

func TestRFQ_Lifecycle_HappyPath(t *testing.T) {
	r, err := domain.NewRFQ(domain.NewRFQInput{
		TenantID:  uuid.New(),
		Requester: "trader-a",
		BaseCCY:   "EUR",
		QuoteCCY:  "USD",
	})
	if err != nil {
		t.Fatalf("NewRFQ: %v", err)
	}
	if r.Status() != domain.RFQRequested {
		t.Fatalf("initial: got %s", r.Status())
	}
	qid := uuid.New()
	if err := r.AttachQuote(qid); err != nil {
		t.Fatalf("AttachQuote: %v", err)
	}
	if r.Status() != domain.RFQQuoted {
		t.Fatalf("after attach: got %s", r.Status())
	}
	if err := r.Accept(qid, "trader-a"); err != nil {
		t.Fatalf("Accept: %v", err)
	}
	if r.Status() != domain.RFQAccepted {
		t.Fatalf("after accept: got %s", r.Status())
	}
	events := r.PendingEvents()
	wantNames := []string{"rfq.requested.v1", "rfq.quoted.v1", "rfq.accepted.v1"}
	if len(events) != len(wantNames) {
		t.Fatalf("events: got %d want %d", len(events), len(wantNames))
	}
	for i, n := range wantNames {
		if events[i].EventName() != n {
			t.Errorf("event[%d]: got %s want %s", i, events[i].EventName(), n)
		}
	}
}

func TestRFQ_AcceptUnknownQuote_Rejected(t *testing.T) {
	r, _ := domain.NewRFQ(domain.NewRFQInput{
		TenantID:  uuid.New(),
		Requester: "trader-a",
		BaseCCY:   "EUR",
		QuoteCCY:  "USD",
	})
	_ = r.AttachQuote(uuid.New())
	err := r.Accept(uuid.New(), "trader-a") // different quote
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}
}

func TestRFQ_RejectRequiresReason(t *testing.T) {
	r, _ := domain.NewRFQ(domain.NewRFQInput{
		TenantID:  uuid.New(),
		Requester: "trader-a",
		BaseCCY:   "EUR",
		QuoteCCY:  "USD",
	})
	if err := r.Reject(""); !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}
}

func TestRFQ_ExpireFromQuoted(t *testing.T) {
	r, _ := domain.NewRFQ(domain.NewRFQInput{
		TenantID:  uuid.New(),
		Requester: "trader-a",
		BaseCCY:   "EUR",
		QuoteCCY:  "USD",
	})
	_ = r.AttachQuote(uuid.New())
	if err := r.Expire(); err != nil {
		t.Fatalf("Expire: %v", err)
	}
	if r.Status() != domain.RFQExpired {
		t.Fatalf("after expire: %s", r.Status())
	}
}
