import { Buffer } from "node:buffer";

const REQUIRED_ENV = [
  "ZSB_JIRA_BASEURL",
  "ZSB_JIRA_EMAIL",
  "ZSB_JIRA_PROJECT_KEY",
  "ZSB_JIRA_ISSUE_TYPE_IDS",
  "ZSB_JIRA_REPORTER_ACCOUNT_ID",
  "ZSB_JIRA_ASSIGNEE_ACCOUNT_ID",
  "ZSB_PARENT_TICKET",
  "ZSB_JIRA_PRIORITY_ID",
  "ZSB_JIRA_LABELS",
] as const;

const ALLOWED_JIRA_BASE_URL = "https://mercadolibre.atlassian.net";
const JIRA_KEY_RE = /^[A-Z][A-Z0-9]+-\d+$/;

type EnvKey = (typeof REQUIRED_ENV)[number];
type JiraConfig = Record<EnvKey, string> & {
  jiraBaseUrl: typeof ALLOWED_JIRA_BASE_URL;
  labels: string[];
  issueTypeMap: Map<string, string>;
};

type AdfDoc = {
  type: "doc";
  version: 1;
  content: Array<{
    type: "paragraph";
    content: Array<{ type: "text"; text: string }>;
  }>;
};

class CliError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "CliError";
  }
}

function usage(): string {
  return 'Usage: create-jira-ticket "<title>" ["<description>"]';
}

function requireEnv(name: EnvKey): string {
  const value = process.env[name];
  if (value === undefined || value.length === 0) {
    throw new CliError(`You must set ${name} first.`);
  }
  return value;
}

function parseLabels(value: string): string[] {
  return value
    .split(",")
    .map((label) => label.trim())
    .filter((label) => label.length > 0);
}

function parseIssueTypeIds(value: string): Map<string, string> {
  const map = new Map<string, string>();
  for (const entry of value.split(",")) {
    const eq = entry.indexOf("=");
    if (eq < 1) continue;
    const label = entry.slice(0, eq).trim();
    const id = entry.slice(eq + 1).trim();
    if (label && id) map.set(label, id);
  }
  if (map.size === 0) {
    throw new CliError("ZSB_JIRA_ISSUE_TYPE_IDS has no valid entries.");
  }
  return map;
}

function loadConfig(): JiraConfig {
  const env = Object.fromEntries(
    REQUIRED_ENV.map((name) => [name, requireEnv(name)])
  ) as Record<EnvKey, string>;

  const jiraBaseUrl = env.ZSB_JIRA_BASEURL.replace(/\/+$/, "");
  if (jiraBaseUrl !== ALLOWED_JIRA_BASE_URL) {
    throw new CliError(
      `ZSB_JIRA_BASEURL must be ${ALLOWED_JIRA_BASE_URL}`
    );
  }

  return {
    ...env,
    jiraBaseUrl,
    labels: parseLabels(env.ZSB_JIRA_LABELS),
    issueTypeMap: parseIssueTypeIds(env.ZSB_JIRA_ISSUE_TYPE_IDS),
  };
}

async function readJiraToken(): Promise<string> {
  const proc = Bun.spawn(["pass", "jira"], {
    stdin: "ignore",
    stdout: "pipe",
    stderr: "pipe",
  });

  const [stdout, stderr, exitCode] = await Promise.all([
    new Response(proc.stdout).text(),
    new Response(proc.stderr).text(),
    proc.exited,
  ]);

  const token = stdout.trim();
  if (exitCode !== 0) {
    const message = stderr.trim() || "pass jira exited with a non-zero status.";
    throw new CliError(`Token obtain failed: ${message}`);
  }
  if (token.length === 0) {
    throw new CliError("Token obtain failed.");
  }

  return token;
}

async function selectIssueTypeId(issueTypeMap: Map<string, string>): Promise<string> {
  const labels = Array.from(issueTypeMap.keys()).join("\n");
  const proc = Bun.spawn(["fzf", "--prompt=Issue type: ", "--no-sort"], {
    stdin: Buffer.from(labels),
    stdout: "pipe",
    stderr: "inherit",
  });

  const [stdout, exitCode] = await Promise.all([
    new Response(proc.stdout).text(),
    proc.exited,
  ]);

  if (exitCode !== 0) {
    throw new CliError("Issue type selection cancelled.");
  }

  const selected = stdout.trim();
  const id = issueTypeMap.get(selected);
  if (!id) {
    throw new CliError(`Unknown issue type: ${selected}`);
  }
  return id;
}

