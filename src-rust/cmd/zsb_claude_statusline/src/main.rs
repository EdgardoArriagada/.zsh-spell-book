use std::io::{self, BufWriter, Read, Write};
use std::process::Command;

use serde::Deserialize;

// ── Color palette (Nordic — 256-color approximations) ─────────────────────────
const DARK: u8 = 236;    // #303030 ≈ #2E3440
const GOLDEN: u8 = 222;  // #FFD787 ≈ #EBCB8B
const BLUE: u8 = 110;    // #87AFD7 ≈ #81A1C1
const PURPLE: u8 = 140;  // #AF87D7 ≈ #B48EAD
const GREEN: u8 = 150;   // #AFD787 ≈ #A3BE8C
const ORANGE: u8 = 173;  // #D7875F ≈ #D08770
const RED: u8 = 131;     // #AF5F5F ≈ #BF616A

const SEP_RIGHT: &str = "\u{e0b0}"; // powerline filled right arrow
const RESET: &str = "\x1b[0m";

// Pre-computed ANSI escape sequences — avoids String allocation on every call.
fn fg(n: u8) -> &'static str {
    match n {
        DARK   => "\x1b[38;5;236m",
        GOLDEN => "\x1b[38;5;222m",
        BLUE   => "\x1b[38;5;110m",
        PURPLE => "\x1b[38;5;140m",
        GREEN  => "\x1b[38;5;150m",
        ORANGE => "\x1b[38;5;173m",
        RED    => "\x1b[38;5;131m",
        _      => "",
    }
}

fn bg(n: u8) -> &'static str {
    match n {
        DARK   => "\x1b[48;5;236m",
        GOLDEN => "\x1b[48;5;222m",
        BLUE   => "\x1b[48;5;110m",
        PURPLE => "\x1b[48;5;140m",
        GREEN  => "\x1b[48;5;150m",
        ORANGE => "\x1b[48;5;173m",
        RED    => "\x1b[48;5;131m",
        _      => "",
    }
}

// ── Input structs ──────────────────────────────────────────────────────────────
#[derive(Deserialize, Default)]
struct ModelInfo {
    display_name: Option<String>,
}

#[derive(Deserialize, Default)]
struct CurrentUsage {
    input_tokens: Option<u64>,
    output_tokens: Option<u64>,
    cache_creation_input_tokens: Option<u64>,
    cache_read_input_tokens: Option<u64>,
}

#[derive(Deserialize, Default)]
struct ContextWindow {
    used_percentage: Option<f64>,
    current_tokens: Option<u64>,
    max_tokens: Option<u64>,
    current_usage: Option<CurrentUsage>,
}

#[derive(Deserialize, Default)]
struct Cost {
    total_cost_usd: Option<f64>,
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

/// Single subprocess returning (branch_name, is_dirty).
/// Replaces the previous git_branch + git_is_dirty two-call pattern.
fn git_info(cwd: &str) -> (Option<String>, bool) {
    let Ok(output) = Command::new("git")
        .args(["-C", cwd, "status", "--porcelain", "--branch"])
        .output()
    else {
        return (None, false);
    };
    if !output.status.success() {
        return (None, false);
    }

    let stdout = String::from_utf8_lossy(&output.stdout);
    let mut lines = stdout.lines();

    // First line: "## main...origin/main", "## main", or "## HEAD (no branch)"
    let branch = lines.next().and_then(|header| {
        let name = header.strip_prefix("## ")?;
        let name = name.split("...").next().unwrap_or(name).trim();
        if name.is_empty() || name.starts_with("HEAD") {
            None
        } else {
            Some(name.to_string())
        }
    });

    // Any remaining lines mean the working tree is dirty.
    let dirty = lines.next().is_some();
    (branch, dirty)
}

// ── Formatting helpers ─────────────────────────────────────────────────────────

fn project_name(cwd: &str) -> &str {
    std::path::Path::new(cwd)
        .file_name()
        .and_then(|n| n.to_str())
        .unwrap_or(cwd)
}

fn capitalize(s: &str) -> String {
    let mut chars = s.chars();
    match chars.next() {
        None => String::new(),
        Some(first) => first.to_uppercase().collect::<String>() + chars.as_str(),
    }
}

/// "claude-sonnet-4-6" → "Sonnet 4", "claude-opus-4-6" → "Opus 4"
fn shorten_model(name: &str) -> String {
    let mut iter = name.split('-').skip(1); // skip "claude"
    let model_type = iter.next().map(capitalize);
    let version = iter.next();
    match (model_type, version) {
        (Some(t), Some(v)) => format!("{t} {v}"),
        (Some(t), None) => t,
        _ => name.to_string(),
    }
}

/// Resolve context usage as a percentage (0–100). Defaults to 0 when context
/// window is present but percentage/tokens are not yet populated.
fn resolve_pct(input: &StatusInput) -> Option<f64> {
    let ctx = input.context_window.as_ref()?;
    Some(ctx.used_percentage.or_else(|| {
        let used = ctx.current_tokens?;
        let max = ctx.max_tokens.filter(|&m| m > 0)?;
        Some(used as f64 / max as f64 * 100.0)
    }).unwrap_or(0.0))
}

fn bar_color(pct: f64) -> u8 {
    if pct == 0.0 {
        BLUE
    } else if pct < 35.0 {
        GREEN
    } else if pct < 70.0 {
        GOLDEN
    } else if pct < 85.0 {
        ORANGE
    } else {
        RED
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
    bg_color: u8,
    fg_color: u8,
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
            out.push_str(bg(seg.bg_color));
            out.push_str(fg(seg.fg_color));
        } else {
            let prev_bg = segments[i - 1].bg_color;
            out.push_str(bg(seg.bg_color));
            out.push_str(fg(prev_bg));
            out.push_str(SEP_RIGHT);
            out.push_str(fg(seg.fg_color));
        }
        out.push_str(&seg.content);
    }
    // Trailing separator (fg = last segment's bg, bg = terminal default)
    let last_bg = segments.last().unwrap().bg_color;
    out.push_str(RESET);
    out.push_str(fg(last_bg));
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

