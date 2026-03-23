use std::io::{self, Read};
use std::process::Command;

use serde::Deserialize;

// ── Color palette (Nordic — AlexvZyl/nordic.nvim) ─────────────────────────────
const DARK: (u8, u8, u8) = (46, 52, 64);       // gray1      #2E3440
const GOLDEN: (u8, u8, u8) = (235, 203, 139);  // yellow.base #EBCB8B
const BLUE: (u8, u8, u8) = (129, 161, 193);    // blue1       #81A1C1
const GREEN: (u8, u8, u8) = (163, 190, 140);   // green.base  #A3BE8C

const SEP_RIGHT: &str = "\u{e0b0}"; // powerline filled right arrow ("")
const RESET: &str = "\x1b[0m";

fn fg((r, g, b): (u8, u8, u8)) -> String {
    format!("\x1b[38;2;{r};{g};{b}m")
}

fn bg((r, g, b): (u8, u8, u8)) -> String {
    format!("\x1b[48;2;{r};{g};{b}m")
}

// ── Input structs ──────────────────────────────────────────────────────────────
#[derive(Deserialize, Default)]
struct ModelInfo {
    display_name: Option<String>,
}

#[derive(Deserialize, Default)]
struct ContextWindow {
    used_percentage: Option<f64>,
    current_tokens: Option<u64>,
    max_tokens: Option<u64>,
}

#[derive(Deserialize, Default)]
struct Cost {
    session_cost_usd: Option<f64>,
    #[allow(dead_code)]
    total_lines_added: Option<i64>,
    #[allow(dead_code)]
    total_lines_removed: Option<i64>,
}

#[derive(Deserialize, Default)]
struct WorktreeInfo {
    branch: Option<String>,
}

#[derive(Deserialize, Default)]
struct StatusInput {
    cwd: Option<String>,
    model: Option<ModelInfo>,
    context_window: Option<ContextWindow>,
    cost: Option<Cost>,
    worktree: Option<WorktreeInfo>,
}

// ── Git helpers ────────────────────────────────────────────────────────────────

fn git_branch(cwd: &str) -> Option<String> {
    Command::new("git")
        .args(["-C", cwd, "rev-parse", "--abbrev-ref", "HEAD"])
        .output()
        .ok()
        .filter(|o| o.status.success())
        .map(|o| String::from_utf8_lossy(&o.stdout).trim().to_string())
        .filter(|s| !s.is_empty() && s != "HEAD")
}

fn git_is_dirty(cwd: &str) -> bool {
    Command::new("git")
        .args(["-C", cwd, "status", "--porcelain"])
        .output()
        .ok()
        .filter(|o| o.status.success())
        .map(|o| !o.stdout.is_empty())
        .unwrap_or(false)
}

// ── Formatting helpers ─────────────────────────────────────────────────────────

fn project_name(cwd: &str) -> &str {
    std::path::Path::new(cwd)
        .file_name()
        .and_then(|n| n.to_str())
        .unwrap_or(cwd)
}

/// "claude-sonnet-4-6" → "Sonnet 4", "claude-opus-4-6" → "Opus 4"
fn shorten_model(name: &str) -> String {
    let parts: Vec<&str> = name.split('-').collect();
    let model_type = parts.get(1).map(|p| {
        let mut c = p.chars();
        match c.next() {
            None => p.to_string(),
            Some(first) => first.to_uppercase().collect::<String>() + c.as_str(),
        }
    });
    match (model_type, parts.get(2)) {
        (Some(t), Some(v)) => format!("{t} {v}"),
        (Some(t), None) => t,
        _ => name.to_string(),
    }
}

fn format_k(n: u64) -> String {
    if n >= 1000 {
        format!("{}k", n / 1000)
    } else {
        n.to_string()
    }
}

// ── Segment rendering ──────────────────────────────────────────────────────────

struct Segment {
    bg_color: (u8, u8, u8),
    fg_color: (u8, u8, u8),
    content: String,
}

