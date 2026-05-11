//! `ProcessRunner` port — abstracts CGI process execution.

use std::io;
use std::path::PathBuf;
use std::time::Duration;

#[derive(Debug, Clone)]
pub struct CgiRunSpec {
    pub interpreter: PathBuf,
    pub script_filename: PathBuf,
    pub cwd: PathBuf,
    pub env: Vec<(String, String)>,
    pub stdin: Vec<u8>,
}

#[derive(Debug, Clone)]
pub struct ProcessOutput {
    pub stdout: Vec<u8>,
    pub stderr: Vec<u8>,
    pub exit_code: Option<i32>,
}

pub trait ProcessRunner {
    fn run_cgi(&self, spec: &CgiRunSpec, timeout: Duration) -> io::Result<ProcessOutput>;
}

#[cfg(test)]
pub mod fake {
    #![allow(clippy::unwrap_used, clippy::expect_used, clippy::panic)]

    use std::cell::RefCell;
    use std::collections::VecDeque;
    use std::io;
    use std::time::Duration;

    use super::{CgiRunSpec, ProcessOutput, ProcessRunner};

    #[derive(Debug, Default)]
    pub struct FakeProcessRunner {
        outputs: RefCell<VecDeque<io::Result<ProcessOutput>>>,
    }

    impl FakeProcessRunner {
        pub fn new() -> Self {
            Self::default()
        }

        pub fn push_output(&self, output: io::Result<ProcessOutput>) {
            self.outputs.borrow_mut().push_back(output);
        }
    }

    impl ProcessRunner for FakeProcessRunner {
        fn run_cgi(&self, _spec: &CgiRunSpec, _timeout: Duration) -> io::Result<ProcessOutput> {
            self.outputs.borrow_mut().pop_front().unwrap_or_else(|| {
                Ok(ProcessOutput {
                    stdout: Vec::new(),
                    stderr: Vec::new(),
                    exit_code: Some(0),
                })
            })
        }
    }
}
