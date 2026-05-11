use std::{collections::HashMap, num::ParseFloatError};

pub struct Flag<'a> {
    pub short_hand: String,
    pub long_hand: String,
    pub desc: &'a str,
}

impl<'a> Flag<'a> {
    pub fn opt_flag(name: &'a str, d: &'a str) -> Self {
        let mut chars = name.chars();
        let first = chars.next().unwrap_or('-');
        let short_hand = format!("-{}", first);
        let long_hand = format!("--{}", name);

        Flag {
            short_hand,
            long_hand,
            desc: d,
        }
    }
}

pub type Callback = fn(&str, &str) -> Result<String, ParseFloatError>;

pub struct FlagsHandler {
    pub flags: HashMap<String, Callback>,
}

impl FlagsHandler {
    pub fn add_flag(&mut self, flag: Flag, func: Callback) {
        // Register both short and long forms with the same callback
        self.flags.insert(flag.short_hand, func);
        self.flags.insert(flag.long_hand, func);
    }

    pub fn exec_func(&self, input: &str, argv: &[&str]) -> Result<String, String> {
        let func = self
            .flags
            .get(input)
            .ok_or_else(|| "flag not found".to_string())?;

        if argv.len() < 2 {
            return Err("not enough arguments".to_string());
        }

        match func(argv[0], argv[1]) {
            Ok(result) => Ok(result),
            Err(e) => Err(e.to_string()),
        }
    }
}

pub fn div(a: &str, b: &str) -> Result<String, ParseFloatError> {
    let a_val: f64 = a.parse()?;
    let b_val: f64 = b.parse()?;
    let result = a_val / b_val;
    Ok(result.to_string())
}

pub fn rem(a: &str, b: &str) -> Result<String, ParseFloatError> {
    let a_val: f64 = a.parse()?;
    let b_val: f64 = b.parse()?;
    let result = a_val % b_val;
    Ok(result.to_string())
}