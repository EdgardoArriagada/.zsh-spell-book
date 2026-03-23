mod args;
use args::Args;
use clap::Parser;
use std::{fmt, process};

#[derive(Debug)]
enum Case {
    Flat,
    Upper,
    Camel,
    Pascal,
    Snake,
    AllCaps,
    Kebab,
    Train,
    Spaced,
    Title,
    HttpHeader,
    None,
}

impl fmt::Display for Case {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Case::Flat => write!(f, "flat"),
            Case::Upper => write!(f, "upper"),
            Case::Camel => write!(f, "camel"),
            Case::Pascal => write!(f, "pascal"),
            Case::Snake => write!(f, "snake"),
            Case::AllCaps => write!(f, "all_caps"),
            Case::Kebab => write!(f, "kebab"),
            Case::Train => write!(f, "train"),
            Case::Spaced => write!(f, "spaced"),
            Case::Title => write!(f, "title"),
            Case::HttpHeader => write!(f, "http_header"),
            Case::None => write!(f, "none"),
        }
    }
}

impl PartialEq for Case {
    fn eq(&self, other: &Self) -> bool {
        match (self, other) {
            (Case::Flat, Case::Flat) => true,
            (Case::Upper, Case::Upper) => true,
            (Case::Camel, Case::Camel) => true,
            (Case::Pascal, Case::Pascal) => true,
            (Case::Snake, Case::Snake) => true,
            (Case::AllCaps, Case::AllCaps) => true,
            (Case::Kebab, Case::Kebab) => true,
            (Case::Train, Case::Train) => true,
            (Case::Spaced, Case::Spaced) => true,
            (Case::Title, Case::Title) => true,
            (Case::HttpHeader, Case::HttpHeader) => true,
            (Case::None, Case::None) => true,
            _ => false,
        }
    }
}

fn stderr(msg: &str) {
    eprintln!("{}", msg);
    process::exit(1);
}

fn stdout(msg: &str) {
    println!("{}", msg);
    process::exit(0);
}

fn main() {
    let args = Args::parse();

    let case = get_case(&args.word);

    if case == Case::None {
        stderr("Invalid input");
    }

    let result = match args {
        Args { flat: true, .. } => flat_case(&args.word),
        Args { upper: true, .. } => upper_case(&args.word),
        Args { camel: true, .. } => camel_case(&args.word, case),
        Args { pascal: true, .. } => pascal_case(&args.word, case),
        Args { snake: true, .. } => snake_case(&args.word, case),
        Args { all_caps: true, .. } => all_caps_case(&args.word, case),
        Args { kebab: true, .. } => kebab_case(&args.word, case),
        Args { train: true, .. } => train_case(&args.word, case),
        Args { spaced: true, .. } => spaced_case(&args.word, case),
        Args { title: true, .. } => title_case(&args.word, case),
        Args {
            http_header: true, ..
        } => http_header_case(&args.word, case),
        _ => case.to_string(),
    };

    stdout(&result);
}

fn is_first_upper(word: &str) -> bool {
    let v: Vec<char> = word.chars().collect();
    v[0].is_uppercase()
}

fn get_case(word: &str) -> Case {
    let contains_dash = word.contains('-');
    let contains_underscore = word.contains('_');
    let contains_space = word.contains(' ');

    if contains_dash && contains_underscore
        || contains_dash && contains_space
        || contains_underscore && contains_space
    {
        return Case::None;
    }

    let is_lowercased = word.to_lowercase() == word;

    if !contains_dash && !contains_underscore {
        if contains_space {
            if is_lowercased {
                return Case::Spaced;
            } else if word.to_uppercase() == word {
                return Case::Spaced;
            } else if is_title_case(&word) {
                return Case::Title;
            }
        } else {

        if is_lowercased {
            return Case::Flat;
        } else if word.to_uppercase() == word {
            return Case::Upper;
        }

        if is_first_upper(&word) {
            return Case::Pascal;
        } else {
            return Case::Camel;
        }
        }

    }

    if contains_underscore {
        if is_lowercased {
            return Case::Snake;
        } else if word.to_uppercase() == word {
            return Case::AllCaps;
        }
    }

    if contains_dash {
        if is_lowercased {
            return Case::Kebab;
        } else if word.to_uppercase() == word {
            return Case::Train;
        } else if is_http_header_case(&word) {
            return Case::HttpHeader;
        }
    }

    return Case::None;
}

