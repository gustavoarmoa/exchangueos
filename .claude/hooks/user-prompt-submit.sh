#!/usr/bin/env bash
# .claude/hooks/user-prompt-submit.sh
# Triggered when user submits prompt — opportunity to inject context
set -euo pipefail

INPUT=$(cat)
PROMPT=$(echo "$INPUT" | jq -r '.prompt // ""' 2>/dev/null || echo "")

# Auto-detect intent and inject context hints (output to stdout = added to prompt)
HINTS=""

if [[ "$PROMPT" == *"trade"* || "$PROMPT" == *"FXTrade"* ]]; then
    HINTS+="\n📌 Context hint: Trade work usually requires fx-domain + iso20022 + bacen-compliance agents in parallel."
fi

if [[ "$PROMPT" == *"CLS"* || "$PROMPT" == *"PayIn"* ]]; then
    HINTS+="\n📌 Context hint: CLS work uses cls-settlement agent + check .base/flows/cls/ for canonical flows."
fi

if [[ "$PROMPT" == *"BACEN"* || "$PROMPT" == *"DEC"* || "$PROMPT" == *"IOF"* ]]; then
    HINTS+="\n📌 Context hint: BACEN compliance uses bacen-compliance agent + check .base/plans/09-compliance/."
fi

if [[ "$PROMPT" == *"ontology"* || "$PROMPT" == *"TTL"* ]]; then
    HINTS+="\n📌 Context hint: Ontology work uses ontology-shacl agent + check .base/aasc/ontology/."
fi

[[ -n "$HINTS" ]] && echo -e "$HINTS"
exit 0
