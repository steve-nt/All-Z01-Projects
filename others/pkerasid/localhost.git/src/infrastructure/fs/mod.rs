//! `OsFileSystem` — real `std::fs` adapter implementing the `FileSystem` port.

use std::io;
use std::path::Path;
use std::time::UNIX_EPOCH;

use crate::application::ports::filesystem::{DirEntry, FileMetadata, FileSystem};

#[derive(Debug, Default)]
pub struct OsFileSystem;

impl OsFileSystem {
    pub fn new() -> Self {
        Self
    }
}

impl FileSystem for OsFileSystem {
    fn stat(&self, path: &Path) -> io::Result<FileMetadata> {
        let m = std::fs::metadata(path)?;
        let modified_secs = m
            .modified()
            .ok()
            .and_then(|t| t.duration_since(UNIX_EPOCH).ok())
            .map_or(0, |d| d.as_secs());
        Ok(FileMetadata {
            is_dir: m.is_dir(),
            is_file: m.is_file(),
            size: m.len(),
            modified_secs,
        })
    }

    fn read_file(&self, path: &Path) -> io::Result<Vec<u8>> {
        std::fs::read(path)
    }

    fn read_dir(&self, path: &Path) -> io::Result<Vec<DirEntry>> {
        let mut entries = Vec::new();
        for entry in std::fs::read_dir(path)? {
            let entry = entry?;
            let meta = entry.metadata()?;
            let name = entry.file_name().into_string().unwrap_or_default();
            let modified_secs = meta
                .modified()
                .ok()
                .and_then(|t| t.duration_since(UNIX_EPOCH).ok())
                .map_or(0, |d| d.as_secs());
            entries.push(DirEntry {
                name,
                is_dir: meta.is_dir(),
                size: meta.len(),
                modified_secs,
            });
        }
        entries.sort_by(|a, b| a.name.cmp(&b.name));
        Ok(entries)
    }

    fn remove_file(&self, path: &Path) -> io::Result<()> {
        std::fs::remove_file(path)
    }

    fn write_file(&self, path: &Path, data: &[u8]) -> io::Result<()> {
        std::fs::write(path, data)
    }
}
