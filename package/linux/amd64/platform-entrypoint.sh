#!/usr/bin/env bash
set -euo pipefail

/usr/bin/update-platform-ca

metadata_url="${PLATFORM_METADATA_URL:-http://metadata/2016-07-29}"
metadata_answer="${PLATFORM_METADATA_ANSWER:-169.254.169.250}"
never_recurse_to="${NEVER_RECURSE_TO:-169.254.169.250}"
answers_file="${PLATFORM_DNS_ANSWERS_FILE:-/etc/internal-dns/answers.json}"
agent_ip=""

load_agent_ip() {
  local attempt
  for attempt in $(seq 1 60); do
    agent_ip="$(curl --fail --silent --show-error --max-time 2 "$metadata_url/self/host/agent_ip" || true)"
    if [[ -n "$agent_ip" && "$agent_ip" != "Not found" ]]; then
      return 0
    fi
    sleep 1
  done
  echo "metadata agent IP was unavailable after 60 attempts" >&2
  return 1
}

if [[ "$metadata_answer" == "agent_ip" || "$never_recurse_to" == "agent_ip" ]]; then
  load_agent_ip
fi
if [[ "$metadata_answer" == "agent_ip" ]]; then
  metadata_answer="$agent_ip"
fi
if [[ "$never_recurse_to" == "agent_ip" ]]; then
  never_recurse_to="$agent_ip"
fi

metadata_mode=()
if [[ "${PLATFORM_METADATA_ENABLED:-false}" == "true" ]]; then
  metadata_mode=(--metadata-url="$metadata_url")
fi

answers_mode=(--answers="$answers_file")
for argument in "$@"; do
  if [[ "$argument" == "--answers" || "$argument" == --answers=* ]]; then
    answers_mode=()
    break
  fi
done

exec "$@" \
  "${metadata_mode[@]}" \
  "${answers_mode[@]}" \
  --metadata-answer="$metadata_answer" \
  --never-recurse-to="$never_recurse_to"
