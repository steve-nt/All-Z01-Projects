use std::ffi::CString;
use std::io;
use std::os::fd::RawFd;
use std::os::unix::ffi::OsStrExt;
use std::time::{Duration, Instant};

use crate::application::ports::process_runner::{CgiRunSpec, ProcessOutput, ProcessRunner};

#[derive(Debug, Default, Clone, Copy)]
pub struct OsProcessRunner;

impl OsProcessRunner {
    pub fn new() -> Self {
        Self
    }
}

impl ProcessRunner for OsProcessRunner {
    #[allow(clippy::too_many_lines)]
    fn run_cgi(&self, spec: &CgiRunSpec, timeout: Duration) -> io::Result<ProcessOutput> {
        let interpreter = cstring_from_bytes(spec.interpreter.as_os_str().as_bytes())?;
        let script = cstring_from_bytes(spec.script_filename.as_os_str().as_bytes())?;
        let cwd = cstring_from_bytes(spec.cwd.as_os_str().as_bytes())?;

        let mut env_store: Vec<CString> = Vec::with_capacity(spec.env.len());
        for (k, v) in &spec.env {
            env_store.push(cstring_from_bytes(format!("{k}={v}").as_bytes())?);
        }
        let mut envp: Vec<*const libc::c_char> = env_store.iter().map(|e| e.as_ptr()).collect();
        envp.push(std::ptr::null());

        let argv_store = [interpreter.clone(), script];
        let mut argv: Vec<*const libc::c_char> = argv_store.iter().map(|a| a.as_ptr()).collect();
        argv.push(std::ptr::null());

        let stdin_pipe = create_pipe()?;
        let stdout_pipe = create_pipe()?;
        let stderr_pipe = create_pipe()?;

        // SAFETY: `fork` is called in a short critical section before any Rust allocations in child.
        let pid = unsafe { libc::fork() };
        if pid < 0 {
            close_fd_pair(stdin_pipe);
            close_fd_pair(stdout_pipe);
            close_fd_pair(stderr_pipe);
            return Err(io::Error::last_os_error());
        }
        if pid == 0 {
            let exec_spec = ExecSpec {
                stdin_pipe,
                stdout_pipe,
                stderr_pipe,
                cwd: &cwd,
                interpreter: &interpreter,
                argv: &argv,
                envp: &envp,
            };
            child_exec(&exec_spec);
        }

        // SAFETY: parent only needs stdin write and stdout/stderr read ends.
        unsafe {
            let _ = libc::close(stdin_pipe.0);
            let _ = libc::close(stdout_pipe.1);
            let _ = libc::close(stderr_pipe.1);
        }
        let mut stdin_write = AutoFd::new(stdin_pipe.1);
        let mut stdout_read = AutoFd::new(stdout_pipe.0);
        let mut stderr_read = AutoFd::new(stderr_pipe.0);

        let epoll = AutoFd::new(create_epoll()?);
        register_fd(epoll.fd(), stdin_write.fd(), libc::EPOLLOUT as u32)?;
        register_fd(
            epoll.fd(),
            stdout_read.fd(),
            (libc::EPOLLIN | libc::EPOLLHUP) as u32,
        )?;
        register_fd(
            epoll.fd(),
            stderr_read.fd(),
            (libc::EPOLLIN | libc::EPOLLHUP) as u32,
        )?;

        let mut out = Vec::new();
        let mut err = Vec::new();
        let mut stdin_offset = 0usize;
        let mut exit_code: Option<i32> = None;
        let start = Instant::now();

        while stdin_write.is_open()
            || stdout_read.is_open()
            || stderr_read.is_open()
            || exit_code.is_none()
        {
            if start.elapsed() >= timeout {
                // SAFETY: valid child pid from successful fork.
                let _ = unsafe { libc::kill(pid, libc::SIGKILL) };
                let _ = reap_child(pid);
                return Err(io::Error::new(
                    io::ErrorKind::TimedOut,
                    "cgi process timed out",
                ));
            }

            if exit_code.is_none() {
                exit_code = reap_child(pid)?;
            }

            let mut events = [libc::epoll_event { events: 0, u64: 0 }; 8];
            let max_events = i32::try_from(events.len())
                .map_err(|_| io::Error::new(io::ErrorKind::InvalidInput, "too many events"))?;
            // SAFETY: valid epoll fd and writable events buffer.
            let n = unsafe { libc::epoll_wait(epoll.fd(), events.as_mut_ptr(), max_events, 50) };
            if n < 0 {
                let err_no = io::Error::last_os_error();
                if err_no.kind() == io::ErrorKind::Interrupted {
                    continue;
                }
                return Err(err_no);
            }

            let n_events = usize::try_from(n)
                .map_err(|_| io::Error::new(io::ErrorKind::InvalidData, "negative epoll count"))?;
            for ev in events.iter().take(n_events) {
                let fd = RawFd::try_from(ev.u64)
                    .map_err(|_| io::Error::new(io::ErrorKind::InvalidData, "invalid epoll fd"))?;
                let flags = ev.events;
                if stdin_write.is_open() && fd == stdin_write.fd() {
                    if stdin_offset < spec.stdin.len() {
                        // SAFETY: buffer pointer valid for remaining bytes.
                        let written = unsafe {
                            libc::write(
                                stdin_write.fd(),
                                spec.stdin[stdin_offset..].as_ptr().cast(),
                                spec.stdin.len() - stdin_offset,
                            )
                        };
                        if written > 0 {
                            stdin_offset += written.cast_unsigned();
                        } else if written < 0 {
                            let e = io::Error::last_os_error();
                            if !matches!(
                                e.kind(),
                                io::ErrorKind::WouldBlock | io::ErrorKind::Interrupted
                            ) {
                                stdin_write.close();
                            }
                        }
                    }
                    if stdin_offset >= spec.stdin.len() || (flags & libc::EPOLLHUP as u32) != 0 {
                        stdin_write.close();
                    }
                } else if stdout_read.is_open() && fd == stdout_read.fd() {
                    if (flags & libc::EPOLLIN as u32) != 0 {
                        read_available(stdout_read.fd(), &mut out)?;
                    }
                    if (flags & libc::EPOLLHUP as u32) != 0 {
                        stdout_read.close();
                    }
                } else if stderr_read.is_open() && fd == stderr_read.fd() {
                    if (flags & libc::EPOLLIN as u32) != 0 {
                        read_available(stderr_read.fd(), &mut err)?;
                    }
                    if (flags & libc::EPOLLHUP as u32) != 0 {
                        stderr_read.close();
                    }
                }
            }

            // Drain reads even if no explicit HUP arrived.
            if stdout_read.is_open()
                && read_available(stdout_read.fd(), &mut out)? == 0
                && exit_code.is_some()
            {
                stdout_read.close();
            }
            if stderr_read.is_open()
                && read_available(stderr_read.fd(), &mut err)? == 0
                && exit_code.is_some()
            {
                stderr_read.close();
            }
        }

        Ok(ProcessOutput {
            stdout: out,
            stderr: err,
            exit_code,
        })
    }
}

