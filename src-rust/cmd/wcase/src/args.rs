use clap::Parser;

/// word case
#[derive(Debug, Parser)]
#[command(author, about, long_about = None)]
pub struct Args {
    /// transform to flat case
    #[clap(long)]
    pub flat: bool,

    /// transform to upper case
    #[clap(long)]
    pub upper: bool,

    /// transform to camel case
    #[clap(long)]
    pub camel: bool,

    /// transform to pascal case
    #[clap(long)]
    pub pascal: bool,

    /// transform to snake case
    #[clap(long)]
    pub snake: bool,

    /// transform to all caps
    #[clap(long)]
    pub all_caps: bool,

    /// transform to kebab case
    #[clap(long)]
    pub kebab: bool,

    /// transform to train case
    #[clap(long)]
    pub train: bool,

    /// transform to spaced case
    #[clap(long)]
    pub spaced: bool,

    // transform to http header case
    #[clap(long)]
    pub http_header: bool,

    /// transform to title case
    #[clap(long)]
    pub title: bool,

    /// word to transform
    #[clap(long, short)]
    pub word: String,
}
