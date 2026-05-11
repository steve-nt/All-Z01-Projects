use chrono::{NaiveDate, Weekday, Duration, Datelike};

pub fn middle_day(year: u32) -> Option<Weekday> {
    let year_i32 = year as i32;

    // Leap year? If yes, no single middle day
    let is_leap = NaiveDate::from_ymd_opt(year_i32, 2, 29).is_some();
    if is_leap {
        return None;
    }

    // 365 days → middle is 183rd day (0-based index 182)
    let jan_1 = NaiveDate::from_ymd_opt(year_i32, 1, 1)?;
    let middle = jan_1 + Duration::days(182);

    Some(middle.weekday())
}