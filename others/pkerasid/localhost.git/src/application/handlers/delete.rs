//! DELETE handler: removes a file under the route root.
//!
//! Returns 204 No Content on success, 404 if the file doesn't exist, 403 if
//! the path resolves to a directory, and 500 on any other I/O error.

use std::path::Path;

use crate::application::ports::filesystem::FileSystem;
use crate::domain::http::response::Response;
use crate::domain::http::status::Status;

pub fn handle(file_path: &Path, fs: &dyn FileSystem) -> Response {
    match fs.stat(file_path) {
        Err(e) if e.kind() == std::io::ErrorKind::NotFound => {
            Response::builder(Status::NOT_FOUND).build()
        }
        Err(_) => Response::builder(Status::INTERNAL_SERVER_ERROR).build(),
        Ok(meta) if meta.is_dir => Response::builder(Status::FORBIDDEN).build(),
        Ok(_) => match fs.remove_file(file_path) {
            Ok(()) => Response::builder(Status::NO_CONTENT).build(),
            Err(e) if e.kind() == std::io::ErrorKind::NotFound => {
                Response::builder(Status::NOT_FOUND).build()
            }
            Err(_) => Response::builder(Status::INTERNAL_SERVER_ERROR).build(),
        },
    }
}

#[cfg(test)]
mod tests {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use std::path::PathBuf;

    use super::*;
    use crate::application::ports::filesystem::fake::FakeFileSystem;

    #[test]
    fn deletes_existing_file_returns_204() {
        let fs = FakeFileSystem::new();
        fs.add_file("/www/upload/old.txt", b"data");
        let resp = handle(&PathBuf::from("/www/upload/old.txt"), &fs);
        assert_eq!(resp.status(), Status::NO_CONTENT);
        // File should be gone.
        assert!(
            fs.stat(std::path::Path::new("/www/upload/old.txt"))
                .is_err()
        );
    }

    #[test]
    fn missing_file_returns_404() {
        let fs = FakeFileSystem::new();
        let resp = handle(&PathBuf::from("/www/upload/ghost.txt"), &fs);
        assert_eq!(resp.status(), Status::NOT_FOUND);
    }

    #[test]
    fn directory_returns_403() {
        let fs = FakeFileSystem::new();
        fs.add_dir("/www/uploads");
        let resp = handle(&PathBuf::from("/www/uploads"), &fs);
        assert_eq!(resp.status(), Status::FORBIDDEN);
    }
}