fn cstring_from_bytes(bytes: &[u8]) -> io::Result<CString> {
    CString::new(bytes)
        .map_err(|_| io::Error::new(io::ErrorKind::InvalidInput, "embedded NUL byte"))
}

fn create_pipe() -> io::Result<(RawFd, RawFd)> {
    let mut fds = [0; 2];
    // SAFETY: pointer to 2-int array is valid.
    let rc = unsafe { libc::pipe2(fds.as_mut_ptr(), libc::O_CLOEXEC | libc::O_NONBLOCK) };
    if rc == 0 {
        Ok((fds[0], fds[1]))
    } else {
        Err(io::Error::last_os_error())
    }
}

fn close_fd_pair(pair: (RawFd, RawFd)) {
    // SAFETY: closing best-effort during error cleanup.
    unsafe {
        let _ = libc::close(pair.0);
        let _ = libc::close(pair.1);
    }
}

fn create_epoll() -> io::Result<RawFd> {
    // SAFETY: direct syscall with valid flags.
    let fd = unsafe { libc::epoll_create1(libc::EPOLL_CLOEXEC) };
    if fd >= 0 {
        Ok(fd)
    } else {
        Err(io::Error::last_os_error())
    }
}

fn register_fd(epoll_fd: RawFd, fd: RawFd, events: u32) -> io::Result<()> {
    let fd_u64 = u64::try_from(fd)
        .map_err(|_| io::Error::new(io::ErrorKind::InvalidInput, "negative fd"))?;
    let mut ev = libc::epoll_event {
        events,
        u64: fd_u64,
    };
    // SAFETY: pointers and fds are valid.
    let rc = unsafe { libc::epoll_ctl(epoll_fd, libc::EPOLL_CTL_ADD, fd, &raw mut ev) };
    if rc == 0 {
        Ok(())
    } else {
        Err(io::Error::last_os_error())
    }
}

