import assert from "node:assert/strict";
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { createHtml } from "./bin/mermaidjs";

const dir = await mkdtemp(join(tmpdir(), "zsb-mermaid-"));
const script = new URL("./bin/mermaidjs", import.meta.url).pathname;
const source = await Bun.file(new URL("./examples/ads-chat.mmd", import.meta.url)).text();
const run = (args, stdin) => Bun.spawnSync([process.execPath, script, ...args], {
  cwd: dir, stdin: Buffer.from(stdin ?? ""), stdout: "pipe", stderr: "pipe",
});

try {
  await writeFile(join(dir, "input.mmd"), source);
  assert.equal(run(["input.mmd", "file.html"]).exitCode, 0);
  assert.equal(run([], source).exitCode, 0);
  assert.equal(run(["-", "pipe.html"], source).exitCode, 0);
  const html = await readFile(join(dir, "file.html"), "utf8");
  assert.equal(html, await readFile(join(dir, "pipe.html"), "utf8"));
  assert.equal(html, await readFile(join(dir, "diagram.html"), "utf8"));
  assert.ok(html.includes("data:text/javascript;base64,"));
  assert.ok(html.includes('securityLevel: "strict"'));
  assert.equal(run(["input.mmd", "file.html"]).exitCode, 1);
  assert.equal(await readFile(join(dir, "file.html"), "utf8"), html);
  assert.equal(run(["missing.mmd"]).exitCode, 1);
  assert.equal(run([], " ").exitCode, 1);
  assert.equal(run(["a", "b", "c"]).exitCode, 1);
  assert.equal(run(["--help"]).exitCode, 0);
  const escaped = await createHtml('flowchart LR\nA["</pre><script>alert(1)</script>&"]');
  assert.ok(escaped.includes("&lt;/pre&gt;&lt;script&gt;alert(1)&lt;/script&gt;&amp;"));
  await assert.rejects(createHtml("x".repeat(50_001)), /50,000/);
  console.log("Mermaid CLI checks passed.");
} finally {
  await rm(dir, { recursive: true, force: true });
}
