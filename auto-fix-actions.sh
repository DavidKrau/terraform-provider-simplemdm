#!/usr/bin/env bash
set -euo pipefail

MAX_ITERATIONS=5
LOG_DIR=".auto-fix-logs"
MODEL_NAME="gpt-4o-mini"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: required command '$1' not found in PATH" >&2
    exit 1
  fi
}

require_cmd git
require_cmd gh
require_cmd jq
require_cmd openai

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

mkdir -p "$LOG_DIR"

default_branch() {
  if gh repo view --json defaultBranchRef -q '.defaultBranchRef.name' >/dev/null 2>&1; then
    gh repo view --json defaultBranchRef -q '.defaultBranchRef.name'
  else
    git symbolic-ref --short refs/remotes/origin/HEAD 2>/dev/null | sed 's@^origin/@@'
  fi
}

NAME_WITH_OWNER=$(gh repo view --json nameWithOwner -q '.nameWithOwner')
if [[ -z "$NAME_WITH_OWNER" ]]; then
  origin_url=$(git remote get-url origin)
  NAME_WITH_OWNER=$(basename -s .git "$origin_url")
  owner_part=$(dirname "$origin_url" | sed 's#.*github.com[:/]##')
  NAME_WITH_OWNER="$owner_part/$NAME_WITH_OWNER"
fi

BASE_BRANCH=$(default_branch)
if [[ -n "$BASE_BRANCH" ]]; then
  git fetch origin "$BASE_BRANCH" >/dev/null 2>&1 || true
  if git show-ref --verify --quiet "refs/heads/$BASE_BRANCH"; then
    git checkout "$BASE_BRANCH"
    git pull --ff-only origin "$BASE_BRANCH" || git pull origin "$BASE_BRANCH"
  else
    git checkout -b "$BASE_BRANCH" "origin/$BASE_BRANCH"
  fi
fi

declare -A SUMMARY

gather_failing_runs() {
  gh api "repos/$NAME_WITH_OWNER/actions/runs" \
    -f status=completed \
    -f conclusion=failure \
    -f per_page=50 \
    --paginate \
    -q '.workflow_runs[] | {id: .id, name: .name, run_number: .run_number, head_branch: .head_branch, event: .event, html_url: .html_url}'
}

get_run_details() {
  local run_id="$1"
  gh api "repos/$NAME_WITH_OWNER/actions/runs/$run_id" \
    -q '{id: .id, name: .name, run_attempt: .run_attempt, head_sha: .head_sha, event: .event, display_title: .display_title, workflow_path: .path, head_branch: .head_branch, html_url: .html_url}'
}

call_openai_for_fix() {
  local prompt_file="$1"
  local jq_payload
  jq_payload=$(jq -n \
    --arg system "You are an autonomous software maintenance bot. Only respond with strict JSON containing keys 'analysis' and 'diff'. The 'diff' must be a unified diff that can be applied with 'git apply --index'. Use an empty string for 'diff' if no changes are required." \
    --arg user "$(<"$prompt_file")" \
    '{messages: [{role: "system", content: $system}, {role: "user", content: $user}]}'
  )

  openai api chat.completions.create \
    -m "$MODEL_NAME" \
    -g "$jq_payload"
}

apply_generated_diff() {
  local diff="$1"
  if [[ -z "$diff" || "$diff" == "null" ]]; then
    echo "No diff returned from OpenAI." >&2
    return 1
  fi
  echo "Applying diff from OpenAI..."
  if ! printf '%s\n' "$diff" | git apply --index; then
    echo "Failed to apply diff with --index, attempting without staging..." >&2
    if ! printf '%s\n' "$diff" | git apply; then
      echo "Failed to apply diff." >&2
      return 1
    fi
    git add -u
  fi
  return 0
}

wait_for_run() {
  local run_id="$1"
  if gh run watch "$run_id" --exit-status; then
    return 0
  fi
  return 1
}

checkout_run_branch() {
  local head_branch="$1"

  if [[ -z "$head_branch" || "$head_branch" == "null" ]]; then
    return 1
  fi

  git fetch origin "$head_branch" >/dev/null 2>&1 || true

  git reset --hard >/dev/null 2>&1 || true
  git clean -fd >/dev/null 2>&1 || true

  if git show-ref --verify --quiet "refs/heads/$head_branch"; then
    git checkout "$head_branch"
  else
    git checkout -b "$head_branch" "origin/$head_branch"
  fi

  if git rev-parse --verify "origin/$head_branch" >/dev/null 2>&1; then
    git reset --hard "origin/$head_branch" >/dev/null 2>&1 || true
  fi

  git clean -fd >/dev/null 2>&1 || true

  return 0
}