    // Segment 2: git branch with dirty indicator — single git subprocess.
    let cwd = input.cwd.as_deref().unwrap_or(".");
    let (git_branch, dirty) = git_info(cwd);

    let branch = input
        .worktree
        .as_ref()
        .and_then(|w| w.branch.as_deref())
        .map(str::to_string)
        .or(git_branch);

    if let Some(branch) = branch {
        let dirty_indicator = if dirty { " \u{00b1}" } else { "" };
        segments.push(Segment {
            bg_color: BLUE,
            fg_color: DARK,
            content: format!(" \u{e0a0} {branch}{dirty_indicator} "),
        });
    }

    render_line(&segments)
}

fn build_line2(input: &StatusInput) -> String {
    let mut segments = Vec::new();

    // Segment 0: progress bar on the LEFT (bordered box style)
    if let Some(pct) = resolve_pct(input) {
        segments.push(build_progress_segment(pct));
    }

    // Segment 1: model (PURPLE bg, DARK fg)
    let model_name = input
        .model
        .as_ref()
        .and_then(|m| m.display_name.as_deref())
        .map(shorten_model)
        .unwrap_or_else(|| "Claude".to_string());
    segments.push(Segment {
        bg_color: PURPLE,
        fg_color: DARK,
        content: format!(" \u{f444} {model_name} "),
    });

    // Segment 2: input · output · cache tokens (GREEN bg, DARK fg)
    let usage = input.context_window.as_ref().and_then(|c| c.current_usage.as_ref());
    let input_tok = usage.and_then(|u| u.input_tokens).unwrap_or(0);
    let output_tok = usage.and_then(|u| u.output_tokens).unwrap_or(0);
    let cache_tok = usage.and_then(|u| {
        match (u.cache_read_input_tokens, u.cache_creation_input_tokens) {
            (Some(r), Some(w)) => Some(r + w),
            (Some(r), None) => Some(r),
            (None, Some(w)) => Some(w),
            _ => None,
        }
    }).unwrap_or(0);

    segments.push(Segment {
        bg_color: GREEN,
        fg_color: DARK,
        content: format!(
            " \u{f019} {}  \u{f093} {}  \u{f1c0} {} ",
            format_k(input_tok),
            format_k(output_tok),
            format_k(cache_tok),
        ),
    });

    // Segment 3: total I+O · session cost (DARK bg, GOLDEN fg)
    let total = input_tok + output_tok;
    let cost_usd = input.cost.as_ref().and_then(|c| c.total_cost_usd).unwrap_or(0.0);

    {
        let mut parts: Vec<String> = Vec::new();
        parts.push(format!("\u{f080} {}", format_k(total)));   // fa-bar-chart = total
        parts.push(format!("\u{f0d6} ${cost_usd:.2}"));        // fa-money     = cost
        segments.push(Segment {
            bg_color: DARK,
            fg_color: GOLDEN,
            content: format!(" {} ", parts.join("  ")),
        });
    }

    render_line(&segments)
}

/// Left-side progress bar segment: block fill style, no box-drawing borders.
fn build_progress_segment(pct: f64) -> Segment {
    const BAR_WIDTH: usize = 10;
    let filled = ((pct / 100.0) * BAR_WIDTH as f64).round() as usize;
    let filled = filled.min(BAR_WIDTH);
    let color = bar_color(pct);
    let bars: String = (0..BAR_WIDTH)
        .map(|i| if i < filled { '\u{2588}' } else { '\u{2591}' })
        .collect();
    Segment {
        bg_color: DARK,
        fg_color: color,
        content: format!(" {bars} {pct:.0}% "),
    }
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

    let stdout = io::stdout();
    let mut out = BufWriter::new(stdout.lock());
    let _ = writeln!(out, "{line1}");
    let _ = write!(out, "{line2}");
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn bar_color_blue_at_zero() {
        assert_eq!(bar_color(0.0), BLUE);
    }

    #[test]
    fn bar_color_green_low_usage() {
        assert_eq!(bar_color(0.1), GREEN);
        assert_eq!(bar_color(20.0), GREEN);
        assert_eq!(bar_color(34.9), GREEN);
    }

    #[test]
    fn bar_color_golden_mid_usage() {
        assert_eq!(bar_color(35.0), GOLDEN);
        assert_eq!(bar_color(50.0), GOLDEN);
        assert_eq!(bar_color(69.9), GOLDEN);
    }

    #[test]
    fn bar_color_orange_high_usage() {
        assert_eq!(bar_color(70.0), ORANGE);
        assert_eq!(bar_color(80.0), ORANGE);
        assert_eq!(bar_color(84.9), ORANGE);
    }

    #[test]
    fn bar_color_red_critical() {
        assert_eq!(bar_color(85.0), RED);
        assert_eq!(bar_color(95.0), RED);
        assert_eq!(bar_color(100.0), RED);
    }
}