fn is_token_capitalized_case(word: &str, token: char) -> bool {
    let mut first = true;
    let mut found_token = false;

    for c in word.chars() {
        if first {
            if !c.is_uppercase() || c == '-' || c == '_' || c == ' ' {
                return false;
            }
            first = false;
            continue;
        }

        if found_token {
            if !c.is_uppercase() {
                return false;
            }
            found_token = false;
            continue;
        }

        if c == token {
            found_token = true;
            continue;
        }

        if c.is_uppercase() {
            return false;
        }
    }

    return true;
}

fn is_http_header_case(word: &str) -> bool {
    is_token_capitalized_case(word, '-')
}

fn is_title_case(word: &str) -> bool {
    is_token_capitalized_case(word, ' ')
}

fn capitalize_first_letter(word: &str) -> String {
    let mut v: Vec<char> = word.chars().collect();
    v[0] = v[0].to_uppercase().nth(0).unwrap();
    return v.into_iter().collect();
}

fn lower_first_letter(word: &str) -> String {
    let mut v: Vec<char> = word.chars().collect();
    v[0] = v[0].to_lowercase().nth(0).unwrap();
    return v.into_iter().collect();
}

fn flat_word(word: &str) -> String {
    word.replace("-", "").replace("_", "").replace(" ", "")
}

fn flat_case(word: &str) -> String {
    flat_word(word).to_lowercase()
}

fn upper_case(word: &str) -> String {
    flat_word(word).to_uppercase()
}

fn camel_case(word: &str, case: Case) -> String {
    match case {
        Case::Camel => return word.to_string(),
        Case::Pascal => return lower_first_letter(word),
        Case::Flat => return word.to_string(),
        Case::Upper => return word.to_string().to_lowercase(),
        _ => (),
    }

    let mut result = String::new();
    let mut first = true;

    for part in word.split(|c| c == '-' || c == '_' || c == ' ') {
        if first {
            result.push_str(&part.to_lowercase());
            first = false;
        } else {
            result.push_str(&part[..1].to_uppercase());
            result.push_str(&part[1..].to_lowercase());
        }
    }

    lower_first_letter(&result)
}

fn pascal_case(word: &str, case: Case) -> String {
    match case {
        Case::Camel => return capitalize_first_letter(word),
        Case::Pascal => return word.to_string(),
        Case::Flat => return capitalize_first_letter(word),
        Case::Upper => return capitalize_first_letter(&word.to_lowercase()),
        _ => capitalize_first_letter(&camel_case(word, case)),
    }
}

fn camel_or_pascal_to_token(word: &str, token: char) -> String {
    let mut result = String::new();
    let mut first = true;

    for c in word.chars() {
        if first {
            result.push(c.to_lowercase().nth(0).unwrap());
            first = false;
            continue;
        }

        if c.is_uppercase() {
            result.push(token);
            result.push(c.to_lowercase().nth(0).unwrap());
        } else {
            result.push(c);
        }
    }

    result
}

fn camel_or_pascal_to_upper_token(word: &str, token: char) -> String {
    let mut result = String::new();
    let mut first = true;

    for c in word.chars() {
        if first {
            result.push(c.to_uppercase().nth(0).unwrap());
            first = false;
            continue;
        }

        if c.is_uppercase() {
            result.push(token);
        }

        result.push(c);
    }

    result
}

fn snake_case(word: &str, case: Case) -> String {
    match case {
        Case::Snake => return word.to_string(),
        Case::AllCaps => return word.to_lowercase(),
        Case::Flat => return word.to_string(),
        Case::Upper => return word.to_lowercase(),
        Case::Kebab => return word.replace("-", "_"),
        Case::Train => return word.replace("-", "_").to_lowercase(),
        Case::Spaced => return word.replace(" ", "_").to_lowercase(),
        Case::Title => return word.replace(" ", "_").to_lowercase(),
        Case::HttpHeader => return word.replace("-", "_").to_lowercase(),
        _ => camel_or_pascal_to_token(word, '_'),
    }
}

