//! `FileSystem` port — abstracts file I/O so handlers are testable with a fake.

use std::io;
use std::path::Path;

#[derive(Debug, Clone)]
pub struct FileMetadata {
    pub is_dir: bool,
    pub is_file: bool,
    pub size: u64,
    /// Seconds since Unix epoch, used for `Last-Modified`.
    pub modified_secs: u64,
}

#[derive(Debug, Clone)]
pub struct DirEntry {
    pub name: String,
    pub is_dir: bool,
    pub size: u64,
    pub modified_secs: u64,
}

pub trait FileSystem {
    fn stat(&self, path: &Path) -> io::Result<FileMetadata>;
    fn read_file(&self, path: &Path) -> io::Result<Vec<u8>>;
    fn read_dir(&self, path: &Path) -> io::Result<Vec<DirEntry>>;
    fn remove_file(&self, path: &Path) -> io::Result<()>;
    fn write_file(&self, path: &Path, data: &[u8]) -> io::Result<()>;
}

// --- in-process fake for unit tests ---

#[cfg(test)]
pub mod fake {
    #![allow(
        clippy::unwrap_used,
        clippy::expect_used,
        clippy::panic,
        clippy::cast_possible_truncation
    )]

    use super::*;
    use std::cell::RefCell;
    use std::collections::HashMap;
    use std::path::PathBuf;

    enum FakeEntry {
        File(Vec<u8>),
        Dir,
    }

    #[derive(Default)]
    pub struct FakeFileSystem {
        entries: RefCell<HashMap<PathBuf, FakeEntry>>,
    }

    impl std::fmt::Debug for FakeFileSystem {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            f.debug_struct("FakeFileSystem").finish_non_exhaustive()
        }
    }

    impl FakeFileSystem {
        pub fn new() -> Self {
            Self::default()
        }

        pub fn add_file(&self, path: impl Into<PathBuf>, content: impl Into<Vec<u8>>) {
            self.entries
                .borrow_mut()
                .insert(path.into(), FakeEntry::File(content.into()));
        }

        pub fn add_dir(&self, path: impl Into<PathBuf>) {
            self.entries
                .borrow_mut()
                .insert(path.into(), FakeEntry::Dir);
        }
    }

    impl FileSystem for FakeFileSystem {
        fn stat(&self, path: &Path) -> io::Result<FileMetadata> {
            match self.entries.borrow().get(path) {
                Some(FakeEntry::File(data)) => Ok(FileMetadata {
                    is_dir: false,
                    is_file: true,
                    size: data.len() as u64,
                    modified_secs: 0,
                }),
                Some(FakeEntry::Dir) => Ok(FileMetadata {
                    is_dir: true,
                    is_file: false,
                    size: 0,
                    modified_secs: 0,
                }),
                None => Err(io::Error::from(io::ErrorKind::NotFound)),
            }
        }

        fn read_file(&self, path: &Path) -> io::Result<Vec<u8>> {
            match self.entries.borrow().get(path) {
                Some(FakeEntry::File(data)) => Ok(data.clone()),
                Some(FakeEntry::Dir) => Err(io::Error::other("is a directory")),
                None => Err(io::Error::from(io::ErrorKind::NotFound)),
            }
        }

        fn read_dir(&self, path: &Path) -> io::Result<Vec<DirEntry>> {
            let entries = self.entries.borrow();
            match entries.get(path) {
                Some(FakeEntry::File(_)) => {
                    return Err(io::Error::other("not a directory"));
                }
                Some(FakeEntry::Dir) | None => {}
            }
            let mut result = Vec::new();
            for (p, entry) in entries.iter() {
                if p.parent() == Some(path) {
                    let name = p
                        .file_name()
                        .and_then(|n| n.to_str())
                        .unwrap_or("")
                        .to_owned();
                    result.push(DirEntry {
                        name,
                        is_dir: matches!(entry, FakeEntry::Dir),
                        size: if let FakeEntry::File(d) = entry {
                            d.len() as u64
                        } else {
                            0
                        },
                        modified_secs: 0,
                    });
                }
            }
            result.sort_by(|a, b| a.name.cmp(&b.name));
            Ok(result)
        }

        fn remove_file(&self, path: &Path) -> io::Result<()> {
            let removed = self.entries.borrow_mut().remove(path);
            match removed {
                Some(FakeEntry::File(_)) => Ok(()),
                Some(entry) => {
                    self.entries.borrow_mut().insert(path.to_owned(), entry);
                    Err(io::Error::other("is a directory"))
                }
                None => Err(io::Error::from(io::ErrorKind::NotFound)),
            }
        }

        fn write_file(&self, path: &Path, data: &[u8]) -> io::Result<()> {
            self.entries
                .borrow_mut()
                .insert(path.to_owned(), FakeEntry::File(data.to_vec()));
            Ok(())
        }
    }
}