/// Renders powerline-chained segments left-to-right, closing with a trailing separator.
fn render_line(segments: &[Segment]) -> String {
    if segments.is_empty() {
        return String::new();
    }
    let mut out = String::new();
    for (i, seg) in segments.iter().enumerate() {
        if i == 0 {
            out.push_str(&bg(seg.bg_color));
            out.push_str(&fg(seg.fg_color));
        } else {
            let prev_bg = segments[i - 1].bg_color;
            out.push_str(&bg(seg.bg_color));
            out.push_str(&fg(prev_bg));
            out.push_str(SEP_RIGHT);
            out.push_str(&fg(seg.fg_color));
        }
        out.push_str(&seg.content);
    }
    // Trailing separator (fg = last segment's bg, bg = terminal default)
    let last_bg = segments.last().unwrap().bg_color;
    out.push_str(RESET);
    out.push_str(&fg(last_bg));
    out.push_str(SEP_RIGHT);
    out.push_str(RESET);
    out
}

// ── Line builders ──────────────────────────────────────────────────────────────

fn build_line1(input: &StatusInput) -> String {
    let mut segments = Vec::new();

    // Segment 1: project name (dark bg, golden fg)
    let project = input.cwd.as_deref().map(project_name).unwrap_or("?");
    segments.push(Segment {
        bg_color: DARK,
        fg_color: GOLDEN,
        content: format!(" \u{e5fc} {project} "),
    });

    // Segment 2: git branch with dirty indicator (coral bg, dark fg)
    let cwd = input.cwd.as_deref().unwrap_or(".");
    let branch = input
        .worktree
        .as_ref()
        .and_then(|w| w.branch.as_deref())
        .map(str::to_string)
        .or_else(|| git_branch(cwd));

    if let Some(branch) = branch {
        let dirty = if git_is_dirty(cwd) { " \u{00b1}" } else { "" };
        segments.push(Segment {
            bg_color: BLUE,
            fg_color: DARK,
            content: format!(" \u{e0a0} {branch}{dirty} "),
        });
    }

    render_line(&segments)
}

fn build_line2(input: &StatusInput) -> String {
    let mut segments = Vec::new();

    // Segment 1: model name (coral bg, dark fg)
    let model_name = input
        .model
        .as_ref()
        .and_then(|m| m.display_name.as_deref())
        .map(shorten_model)
        .unwrap_or_else(|| "Claude".to_string());
    segments.push(Segment {
        bg_color: BLUE,
        fg_color: DARK,
        content: format!(" \u{f444} {model_name} "),
    });

    // Segment 2: context tokens + cost (lavender bg, dark fg)
    let ctx = input.context_window.as_ref();
    let cost = input.cost.as_ref();

    // Resolve token counts: prefer explicit fields, fall back to percentage estimate
    let (used_tokens, max_tokens) = ctx
        .and_then(|c| c.current_tokens.zip(c.max_tokens))
        .or_else(|| {
            ctx.and_then(|c| c.used_percentage).map(|pct| {
                let max: u64 = 200_000;
                let used = (pct / 100.0 * max as f64) as u64;
                (used, max)
            })
        })
        .unwrap_or((0, 0));

    let cost_usd = cost.and_then(|c| c.session_cost_usd);

    if used_tokens > 0 || cost_usd.is_some() {
        let mut info = format_k(used_tokens);
        if let Some(usd) = cost_usd {
            info.push_str(&format!(" \u{00b7} ${usd:.2}"));
        }
        if max_tokens > used_tokens {
            info.push_str(&format!(" {} left", format_k(max_tokens - used_tokens)));
        }
        segments.push(Segment {
            bg_color: GREEN,
            fg_color: DARK,
            content: format!(" \u{f0d6} {info} "),
        });
    }

    render_line(&segments)
}

// ── Main ───────────────────────────────────────────────────────────────────────

fn main() {
    let mut raw = String::new();
    if io::stdin().read_to_string(&mut raw).is_err() {
        return;
    }

    let input: StatusInput = serde_json::from_str(&raw).unwrap_or_default();

    let line1 = build_line1(&input);
    let line2 = build_line2(&input);

    println!("{line1}");
    print!("{line2}");
}