fn all_caps_case(word: &str, case: Case) -> String {
    match case {
        Case::Snake => return word.to_uppercase(),
        Case::AllCaps => return word.to_string(),
        Case::Flat => return word.to_uppercase(),
        Case::Upper => return word.to_string(),
        Case::Kebab => return word.replace("-", "_").to_uppercase(),
        Case::Train => return word.replace("-", "_"),
        Case::Spaced => return word.replace(" ", "_").to_uppercase(),
        Case::Title => return word.replace(" ", "_").to_uppercase(),
        Case::HttpHeader => return word.replace("-", "_").to_uppercase(),
        _ => camel_or_pascal_to_token(word, '_').to_uppercase(),
    }
}

fn kebab_case(word: &str, case: Case) -> String {
    match case {
        Case::Snake => return word.replace("_", "-"),
        Case::AllCaps => return word.replace("_", "-").to_lowercase(),
        Case::Flat => return word.to_string(),
        Case::Upper => return word.to_lowercase(),
        Case::Kebab => return word.to_string(),
        Case::Train => return word.to_lowercase(),
        Case::Spaced => return word.replace(" ", "-"),
        Case::Title => return word.replace(" ", "-").to_lowercase(),
        Case::HttpHeader => return word.to_lowercase(),
        _ => camel_or_pascal_to_token(word, '-'),
    }
}

fn train_case(word: &str, case: Case) -> String {
    match case {
        Case::Snake => return word.replace("_", "-").to_uppercase(),
        Case::AllCaps => return word.replace("_", "-"),
        Case::Flat => return word.to_uppercase(),
        Case::Upper => return word.to_string(),
        Case::Kebab => return word.to_uppercase(),
        Case::Train => return word.to_string(),
        Case::Spaced => return word.replace(" ", "-").to_uppercase(),
        Case::Title => return word.replace(" ", "-").to_uppercase(),
        Case::HttpHeader => return word.replace("_", "-").to_uppercase(),
        _ => camel_or_pascal_to_token(word, '-').to_uppercase(),
    }
}

fn spaced_case(word: &str, case: Case) -> String {
    match case {
        Case::Snake => return word.replace("_", " "),
        Case::AllCaps => return word.replace("_", " ").to_lowercase(),
        Case::Flat => return word.to_string(),
        Case::Upper => return word.to_lowercase(),
        Case::Kebab => return word.replace("-", " "),
        Case::Train => return word.replace("-", " ").to_lowercase(),
        Case::Spaced => return word.to_string(),
        Case::Title => return word.to_lowercase(),
        Case::HttpHeader => return word.replace("-", " ").to_lowercase(),
        _ => camel_or_pascal_to_token(word, ' '),
    }
}

fn token_to_token_capitalized_case(word: &str, token: char) -> String {
    let mut result = String::new();
    let mut first = true;

    for part in word
        .to_lowercase()
        .split(|c| c == '-' || c == '_' || c == ' ')
    {
        if first {
            result.push_str(&part[..1].to_uppercase());
            result.push_str(&part[1..].to_lowercase());
            first = false;
        } else {
            result.push(token);
            result.push_str(&part[..1].to_uppercase());
            result.push_str(&part[1..].to_lowercase());
        }
    }

    result
}

fn title_case(word: &str, case: Case) -> String {
    match case {
        Case::Flat => return capitalize_first_letter(&word),
        Case::Upper => return capitalize_first_letter(&word.to_lowercase()),
        Case::Camel => return camel_or_pascal_to_upper_token(word, ' '),
        Case::Pascal => return camel_or_pascal_to_upper_token(word, ' '),
        _ => token_to_token_capitalized_case(word, ' '),
    }
}

