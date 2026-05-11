use std::path::PathBuf;
use std::process::ExitCode;

use localhost::interface::config_parser;
use localhost::{VERSION, log_error, log_info, logger};

const DEFAULT_CONFIG: &str = "config/default.toml";

fn main() -> ExitCode {
    logger::set_max_level(logger::Level::Info);

    let config_path = std::env::args()
        .nth(1)
        .map_or_else(|| PathBuf::from(DEFAULT_CONFIG), PathBuf::from);

    log_info!(
        "localhost {VERSION} starting; config={}",
        config_path.display()
    );

    let cfg = match config_parser::load_file(&config_path) {
        Ok(c) => c,
        Err(e) => {
            log_error!("config load failed: {e}");
            return ExitCode::from(2);
        }
    };

    log_info!(
        "loaded {} host(s) on {} listener(s)",
        cfg.hosts().len(),
        cfg.distinct_listeners().len()
    );
    for addr in cfg.distinct_listeners() {
        log_info!("listener: {addr}");
    }

    match start_event_loop(&cfg) {
        Ok(()) => ExitCode::SUCCESS,
        Err(e) => {
            log_error!("event loop exited: {e}");
            ExitCode::from(1)
        }
    }
}

#[cfg(target_os = "linux")]
fn start_event_loop(cfg: &localhost::domain::config::server::ServerConfig) -> std::io::Result<()> {
    use std::rc::Rc;

    use std::cell::RefCell;

    use localhost::application::event_loop::EventLoop;
    use localhost::application::request_pipeline::PipelineContext;
    use localhost::infrastructure::cgi::OsProcessRunner;
    use localhost::infrastructure::clock::SystemClock;
    use localhost::infrastructure::fs::OsFileSystem;
    use localhost::infrastructure::reactor::epoll::EpollReactor;
    use localhost::infrastructure::session_store::memory::MemorySessionStore;

    let reactor = EpollReactor::new()?;
    let clock = SystemClock::new();
    let pipeline = Rc::new(PipelineContext {
        config: Rc::new(cfg.clone()),
        fs: Rc::new(OsFileSystem::new()),
        process_runner: Rc::new(OsProcessRunner::new()),
        session_store: Some(Rc::new(RefCell::new(MemorySessionStore::new()))),
    });
    let mut el = EventLoop::new(reactor, clock).with_pipeline(pipeline);
    for addr in cfg.distinct_listeners() {
        el.bind(addr)?;
        log_info!("bound {addr}");
    }
    log_info!("event loop running");
    el.run()
}

#[cfg(not(target_os = "linux"))]
fn start_event_loop(_cfg: &localhost::domain::config::server::ServerConfig) -> std::io::Result<()> {
    Err(std::io::Error::other(
        "event loop requires Linux (epoll); the binary only runs on Linux",
    ))
}
