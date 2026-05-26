// Package domain — RefData bounded context.
//
// Holds 4 aggregates: Currency, Calendar, BICRecord, SSI.
// All are read-mostly; mutations are admin-only and go through application services.
//
// Conventions (cite .claude/rules/modules-domain.md):
//
//   - Aggregates expose constructors that validate; private fields, pointer-method receivers.
//   - No infrastructure imports.
//   - Decimal types via shopspring/decimal where money/rate appear.
package domain
