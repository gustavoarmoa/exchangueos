---
description: Run full security scan (SAST + SCA + secrets + container + IaC)
allowed-tools: [Bash]
---

# /security-scan

Roda full security pipeline (Tier 2 equivalent):
- gosec (SAST)
- govulncheck (CVE Go deps)
- osv-scanner (CVE lockfiles)
- gitleaks (secrets)
- trivy fs (filesystem CVE)
- trivy image (container scan se imagem built)
- checkov + tfsec (IaC se Terraform mudou)