function buildAdfDescription(description: string): AdfDoc {
  return {
    type: "doc",
    version: 1,
    content: description.split("\n").map((line) => ({
      type: "paragraph",
      content:
        line.length === 0
          ? []
          : [
              {
                type: "text",
                text: line,
              },
            ],
    })),
  };
}

function buildPayload(config: JiraConfig, title: string, description: string, issueTypeId: string) {
  return {
    fields: {
      project: { key: config.ZSB_JIRA_PROJECT_KEY },
      issuetype: { id: issueTypeId },
      reporter: { accountId: config.ZSB_JIRA_REPORTER_ACCOUNT_ID },
      assignee: { accountId: config.ZSB_JIRA_ASSIGNEE_ACCOUNT_ID },
      parent: { key: config.ZSB_PARENT_TICKET },
      priority: { id: config.ZSB_JIRA_PRIORITY_ID },
      labels: config.labels,
      summary: title,
      description: buildAdfDescription(description),
    },
  };
}

function parseJsonObject(body: string): Record<string, unknown> | null {
  try {
    const parsed = JSON.parse(body);
    if (parsed && typeof parsed === "object" && !Array.isArray(parsed)) {
      return parsed as Record<string, unknown>;
    }
  } catch {
    return null;
  }

  return null;
}

function parseJiraError(body: string): string {
  const parsed = parseJsonObject(body);
  if (!parsed) {
    return "Unexpected Jira response.";
  }

  const errorMessages = parsed.errorMessages;
  if (Array.isArray(errorMessages) && errorMessages.length > 0) {
    return errorMessages.map((message) => String(message)).join("; ");
  }

  const errors = parsed.errors;
  if (errors && typeof errors === "object" && !Array.isArray(errors)) {
    const messages = Object.entries(errors).map(
      ([field, message]) => `${field}: ${String(message)}`
    );
    if (messages.length > 0) {
      return messages.join("; ");
    }
  }

  return "Unexpected Jira response.";
}

function extractIssueKey(body: string): string {
  const parsed = parseJsonObject(body);
  const key = typeof parsed?.key === "string" ? parsed.key : "";
  if (!JIRA_KEY_RE.test(key)) {
    throw new CliError(
      "Jira ticket creation succeeded but no valid key was returned."
    );
  }

  return key;
}

async function createJiraTicket(
  config: JiraConfig,
  token: string,
  title: string,
  description: string,
  issueTypeId: string
): Promise<string> {
  const auth = Buffer.from(`${config.ZSB_JIRA_EMAIL}:${token}`, "utf8").toString(
    "base64"
  );

  const response = await fetch(`${config.jiraBaseUrl}/rest/api/3/issue`, {
    method: "POST",
    headers: {
      Authorization: `Basic ${auth}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(buildPayload(config, title, description, issueTypeId)),
  });

  const body = await response.text();
  if (response.status !== 201) {
    throw new CliError(
      `Jira ticket creation failed (${response.status}): ${parseJiraError(body)}`
    );
  }

  return extractIssueKey(body);
}

async function main(): Promise<void> {
  const args = process.argv.slice(2);
  if (args.length < 1 || args.length > 2) {
    throw new CliError(usage());
  }

  const [title, description = ""] = args;
  if (title.length === 0) {
    throw new CliError("Title cannot be empty.");
  }

  const config = loadConfig();
  const issueTypeId = await selectIssueTypeId(config.issueTypeMap);
  const token = await readJiraToken();
  const key = await createJiraTicket(config, token, title, description, issueTypeId);

  console.log(`${config.jiraBaseUrl}/browse/${key}`);
}

main().catch((error) => {
  if (error instanceof CliError) {
    console.error(error.message);
  } else {
    const message = error instanceof Error ? error.message : String(error);
    console.error(`Unexpected error: ${message}`);
  }
  process.exit(1);
});
