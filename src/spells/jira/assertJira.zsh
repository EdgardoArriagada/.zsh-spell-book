${zsb}.assertJira() {
  ${zsb}.assertIsSet 'ZSB_TICKETS_DIR'
  ${zsb}.assertIsSet 'ZSB_CURRENT_TICKET'
  ${zsb}.assertIsSet 'ZSB_PARENT_TICKET'
  ${zsb}.assertIsSet 'ZSB_JIRA_EMAIL'
  ${zsb}.assertIsSet 'ZSB_JIRA_BASEURL'
  ${zsb}.assertIsSet 'ZSB_CURRENT_LABEL'
}
