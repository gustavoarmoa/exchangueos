# TDD Output Style

> Para responses durante TDD workflow

## Structure

```
## RED Phase
- Test that fails first: ...
- File: tests/...
- Run: `go test -run TestXxx ./...`
- Expected: FAIL ❌

## GREEN Phase  
- Minimal implementation: ...
- File: modules/...
- Run: same test
- Expected: PASS ✅

## REFACTOR Phase
- Cleanups: ...
- Run all tests: green ✅

## Coverage Impact
- Before: X% / After: Y%
```
