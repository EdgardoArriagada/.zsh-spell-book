hw() {
  ${zsb}.assertJira
  zparseopts -D -E -F -- p=parent || return 1
  local msg=${*:?Error: You must provide a message.}

  local token="$(pass jira)"
  [[ -z "$token" ]] && ${zsb}.throw "Token obtain failed"

  local response body httpStatus
  response=$(curl -s -w "\n%{http_code}" -u "${ZSB_JIRA_EMAIL}:${token}" \
   -H "Content-Type: application/json" "$ZSB_JIRA_BASEURL/rest/api/3/issue/$ZSB_CURRENT_TICKET/comment" \
   -X POST --data \
      "{
      \"body\": {
          \"type\": \"doc\",
          \"version\": 1,
          \"content\": [
            {
              \"type\": \"paragraph\",
              \"content\": [
                {
                  \"text\": \"${msg}\",
                  \"type\": \"text\"
                }
              ]
            }
          ]
        }
      }")
  body="${response%$'\n'*}"
  httpStatus="${response##*$'\n'}"
  echo "$httpStatus"
  [[ "$httpStatus" != 2* ]] && echo "$body"
}

compdef "_${zsb}.singleComp '-p'" hw