find_run_for_commit() {
  local workflow_path="$1"
  local head_branch="$2"
  local head_sha="$3"

  local attempt=0
  local max_attempts=30
  local runs_json
  local new_run_id

  while (( attempt < max_attempts )); do
    runs_json=$(gh api "repos/$NAME_WITH_OWNER/actions/runs" -f branch="$head_branch" -f per_page=50)
    if [[ -n "$runs_json" ]]; then
      new_run_id=$(jq -r --arg path "$workflow_path" --arg sha "$head_sha" \
        '.workflow_runs[] | select(.path == $path and .head_sha == $sha) | .id' <<<"$runs_json" | head -n1)
      if [[ -n "$new_run_id" && "$new_run_id" != "null" ]]; then
        echo "$new_run_id"
        return 0
      fi
    fi
    sleep 10
    ((attempt++))
  done

  return 1
}

ITERATION=1
while (( ITERATION <= MAX_ITERATIONS )); do
  echo "Iteration $ITERATION/$MAX_ITERATIONS" >&2
  mapfile -t RUNS < <(gather_failing_runs | jq -rc)
  if (( ${#RUNS[@]} == 0 )); then
    echo "No failing workflows detected." >&2
    break
  fi

  ANY_CHANGE=0
  for run_json in "${RUNS[@]}"; do
    run_id=$(jq -r '.id' <<<"$run_json")
    run_name=$(jq -r '.name' <<<"$run_json")
    run_number=$(jq -r '.run_number' <<<"$run_json")
    echo "Processing run $run_name (#$run_number, id $run_id)" >&2

    detail_json=$(get_run_details "$run_id")
    display_title=$(jq -r '.display_title' <<<"$detail_json")
    workflow_path=$(jq -r '.workflow_path' <<<"$detail_json")
    head_branch=$(jq -r '.head_branch' <<<"$detail_json")

    if [[ -z "$head_branch" || "$head_branch" == "null" ]]; then
      if [[ -n "$BASE_BRANCH" ]]; then
        head_branch="$BASE_BRANCH"
      else
        echo "Unable to determine branch for run $run_id" >&2
        SUMMARY["$run_name"]="Unknown branch"
        continue
      fi
    fi

    if ! checkout_run_branch "$head_branch"; then
      echo "Failed to checkout branch $head_branch for run $run_id" >&2
      SUMMARY["$run_name"]="Failed to checkout branch"
      continue
    fi

    log_file="$LOG_DIR/run-${run_id}.log"
    gh run view "$run_id" --log >"$log_file"

    prompt_file=$(mktemp)
    {
      echo "Repository: $NAME_WITH_OWNER"
      echo "Run ID: $run_id"
      echo "Workflow: $run_name"
      echo "Display Title: $display_title"
      echo "Workflow File: $workflow_path"
      echo "Git Status:"
      git status -sb
      echo
      echo "Recent log excerpt:";
      tail -n 400 "$log_file"
    } >"$prompt_file"

    response=$(call_openai_for_fix "$prompt_file")
    rm -f "$prompt_file"

    content=$(echo "$response" | jq -r '.choices[0].message.content')
    if [[ -z "$content" || "$content" == "null" ]]; then
      echo "OpenAI returned empty response" >&2
      SUMMARY["$run_name"]="No response from OpenAI"
      continue
    fi

    analysis=$(printf '%s\n' "$content" | jq -r '.analysis // ""') || analysis="$content"
    diff=$(printf '%s\n' "$content" | jq -r '.diff // ""') || diff=""

    echo "Analysis from OpenAI:" >&2
    printf '%s\n' "$analysis" >&2

    if apply_generated_diff "$diff"; then
      if git status --porcelain | grep -q .; then
        ANY_CHANGE=1
        commit_msg="Auto-fix GitHub Actions failure: $display_title"
        git commit -am "$commit_msg"
      fi

      git push origin "$head_branch"

      new_sha=$(git rev-parse HEAD)
      echo "Waiting for new workflow run on branch $head_branch to pick up commit $new_sha" >&2

      if new_run_id=$(find_run_for_commit "$workflow_path" "$head_branch" "$new_sha"); then
        echo "Monitoring new run $new_run_id for workflow $run_name" >&2
        if wait_for_run "$new_run_id"; then
          SUMMARY["$run_name"]="Fixed on iteration $ITERATION (run $new_run_id)"
        else
          SUMMARY["$run_name"]="Still failing after iteration $ITERATION (run $new_run_id)"
        fi
      else
        SUMMARY["$run_name"]="No new run detected for commit $new_sha"
      fi
    else
      SUMMARY["$run_name"]="Failed to apply diff"
    fi
  done

  if (( ANY_CHANGE == 0 )); then
    echo "No changes were applied in this iteration." >&2
  fi

  ((ITERATION++))
done

if (( ITERATION > MAX_ITERATIONS )); then
  echo "Reached iteration limit ($MAX_ITERATIONS)." >&2
fi

echo "Summary of workflow fixes:"
for key in "${!SUMMARY[@]}"; do
  printf ' - %s: %s\n' "$key" "${SUMMARY[$key]}"
done
