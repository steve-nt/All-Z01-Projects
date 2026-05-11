use std::fs::File;
use std::io::Write;
use std::path::Path;

pub fn open_or_create<P: AsRef<Path>>(path: &P, content: &str) {
    let mut file = File::options()
        .append(true)
        .create(true)
        .open(path)
        .unwrap();

    file.write_all(content.as_bytes()).unwrap();
}