hw() {
  ${zsb}.validateJira
  zparseopts -D -E -F -- p=parent || return 1
  local msg=${*?Error: You must provide a message.}

  local ticket
  [[ -n "$parent" ]] && ticket="$ZSB_PARENT_TICKET" || ticket="$ZSB_CURRENT_TICKET"

  local token="$(pass jira)"
  [[ -z "$token" ]] && ${zsb}.throw "Token obtain failed"

  curl -s -o /dev/null -w "%{http_code}" -u "${ZSB_JIRA_EMAIL}:${token}" \
   -H "Content-Type: application/json" "${ZSB_JIRA_BASEURL}/rest/api/3/issue/${ticket}/comment" \
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
      }"
}

compdef "_${zsb}.singleComp '-p'" hw
