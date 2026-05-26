// Package application — Compliance use cases (BACEN classification + IOF + reporting + screening).
package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/revenu-tech/exchangeos/modules/compliance/domain"
	"github.com/revenu-tech/exchangeos/pkg/bacen"
)

// Repositories — split per aggregate to keep each focused.
type ClassificationRepo interface {
	Save(ctx context.Context, c *domain.Classification) error
	GetByTrade(ctx context.Context, tradeID uuid.UUID) (*domain.Classification, error)
}

type IOFRepo interface {
	Save(ctx context.Context, i *domain.IOFComputation) error
	GetByTrade(ctx context.Context, tradeID uuid.UUID) (*domain.IOFComputation, error)
}

type ReportRepo interface {
	Save(ctx context.Context, r *domain.BACENReport) error
	Get(ctx context.Context, id uuid.UUID) (*domain.BACENReport, error)
}

type ScreeningRepo interface {
	Save(ctx context.Context, s *domain.ScreeningResult) error
}

var (
	ErrInvalidInput = errors.New("compliance-app: invalid input")
	ErrNotFound     = errors.New("compliance-app: not found")
)

// Service exposes compliance use cases.
type Service struct {
	classifier *bacen.Classifier
	iofCalc    *bacen.IOFCalculator
	classRepo  ClassificationRepo
	iofRepo    IOFRepo
	reportRepo ReportRepo
	screenRepo ScreeningRepo
}

func NewService(
	classifier *bacen.Classifier,
	iofCalc *bacen.IOFCalculator,
	classRepo ClassificationRepo,
	iofRepo IOFRepo,
	reportRepo ReportRepo,
	screenRepo ScreeningRepo,
) *Service {
	return &Service{
		classifier: classifier,
		iofCalc:    iofCalc,
		classRepo:  classRepo,
		iofRepo:    iofRepo,
		reportRepo: reportRepo,
		screenRepo: screenRepo,
	}
}

// ClassifyOperation runs the BACEN classifier against a free-text hint (or accepts
// an explicit code) and persists the Classification.
func (s *Service) ClassifyOperation(ctx context.Context, tenantID, tradeID uuid.UUID, codeOrHint string) (*domain.Classification, error) {
	if tenantID == uuid.Nil || tradeID == uuid.Nil {
		return nil, ErrInvalidInput
	}
	// Try exact-code lookup first; fall back to free-text classifier.
	nc, ok := s.classifier.ByCode(codeOrHint)
	if !ok {
		var err error
		nc, err = s.classifier.Classify(codeOrHint)
		if err != nil {
			return nil, err
		}
	}
	c, err := domain.NewClassification(domain.NewClassificationInput{
		TenantID: tenantID, TradeID: tradeID,
		Code: nc.Code, Description: nc.Description,
		Nature: domain.Nature(nc.Nature),
	})
	if err != nil {
		return nil, err
	}
	if err := s.classRepo.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// ComputeIOF runs the IOF calculator and persists.
func (s *Service) ComputeIOF(ctx context.Context, tenantID, tradeID uuid.UUID, opType string, notional decimal.Decimal, notionalCCY string) (*domain.IOFComputation, error) {
	if tenantID == uuid.Nil || tradeID == uuid.Nil {
		return nil, ErrInvalidInput
	}
	rate, _, err := s.iofCalc.Compute(opType, notional)
	if err != nil {
		return nil, err
	}
	iof, err := domain.NewIOFComputation(domain.NewIOFInput{
		TenantID: tenantID, TradeID: tradeID,
		OperationType: opType, Notional: notional, NotionalCCY: notionalCCY, Rate: rate,
	})
	if err != nil {
		return nil, err
	}
	if err := s.iofRepo.Save(ctx, iof); err != nil {
		return nil, err
	}
	return iof, nil
}

// SubmitBACENReport persists + immediately marks SUBMITTED. Real submission happens
// downstream (cmd/worker reacts to status change).
func (s *Service) SubmitBACENReport(ctx context.Context, in domain.NewBACENReportInput) (*domain.BACENReport, error) {
	r, err := domain.NewBACENReport(in)
	if err != nil {
		return nil, err
	}
	if err := s.reportRepo.Save(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

// ScreenCounterparty runs the screening (currently a stub — production replaces
// with calls to OFAC/UN/EU/COAF list providers) and persists.
func (s *Service) ScreenCounterparty(ctx context.Context, in domain.NewScreeningInput) (*domain.ScreeningResult, error) {
	res, err := domain.NewScreeningResult(in)
	if err != nil {
		return nil, err
	}
	if err := s.screenRepo.Save(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}
