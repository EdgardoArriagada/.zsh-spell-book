# Bun scripts

Install dependencies:

```sh
cd src-bun
bun install --frozen-lockfile
```

Generate a self-contained HTML diagram from a file or stdin:

```sh
bun bin/mermaidjs examples/ads-chat.mmd ads-chat.html
cat examples/ads-chat.mmd | bun bin/mermaidjs - diagram.html
```

Open the HTML in a browser. Mermaid parses and renders there; invalid syntax shows an error.
No server or internet connection required. HTML input is escaped and Mermaid runs in strict mode.
Input must contain 1–50,000 characters. Existing output files are never overwritten.
Omit both arguments to read stdin and write `diagram.html`.

Run checks:

```sh
bun mermaid.test.js
```