fn http_header_case(word: &str, case: Case) -> String {
    match case {
        Case::Flat => return capitalize_first_letter(&word),
        Case::Upper => return capitalize_first_letter(&word.to_lowercase()),
        Case::Camel => return camel_or_pascal_to_upper_token(word, '-'),
        Case::Pascal => return camel_or_pascal_to_upper_token(word, '-'),
        _ => token_to_token_capitalized_case(word, '-'),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    static FLAT: &str = "helloworld";
    static UPPER: &str = "HELLOWORLD";
    static CAMEL: &str = "helloWorld";
    static BROKEN_PASCAL: &str = "Helloworld";
    static PASCAL: &str = "HelloWorld";
    static SNAKE: &str = "hello_world";
    static ALL_CAPS: &str = "HELLO_WORLD";
    static KEBAB: &str = "hello-world";
    static TRAIN: &str = "HELLO-WORLD";
    static SPACED: &str = "hello world";
    static TITLE: &str = "Hello World";
    static HTTP_HEADER: &str = "Hello-World";

    #[test]
    fn test_flat_case() {
        assert_eq!(flat_case(FLAT), FLAT);
        assert_eq!(flat_case(UPPER), FLAT);
        assert_eq!(flat_case(CAMEL), FLAT);
        assert_eq!(flat_case(PASCAL), FLAT);
        assert_eq!(flat_case(SNAKE), FLAT);
        assert_eq!(flat_case(ALL_CAPS), FLAT);
        assert_eq!(flat_case(KEBAB), FLAT);
        assert_eq!(flat_case(TRAIN), FLAT);
        assert_eq!(flat_case(SPACED), FLAT);
        assert_eq!(flat_case(TITLE), FLAT);
        assert_eq!(flat_case(HTTP_HEADER), FLAT);
    }

    #[test]
    fn test_upper_case() {
        assert_eq!(upper_case(FLAT), UPPER);
        assert_eq!(upper_case(UPPER), UPPER);
        assert_eq!(upper_case(CAMEL), UPPER);
        assert_eq!(upper_case(PASCAL), UPPER);
        assert_eq!(upper_case(SNAKE), UPPER);
        assert_eq!(upper_case(ALL_CAPS), UPPER);
        assert_eq!(upper_case(KEBAB), UPPER);
        assert_eq!(upper_case(TRAIN), UPPER);
        assert_eq!(upper_case(SPACED), UPPER);
        assert_eq!(upper_case(TITLE), UPPER);
        assert_eq!(upper_case(HTTP_HEADER), UPPER);
    }

    #[test]
    fn test_camel_case() {
        fn camel_case_helper(word: &str) -> String {
            camel_case(word, get_case(word))
        }

        assert_eq!(camel_case_helper(FLAT), FLAT);
        assert_eq!(camel_case_helper(UPPER), FLAT);
        assert_eq!(camel_case_helper(CAMEL), CAMEL);
        assert_eq!(camel_case_helper(PASCAL), CAMEL);
        assert_eq!(camel_case_helper(SNAKE), CAMEL);
        assert_eq!(camel_case_helper(ALL_CAPS), CAMEL);
        assert_eq!(camel_case_helper(KEBAB), CAMEL);
        assert_eq!(camel_case_helper(TRAIN), CAMEL);
        assert_eq!(camel_case_helper(SPACED), CAMEL);
        assert_eq!(camel_case_helper(TITLE), CAMEL);
        assert_eq!(camel_case_helper(HTTP_HEADER), CAMEL);
    }

    #[test]
    fn test_pascal_case() {
        fn pascal_case_helper(word: &str) -> String {
            pascal_case(word, get_case(word))
        }

        assert_eq!(pascal_case_helper(FLAT), BROKEN_PASCAL);
        assert_eq!(pascal_case_helper(UPPER), BROKEN_PASCAL);
        assert_eq!(pascal_case_helper(CAMEL), PASCAL);
        assert_eq!(pascal_case_helper(PASCAL), PASCAL);
        assert_eq!(pascal_case_helper(SNAKE), PASCAL);
        assert_eq!(pascal_case_helper(ALL_CAPS), PASCAL);
        assert_eq!(pascal_case_helper(KEBAB), PASCAL);
        assert_eq!(pascal_case_helper(TRAIN), PASCAL);
        assert_eq!(pascal_case_helper(SPACED), PASCAL);
        assert_eq!(pascal_case_helper(TITLE), PASCAL);
        assert_eq!(pascal_case_helper(HTTP_HEADER), PASCAL);
    }

    #[test]
    fn test_snake_case() {
        fn snake_case_helper(word: &str) -> String {
            snake_case(word, get_case(word))
        }

        assert_eq!(snake_case_helper(FLAT), FLAT);
        assert_eq!(snake_case_helper(UPPER), FLAT);
        assert_eq!(snake_case_helper(CAMEL), SNAKE);
        assert_eq!(snake_case_helper(PASCAL), SNAKE);
        assert_eq!(snake_case_helper(SNAKE), SNAKE);
        assert_eq!(snake_case_helper(ALL_CAPS), SNAKE);
        assert_eq!(snake_case_helper(KEBAB), SNAKE);
        assert_eq!(snake_case_helper(TRAIN), SNAKE);
        assert_eq!(snake_case_helper(SPACED), SNAKE);
        assert_eq!(snake_case_helper(TITLE), SNAKE);
        assert_eq!(snake_case_helper(HTTP_HEADER), SNAKE);
    }

    #[test]
    fn test_all_caps_case() {
        fn all_caps_case_helper(word: &str) -> String {
            all_caps_case(word, get_case(word))
        }

        assert_eq!(all_caps_case_helper(FLAT), UPPER);
        assert_eq!(all_caps_case_helper(UPPER), UPPER);
        assert_eq!(all_caps_case_helper(CAMEL), ALL_CAPS);
        assert_eq!(all_caps_case_helper(PASCAL), ALL_CAPS);
        assert_eq!(all_caps_case_helper(SNAKE), ALL_CAPS);
        assert_eq!(all_caps_case_helper(ALL_CAPS), ALL_CAPS);
        assert_eq!(all_caps_case_helper(KEBAB), ALL_CAPS);
        assert_eq!(all_caps_case_helper(TRAIN), ALL_CAPS);
        assert_eq!(all_caps_case_helper(SPACED), ALL_CAPS);
        assert_eq!(all_caps_case_helper(TITLE), ALL_CAPS);
        assert_eq!(all_caps_case_helper(HTTP_HEADER), ALL_CAPS);
    }

    #[test]
    fn test_kebab_case() {
        fn kebab_case_helper(word: &str) -> String {
            kebab_case(word, get_case(word))
        }

        assert_eq!(kebab_case_helper(FLAT), FLAT);
        assert_eq!(kebab_case_helper(UPPER), FLAT);
        assert_eq!(kebab_case_helper(CAMEL), KEBAB);
        assert_eq!(kebab_case_helper(PASCAL), KEBAB);
        assert_eq!(kebab_case_helper(SNAKE), KEBAB);
        assert_eq!(kebab_case_helper(ALL_CAPS), KEBAB);
        assert_eq!(kebab_case_helper(KEBAB), KEBAB);
        assert_eq!(kebab_case_helper(TRAIN), KEBAB);
        assert_eq!(kebab_case_helper(SPACED), KEBAB);
        assert_eq!(kebab_case_helper(TITLE), KEBAB);
        assert_eq!(kebab_case_helper(HTTP_HEADER), KEBAB);
    }

    #[test]
    fn test_train_case() {
        fn train_case_helper(word: &str) -> String {
            train_case(word, get_case(word))
        }

        assert_eq!(train_case_helper(FLAT), UPPER);
        assert_eq!(train_case_helper(UPPER), UPPER);
        assert_eq!(train_case_helper(CAMEL), TRAIN);
        assert_eq!(train_case_helper(PASCAL), TRAIN);
        assert_eq!(train_case_helper(SNAKE), TRAIN);
        assert_eq!(train_case_helper(ALL_CAPS), TRAIN);
        assert_eq!(train_case_helper(KEBAB), TRAIN);
        assert_eq!(train_case_helper(TRAIN), TRAIN);
        assert_eq!(train_case_helper(SPACED), TRAIN);
        assert_eq!(train_case_helper(TITLE), TRAIN);
        assert_eq!(train_case_helper(HTTP_HEADER), TRAIN);
    }

    #[test]
    fn test_spaced_case() {
        fn spaced_case_helper(word: &str) -> String {
            spaced_case(word, get_case(word))
        }

        assert_eq!(spaced_case_helper(FLAT), FLAT);
        assert_eq!(spaced_case_helper(UPPER), FLAT);
        assert_eq!(spaced_case_helper(CAMEL), SPACED);
        assert_eq!(spaced_case_helper(PASCAL), SPACED);
        assert_eq!(spaced_case_helper(SNAKE), SPACED);
        assert_eq!(spaced_case_helper(ALL_CAPS), SPACED);
        assert_eq!(spaced_case_helper(KEBAB), SPACED);
        assert_eq!(spaced_case_helper(TRAIN), SPACED);
        assert_eq!(spaced_case_helper(SPACED), SPACED);
        assert_eq!(spaced_case_helper(TITLE), SPACED);
        assert_eq!(spaced_case_helper(HTTP_HEADER), SPACED);
    }

    #[test]
    fn test_title_case() {
        fn title_case_helper(word: &str) -> String {
            title_case(word, get_case(word))
        }

        assert_eq!(title_case_helper(FLAT), BROKEN_PASCAL);
        assert_eq!(title_case_helper(UPPER), BROKEN_PASCAL);
        assert_eq!(title_case_helper(CAMEL), TITLE);
        assert_eq!(title_case_helper(PASCAL), TITLE);
        assert_eq!(title_case_helper(SNAKE), TITLE);
        assert_eq!(title_case_helper(ALL_CAPS), TITLE);
        assert_eq!(title_case_helper(KEBAB), TITLE);
        assert_eq!(title_case_helper(TRAIN), TITLE);
        assert_eq!(title_case_helper(SPACED), TITLE);
        assert_eq!(title_case_helper(TITLE), TITLE);
        assert_eq!(title_case_helper(HTTP_HEADER), TITLE);
    }

    #[test]
    fn test_http_header_case() {
        fn http_header_case_helper(word: &str) -> String {
            http_header_case(word, get_case(word))
        }

        assert_eq!(http_header_case_helper(FLAT), BROKEN_PASCAL);
        assert_eq!(http_header_case_helper(UPPER), BROKEN_PASCAL);
        assert_eq!(http_header_case_helper(CAMEL), HTTP_HEADER);
        assert_eq!(http_header_case_helper(PASCAL), HTTP_HEADER);
        assert_eq!(http_header_case_helper(SNAKE), HTTP_HEADER);
        assert_eq!(http_header_case_helper(ALL_CAPS), HTTP_HEADER);
        assert_eq!(http_header_case_helper(KEBAB), HTTP_HEADER);
        assert_eq!(http_header_case_helper(TRAIN), HTTP_HEADER);
        assert_eq!(http_header_case_helper(SPACED), HTTP_HEADER);
        assert_eq!(http_header_case_helper(TITLE), HTTP_HEADER);
        assert_eq!(http_header_case_helper(HTTP_HEADER), HTTP_HEADER);
    }

    #[test]
    fn test_get_case() {
        assert_eq!(get_case(FLAT), Case::Flat);
        assert_eq!(get_case(UPPER), Case::Upper);
        assert_eq!(get_case(CAMEL), Case::Camel);
        assert_eq!(get_case(PASCAL), Case::Pascal);
        assert_eq!(get_case(SNAKE), Case::Snake);
        assert_eq!(get_case(ALL_CAPS), Case::AllCaps);
        assert_eq!(get_case(KEBAB), Case::Kebab);
        assert_eq!(get_case(TRAIN), Case::Train);
        assert_eq!(get_case(SPACED), Case::Spaced);
        assert_eq!(get_case(TITLE), Case::Title);
        assert_eq!(get_case(HTTP_HEADER), Case::HttpHeader);

        assert_eq!(get_case("hello-new_world"), Case::None);
        assert_eq!(get_case("hello-World"), Case::None);
        assert_eq!(get_case("hello new-world"), Case::None);
        assert_eq!(get_case("hello_new-world of programming"), Case::None);
        assert_eq!(get_case("Broken-HttP-Header"), Case::None);
        assert_eq!(get_case("Broken space case"), Case::None);
    }
}
