use std::collections::HashMap;

fn leap(y: i32) -> bool {(y % 4 == 0 && y % 100 != 0) || y % 400 == 0}
fn wday(y: i32, m: i32, d: i32) -> i32 {
	let (mut mm, mut yy) = (m, y);
	if mm < 3 { mm += 12; yy -= 1; }
	let k = yy % 100;
	let j = yy / 100;
	let h = (d + (13 * (mm + 1)) / 5 + k + k / 4 + j / 4 + 5 * j) % 7; // 0=Sat
	match h {0=>6,1=>7,2=>1,3=>2,4=>3,5=>4,6=>5,_=>1}
}
fn weeks_in_year(y: i32) -> i32 {
	let jan1 = wday(y, 1, 1);
	if jan1 == 4 || (jan1 == 3 && leap(y)) { 53 } else { 52 }
}
fn iso_week(y: i32, m: i32, d: i32) -> (i32, i32) {
	let cum = [0,31,59,90,120,151,181,212,243,273,304,334];
	let doy = cum[(m - 1) as usize] + d + if m > 2 && leap(y) { 1 } else { 0 };
	let dow = wday(y, m, d); // Mon=1..7
	let mut week = (doy - dow + 10) / 7;
	let mut iso_y = y;
	let this = weeks_in_year(y);
	if week < 1 { iso_y -= 1; week = weeks_in_year(iso_y); }
	else if week > this { iso_y += 1; week = 1; }
	(iso_y, week)
}

fn parse_ymd(s: &str) -> Option<(i32, i32, i32)> {
	if s.len() < 10 { return None; }
	let y = s[0..4].parse().ok()?;
	let m = s[5..7].parse().ok()?;
	let d = s[8..10].parse().ok()?;
	Some((y, m, d))
}

fn for_each_commit<'a, F: FnMut(&'a json::JsonValue)>(data: &'a json::JsonValue, mut f: F) {
	if data.is_array() { for c in data.members() { f(c); } return; }
	let a = &data["commits"];
	if a.is_array() { for c in a.members() { f(c); } }
}

pub fn commits_per_week(data: &json::JsonValue) -> HashMap<String, u32> {
	let mut map: HashMap<String, u32> = HashMap::new();
	for_each_commit(data, |c| {
		let date = c["commit"]["author"]["date"].as_str().or_else(|| c["commit"]["committer"]["date"].as_str());
		if let Some(date) = date {
			if let Some((y, m, d)) = parse_ymd(date) {
				let (yy, w) = iso_week(y, m, d);
				let k = format!("{}-W{}", yy, w);
				*map.entry(k).or_insert(0) += 1;
			}
		}
	});
	map
}

pub fn commits_per_author(data: &json::JsonValue) -> HashMap<String, u32> {
	let mut map: HashMap<String, u32> = HashMap::new();
	for_each_commit(data, |c| {
		let login = c["author"]["login"].as_str().or_else(|| c["committer"]["login"].as_str());
		if let Some(l) = login { *map.entry(l.to_string()).or_insert(0) += 1; }
	});
	map
}
