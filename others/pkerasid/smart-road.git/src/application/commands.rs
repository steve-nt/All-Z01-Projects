use crate::application::spawn::SpawnRequest;

/// Commands accepted by the application layer from infrastructure adapters.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Command {
    SpawnVehicle(SpawnRequest),
}
