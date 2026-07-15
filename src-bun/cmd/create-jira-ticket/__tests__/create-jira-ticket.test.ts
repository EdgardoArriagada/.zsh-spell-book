import { describe, expect, test } from "bun:test";
import { Buffer } from "node:buffer";
import { chmod, mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const testDir = dirname(fileURLToPath(import.meta.url));
const commandPath = join(testDir, "..", "..", "..", "bin", "create-jira-ticket");

const defaultJiraEnv = {
  ZSB_JIRA_BASEURL: "https://mercadolibre.atlassian.net/",
  ZSB_JIRA_EMAIL: "dev@example.com",
  ZSB_JIRA_PROJECT_KEY: "ABC",
  ZSB_JIRA_ISSUE_TYPE_IDS: "Epic=10000,Task=10101",
  ZSB_JIRA_REPORTER_ACCOUNT_ID: "reporter-1",
  ZSB_JIRA_ASSIGNEE_ACCOUNT_ID: "assignee-1",
  ZSB_PARENT_TICKET: "ABC-1",
  ZSB_JIRA_PRIORITY_ID: "3",
  ZSB_JIRA_LABELS: " foo,bar,, baz ",
};

type RecordedRequest = {
  url: string;
  method: string;
  headers: Record<string, string>;
  body: any;
};

type RunOptions = {
  args?: string[];
  env?: Record<string, string | undefined>;
  fetchStatus?: number;
  fetchBody?: unknown;
  passOutput?: string;
  passExitCode?: number;
  passStderr?: string;
  fzfOutput?: string;
  fzfExitCode?: number;
};

async function writeFakePass(binDir: string): Promise<void> {
  const passPath = join(binDir, "pass");
  await writeFile(
    passPath,
    `#!/usr/bin/env sh
if [ "$1" != "jira" ]; then
  exit 64
fi
printf "%s" "\${ZSB_TEST_PASS_STDERR:-}" >&2
printf "%s" "\${ZSB_TEST_PASS_OUTPUT-secret-token}"
exit "\${ZSB_TEST_PASS_EXIT_CODE:-0}"
`
  );
  await chmod(passPath, 0o755);
}

async function writeFakeFzf(binDir: string): Promise<void> {
  const fzfPath = join(binDir, "fzf");
  await writeFile(
    fzfPath,
    `#!/usr/bin/env sh
cat > /dev/null
if [ "\${ZSB_TEST_FZF_EXIT_CODE:-0}" != "0" ]; then
  exit "\${ZSB_TEST_FZF_EXIT_CODE}"
fi
printf "%s\\n" "\${ZSB_TEST_FZF_OUTPUT:-Epic}"
exit 0
`
  );
  await chmod(fzfPath, 0o755);
}

async function writeFetchPreload(preloadPath: string): Promise<void> {
  await writeFile(
    preloadPath,
    `const requestsFile = process.env.ZSB_TEST_REQUESTS_FILE;

globalThis.fetch = async (url, init = {}) => {
  const headers = Object.fromEntries(new Headers(init.headers ?? {}).entries());
  const body = typeof init.body === "string" ? JSON.parse(init.body) : init.body ?? null;
  const request = {
    url: String(url),
    method: init.method ?? "GET",
    headers,
    body,
  };

  let requests = [];
  try {
    requests = JSON.parse(await Bun.file(requestsFile).text());
  } catch {
    requests = [];
  }
  requests.push(request);
  await Bun.write(requestsFile, JSON.stringify(requests, null, 2));

  return new Response(process.env.ZSB_TEST_FETCH_BODY ?? JSON.stringify({ key: "ABC-123" }), {
    status: Number(process.env.ZSB_TEST_FETCH_STATUS ?? "201"),
    headers: { "Content-Type": "application/json" },
  });
};
`
  );
}

async function runCreateJiraTicket(options: RunOptions = {}) {
  const tmp = await mkdtemp(join(tmpdir(), "zsb-create-jira-ticket-"));
  const binDir = join(tmp, "bin");
  const requestsFile = join(tmp, "requests.json");
  const preloadPath = join(tmp, "mock-fetch.ts");

  await mkdir(binDir);
  await writeFile(requestsFile, "[]");
  await writeFakePass(binDir);
  await writeFakeFzf(binDir);
  await writeFetchPreload(preloadPath);

  const env: Record<string, string> = {
    PATH: `${binDir}:${process.env.PATH ?? ""}`,
    HOME: tmp,
    TMPDIR: tmp,
    ZSB_TEST_REQUESTS_FILE: requestsFile,
    ZSB_TEST_FETCH_STATUS: String(options.fetchStatus ?? 201),
    ZSB_TEST_FETCH_BODY:
      options.fetchBody === undefined
        ? JSON.stringify({ key: "ABC-123" })
        : typeof options.fetchBody === "string"
          ? options.fetchBody
          : JSON.stringify(options.fetchBody),
    ZSB_TEST_PASS_OUTPUT: options.passOutput ?? "secret-token",
    ZSB_TEST_PASS_EXIT_CODE: String(options.passExitCode ?? 0),
    ZSB_TEST_PASS_STDERR: options.passStderr ?? "",
    ZSB_TEST_FZF_OUTPUT: options.fzfOutput ?? "Epic",
    ZSB_TEST_FZF_EXIT_CODE: String(options.fzfExitCode ?? 0),
    ...defaultJiraEnv,
  };

  for (const [name, value] of Object.entries(options.env ?? {})) {
    if (value === undefined) {
      delete env[name];
    } else {
      env[name] = value;
    }
  }

  try {
    const proc = Bun.spawn(
      [
        process.execPath,
        "--no-env-file",
        "--preload",
        preloadPath,
        commandPath,
        ...(options.args ?? ["Ticket title", "Line one\n\nLine three"]),
      ],
      {
        env,
        stdin: "ignore",
        stdout: "pipe",
        stderr: "pipe",
      }
    );

    const [stdout, stderr, exitCode] = await Promise.all([
      new Response(proc.stdout).text(),
      new Response(proc.stderr).text(),
      proc.exited,
    ]);
    const requests = JSON.parse(
      await readFile(requestsFile, "utf8")
    ) as RecordedRequest[];

    return { stdout, stderr, exitCode, requests };
  } finally {
    await rm(tmp, { recursive: true, force: true });
  }
}

describe("create-jira-ticket", () => {
  test("creates the issue with the expected payload and prints only the URL", async () => {
    const result = await runCreateJiraTicket();

    expect(result.exitCode).toBe(0);
    expect(result.stderr).toBe("");
    expect(result.stdout).toBe(
      "https://mercadolibre.atlassian.net/browse/ABC-123\n"
    );
    expect(result.requests).toHaveLength(1);
    expect(result.requests[0].url).toBe(
      "https://mercadolibre.atlassian.net/rest/api/3/issue"
    );
    expect(result.requests[0].method).toBe("POST");
    expect(result.requests[0].headers.authorization).toBe(
      `Basic ${Buffer.from("dev@example.com:secret-token").toString("base64")}`
    );
    expect(result.requests[0].body.fields.issuetype).toEqual({ id: "10000" });
    expect(result.requests[0].body.fields.labels).toEqual(["foo", "bar", "baz"]);
    expect(result.requests[0].body.fields.description).toEqual({
      type: "doc",
      version: 1,
      content: [
        {
          type: "paragraph",
          content: [{ type: "text", text: "Line one" }],
        },
        {
          type: "paragraph",
          content: [],
        },
        {
          type: "paragraph",
          content: [{ type: "text", text: "Line three" }],
        },
      ],
    });
    expect(
      result.requests.some((request) => request.url.includes("transitions"))
    ).toBe(false);
  });

  test("uses the issue type id matching the fzf-selected label", async () => {
    const result = await runCreateJiraTicket({ fzfOutput: "Task" });

    expect(result.exitCode).toBe(0);
    expect(result.requests[0].body.fields.issuetype).toEqual({ id: "10101" });
  });

  test("fails when fzf selection is cancelled", async () => {
    const result = await runCreateJiraTicket({ fzfExitCode: 130 });

    expect(result.exitCode).not.toBe(0);
    expect(result.stdout).toBe("");
    expect(result.stderr).toContain("Issue type selection cancelled.");
    expect(result.requests).toHaveLength(0);
  });

  test("rejects the wrong argument count", async () => {
    const result = await runCreateJiraTicket({ args: [] });

    expect(result.exitCode).not.toBe(0);
    expect(result.stdout).toBe("");
    expect(result.stderr).toContain(
      'Usage: create-jira-ticket "<title>" ["<description>"]'
    );
    expect(result.requests).toHaveLength(0);
  });

  test("rejects missing required environment variables", async () => {
    const result = await runCreateJiraTicket({
      env: { ZSB_JIRA_EMAIL: undefined },
    });

    expect(result.exitCode).not.toBe(0);
    expect(result.stdout).toBe("");
    expect(result.stderr).toContain("You must set ZSB_JIRA_EMAIL first.");
    expect(result.requests).toHaveLength(0);
  });

  test("rejects missing ZSB_JIRA_ISSUE_TYPE_IDS", async () => {
    const result = await runCreateJiraTicket({
      env: { ZSB_JIRA_ISSUE_TYPE_IDS: undefined },
    });

    expect(result.exitCode).not.toBe(0);
    expect(result.stderr).toContain("You must set ZSB_JIRA_ISSUE_TYPE_IDS first.");
    expect(result.requests).toHaveLength(0);
  });

  test("rejects an empty token from pass", async () => {
    const result = await runCreateJiraTicket({ passOutput: "" });

    expect(result.exitCode).not.toBe(0);
    expect(result.stdout).toBe("");
    expect(result.stderr).toContain("Token obtain failed.");
    expect(result.requests).toHaveLength(0);
  });

  test("reports parsed Jira field errors", async () => {
    const result = await runCreateJiraTicket({
      fetchStatus: 400,
      fetchBody: {
        errors: {
          summary: "Summary is required.",
          customfield_10115: "Sprint id is invalid.",
        },
      },
    });

    expect(result.exitCode).not.toBe(0);
    expect(result.stdout).toBe("");
    expect(result.stderr).toContain(
      "Jira ticket creation failed (400): summary: Summary is required.; customfield_10115: Sprint id is invalid."
    );
    expect(result.requests).toHaveLength(1);
  });

  test("fails when Jira creation succeeds without a valid issue key", async () => {
    const result = await runCreateJiraTicket({
      fetchBody: { id: "10000" },
    });

    expect(result.exitCode).not.toBe(0);
    expect(result.stdout).toBe("");
    expect(result.stderr).toContain(
      "Jira ticket creation succeeded but no valid key was returned."
    );
    expect(result.requests).toHaveLength(1);
  });
});