fn read_available(fd: RawFd, out: &mut Vec<u8>) -> io::Result<usize> {
    let mut total = 0usize;
    let mut buf = [0u8; 4096];
    loop {
        // SAFETY: target buffer is valid and writable.
        let n = unsafe { libc::read(fd, buf.as_mut_ptr().cast(), buf.len()) };
        match n.cmp(&0) {
            std::cmp::Ordering::Greater => {
                let read = n.cast_unsigned();
                out.extend_from_slice(&buf[..read]);
                total += read;
            }
            std::cmp::Ordering::Equal => return Ok(total),
            std::cmp::Ordering::Less => {
                let e = io::Error::last_os_error();
                if matches!(
                    e.kind(),
                    io::ErrorKind::WouldBlock | io::ErrorKind::Interrupted
                ) {
                    return Ok(total);
                }
                return Err(e);
            }
        }
    }
}

fn reap_child(pid: libc::pid_t) -> io::Result<Option<i32>> {
    let mut status = 0;
    // SAFETY: `status` pointer valid; pid from `fork`.
    let rc = unsafe { libc::waitpid(pid, &raw mut status, libc::WNOHANG) };
    if rc == 0 {
        Ok(None)
    } else if rc < 0 {
        let e = io::Error::last_os_error();
        if e.kind() == io::ErrorKind::Interrupted {
            Ok(None)
        } else {
            Err(e)
        }
    } else if libc::WIFEXITED(status) {
        Ok(Some(libc::WEXITSTATUS(status)))
    } else {
        Ok(None)
    }
}

struct ExecSpec<'a> {
    stdin_pipe: (RawFd, RawFd),
    stdout_pipe: (RawFd, RawFd),
    stderr_pipe: (RawFd, RawFd),
    cwd: &'a CString,
    interpreter: &'a CString,
    argv: &'a [*const libc::c_char],
    envp: &'a [*const libc::c_char],
}

fn child_exec(spec: &ExecSpec<'_>) -> ! {
    // SAFETY: child process only; dup2/chdir/execve with validated pointers.
    unsafe {
        let _ = libc::dup2(spec.stdin_pipe.0, libc::STDIN_FILENO);
        let _ = libc::dup2(spec.stdout_pipe.1, libc::STDOUT_FILENO);
        let _ = libc::dup2(spec.stderr_pipe.1, libc::STDERR_FILENO);

        let _ = libc::close(spec.stdin_pipe.0);
        let _ = libc::close(spec.stdin_pipe.1);
        let _ = libc::close(spec.stdout_pipe.0);
        let _ = libc::close(spec.stdout_pipe.1);
        let _ = libc::close(spec.stderr_pipe.0);
        let _ = libc::close(spec.stderr_pipe.1);

        let _ = libc::chdir(spec.cwd.as_ptr());
        libc::execve(
            spec.interpreter.as_ptr(),
            spec.argv.as_ptr(),
            spec.envp.as_ptr(),
        );
        libc::_exit(127);
    }
}

#[derive(Debug)]
struct AutoFd(Option<RawFd>);

impl AutoFd {
    fn new(fd: RawFd) -> Self {
        Self(Some(fd))
    }

    fn fd(&self) -> RawFd {
        self.0.unwrap_or(-1)
    }

    fn is_open(&self) -> bool {
        self.0.is_some()
    }

    fn close(&mut self) {
        if let Some(fd) = self.0.take() {
            // SAFETY: close best-effort on owned fd.
            unsafe {
                let _ = libc::close(fd);
            }
        }
    }
}

impl Drop for AutoFd {
    fn drop(&mut self) {
        self.close();
    }
}
