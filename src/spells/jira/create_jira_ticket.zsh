create_jira_ticket() {
  if (( $# != 2 )); then
    ${zsb}.throw 'Usage: create_jira_ticket "<title>" "<description>"'
  fi

  local title="$1"
  local description="$2"

  [[ -n "$title" ]] || ${zsb}.throw "Title cannot be empty."
  [[ -n "$description" ]] || ${zsb}.throw "Description cannot be empty."

  ${zsb}.assertIsSet 'ZSB_JIRA_EMAIL'
  ${zsb}.assertIsSet 'ZSB_JIRA_BASEURL'
  ${zsb}.assertIsSet 'ZSB_JIRA_PROJECT_KEY'
  ${zsb}.assertIsSet 'ZSB_JIRA_ISSUE_TYPE_ID'
  ${zsb}.assertIsSet 'ZSB_JIRA_REPORTER_ACCOUNT_ID'
  ${zsb}.assertIsSet 'ZSB_JIRA_ASSIGNEE_ACCOUNT_ID'
  ${zsb}.assertIsSet 'ZSB_JIRA_PARENT_KEY'
  ${zsb}.assertIsSet 'ZSB_JIRA_PRIORITY_ID'
  ${zsb}.assertIsSet 'ZSB_JIRA_LABELS'
  ${zsb}.assertIsSet 'ZSB_JIRA_SPRINT_ID'

  command -v jq >/dev/null 2>&1 || ${zsb}.throw "You must install jq to create Jira payloads."

  local jira_baseurl="${ZSB_JIRA_BASEURL%/}"
  [[ "$jira_baseurl" == "https://mercadolibre.atlassian.net" ]] || \
    ${zsb}.throw "ZSB_JIRA_BASEURL must be https://mercadolibre.atlassian.net"

  local token="$(pass jira)"
  [[ -z "$token" ]] && ${zsb}.throw "Token obtain failed"

  local payload
  payload="$(jq -n \
    --arg project_key "$ZSB_JIRA_PROJECT_KEY" \
    --arg issue_type_id "$ZSB_JIRA_ISSUE_TYPE_ID" \
    --arg reporter_account_id "$ZSB_JIRA_REPORTER_ACCOUNT_ID" \
    --arg assignee_account_id "$ZSB_JIRA_ASSIGNEE_ACCOUNT_ID" \
    --arg parent_key "$ZSB_JIRA_PARENT_KEY" \
    --arg priority_id "$ZSB_JIRA_PRIORITY_ID" \
    --arg labels "$ZSB_JIRA_LABELS" \
    --arg sprint_id "$ZSB_JIRA_SPRINT_ID" \
    --arg title "$title" \
    --arg description "$description" \
    '
    def csv_strings($text):
      $text
      | split(",")
      | map(gsub("^\\s+|\\s+$"; ""))
      | map(select(. != ""));

    def adf($text):
      {
        type: "doc",
        version: 1,
        content: (
          $text
          | split("\n")
          | map({
              type: "paragraph",
              content: (if . == "" then [] else [{type: "text", text: .}] end)
            })
        )
      };

    {
      fields: {
        project: {key: $project_key},
        issuetype: {id: $issue_type_id},
        reporter: {accountId: $reporter_account_id},
        assignee: {accountId: $assignee_account_id},
        parent: {key: $parent_key},
        priority: {id: $priority_id},
        labels: csv_strings($labels),
        customfield_10115: ($sprint_id | gsub("^\\s+|\\s+$"; "") | tonumber),
        summary: $title,
        description: adf($description)
      }
    }')" || ${zsb}.throw "Could not build Jira payload. Check ZSB_JIRA_SPRINT_ID."

  local response
  response="$(curl -sS -w $'\n%{http_code}' -u "${ZSB_JIRA_EMAIL}:${token}" \
    -H "Content-Type: application/json" \
    "$jira_baseurl/rest/api/3/issue" \
    -X POST --data "$payload")"

  local http_code="${response##*$'\n'}"
  local body="${response%$'\n'$http_code}"

  if [[ "$http_code" != "201" ]]; then
    local jira_error="$(jq -r '
      if (.errorMessages // []) | length > 0 then
        .errorMessages | join("; ")
      elif (.errors // {}) | length > 0 then
        .errors | to_entries | map("\(.key): \(.value)") | join("; ")
      else
        empty
      end
    ' <<< "$body" 2>/dev/null)"

    [[ -n "$jira_error" ]] || jira_error="Unexpected Jira response."
    ${zsb}.throw "Jira ticket creation failed ($http_code): $jira_error"
  fi

  local key="$(jq -r '.key // empty' <<< "$body")"
  [[ -n "$key" ]] || ${zsb}.throw "Jira ticket creation succeeded but no key was returned."

  print -r -- "$jira_baseurl/browse/$key"
}

_${zsb}.nocompletion create_jira_ticket
