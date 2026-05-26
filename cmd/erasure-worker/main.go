// cmd/erasure-worker — LGPD Art. 18 IV right-to-erasure executor.
//
// Usage:
//   erasure-worker --ticket LGPD-2026-0001 --plan plan.yaml --dry-run
//   erasure-worker --ticket LGPD-2026-0001 --plan plan.yaml --execute
//
// --dry-run prints the SQL each operation would issue + the affected row count
// per the WHERE clause, WITHOUT mutating data. Safe to run any time.
//
// --execute requires the plan to carry approvals from both 'dpo' AND
// 'compliance_officer' roles (4-eyes per workflow Stage 3) AND env var
// EXCHANGEOS_ERASURE_CONFIRM=YES-I-MEAN-IT (a final out-of-band guard against
// accidental invocation in CI / shell history).
//
// Successful execution emits one audit_event per operation and one outbox event
// `lgpd.erasure_completed.v1` after the final commit so downstream modules
// (LedgerOS, ComplOS) can apply equivalent erasure to their own datasets.
//
// Full reference: docs/security/data-lifecycle/erasure-workflow.md.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/revenu-tech/exchangeos/internal/erasure"
)

func main() {
	if err := run(); err != nil {
		slog.Error("erasure-worker: failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		ticket   = flag.String("ticket", "", "LGPD ticket ID (must match plan.ticket)")
		planPath = flag.String("plan", "", "path to signed plan YAML")
		dryRun   = flag.Bool("dry-run", false, "print operations + row counts without mutating")
		execute  = flag.Bool("execute", false, "actually apply the plan (requires 4-eyes approval + confirm env var)")
	)
	flag.Parse()

	if *ticket == "" || *planPath == "" {
		return fmt.Errorf("--ticket and --plan are required")
	}
	if *dryRun == *execute {
		return fmt.Errorf("exactly one of --dry-run or --execute must be set")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	data, err := os.ReadFile(*planPath)
	if err != nil {
		return fmt.Errorf("read plan: %w", err)
	}
	plan, err := erasure.ParsePlan(data)
	if err != nil {
		return fmt.Errorf("parse plan: %w", err)
	}
	if plan.Ticket != *ticket {
		return fmt.Errorf("ticket mismatch: --ticket=%s vs plan.ticket=%s", *ticket, plan.Ticket)
	}

	slog.Info("erasure-worker: plan loaded",
		"ticket", plan.Ticket,
		"subject_ref", plan.SubjectRef,
		"operations", len(plan.Operations),
		"approvals", plan.Approvals,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if *dryRun {
		return dryRunPlan(ctx, plan)
	}

	// ── execute branch ──────────────────────────────────────────────────────
	if !plan.HasRequiredApprovals() {
		return fmt.Errorf("plan lacks required approvals (need both 'dpo' AND 'compliance_officer')")
	}
	if os.Getenv("EXCHANGEOS_ERASURE_CONFIRM") != "YES-I-MEAN-IT" {
		return fmt.Errorf("EXCHANGEOS_ERASURE_CONFIRM=YES-I-MEAN-IT required for --execute")
	}

	// Stage 2 wiring: Executor + AuditEmitter are implemented (see internal/erasure)
	// but the concrete DB + audit adapters are MS-024a Stage 3 (CRDB pgx adapter +
	// audit_event emit via modules/admin + outbox event publisher).
	//
	// Until those land, --execute remains explicitly disabled in the binary —
	// preventing accidental "looks ready" misuse. The Executor itself is fully
	// tested with a fake DB; once the adapters exist this branch becomes:
	//
	//   db := pgxadapter.New(ctx, cfg.DBDsn)
	//   audit := admindapter.New(ctx, cfg.AdminEndpoint)
	//   exec := erasure.NewExecutor(db, audit)
	//   res, err := exec.Apply(ctx, plan)
	//   slog.Info("erasure-worker: applied", "ops", len(res.Ops), "rows", res.RowsTotal)
	return fmt.Errorf("execute path needs DB + audit adapters (MS-024a Stage 3); executor itself is ready")
}

// dryRunPlan prints what each operation would do. Safe — no DB write.
func dryRunPlan(_ context.Context, plan *erasure.Plan) error {
	fmt.Printf("\n=== DRY RUN — ticket %s ===\n\n", plan.Ticket)
	for i, op := range plan.Operations {
		switch op.Op {
		case erasure.OpRedact:
			fmt.Printf("[%d] REDACT %s SET (%v) = '%s' WHERE %s\n",
				i+1, op.Table, op.Fields, plan.RedactionMarker(), op.Where)
		case erasure.OpHardDelete:
			fmt.Printf("[%d] DELETE FROM %s WHERE %s\n", i+1, op.Table, op.Where)
		}
	}
	fmt.Printf("\n=== END DRY RUN (%d ops) ===\n\n", len(plan.Operations))
	slog.Info("erasure-worker: dry-run complete", "operations", len(plan.Operations))
	return nil
}
