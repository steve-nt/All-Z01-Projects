use std::time::{SystemTime, UNIX_EPOCH};

fn format_utc_datetime() -> String {
    let secs = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_secs() as i64;

    const SECS_PER_DAY: i64 = 86400;
    let days = secs / SECS_PER_DAY;
    let time_secs = (secs % SECS_PER_DAY + SECS_PER_DAY) % SECS_PER_DAY;

    // Unix epoch to date (simplified UTC)
    let mut d = days + 719468;
    let era = (if d >= 0 { d } else { d - 146096 }) / 146097;
    let day_of_era = d - era * 146097;
    let year = (day_of_era - day_of_era / 1460 + day_of_era / 36524 - day_of_era / 146096) / 365;
    let y = year + era * 400;
    let day_of_year = day_of_era - (365 * year + year / 4 - year / 100);
    let mp = (5 * day_of_year + 2) / 153;
    let day = day_of_year - (153 * mp + 2) / 5 + 1;
    let month = mp + (if mp < 10 { 3 } else { -9 });
    let year = y + (if month <= 2 { 1 } else { 0 });

    let hour = (time_secs / 3600) % 24;
    let min = (time_secs / 60) % 60;
    let sec = time_secs % 60;

    format!(
        "{:04}-{:02}-{:02} {:02}:{:02}:{:02}",
        year, month, day, hour, min, sec
    )
}

#[derive(Debug, Eq, PartialEq)]
pub struct FormError {
    pub form_values: (&'static str, String),
    pub date: String,
    pub err: &'static str,
}

impl FormError {
    pub fn new(field_name: &'static str, field_value: String, err: &'static str) -> Self {
        let date = format_utc_datetime();

        FormError {
            form_values: (field_name, field_value),
            date,
            err,
        }
    }
}

#[derive(Debug, Eq, PartialEq)]
pub struct Form {
    pub name: String,
    pub password: String,
}

impl Form {
    pub fn validate(&self) -> Result<(), FormError> {
        // 1. Name must not be empty
        if self.name.is_empty() {
            return Err(FormError::new("name", self.name.clone(), "Username is empty"));
        }

        // 2. Password must be at least 8 chars
        if self.password.chars().count() < 8 {
            return Err(FormError::new(
                "password",
                self.password.clone(),
                "Password should be at least 8 characters long",
            ));
        }

        // 3. Password must contain letters, digits, and symbols
        let mut has_letter = false;
        let mut has_digit = false;
        let mut has_symbol = false;

        for c in self.password.chars() {
            if c.is_ascii_alphabetic() {
                has_letter = true;
            } else if c.is_ascii_digit() {
                has_digit = true;
            } else if c.is_ascii() && !c.is_ascii_alphanumeric() {
                has_symbol = true;
            }
        }

        if !(has_letter && has_digit && has_symbol) {
            return Err(FormError::new(
                "password",
                self.password.clone(),
                "Password should be a combination of ASCII numbers, letters and symbols",
            ));
        }

        Ok(())
    }
}