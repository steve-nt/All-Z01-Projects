use rand::Rng;
use rand::rngs::ThreadRng;

use crate::application::spawn::SpawnRequest;
use crate::domain::lane::{Direction, Route};
use crate::domain::vehicle::SpeedTier;

#[derive(Debug)]
pub struct VehicleGenerator<R = ThreadRng> {
    rng: R,
}

impl VehicleGenerator<ThreadRng> {
    #[must_use]
    pub fn new() -> Self {
        Self {
            rng: rand::thread_rng(),
        }
    }
}

impl Default for VehicleGenerator<ThreadRng> {
    fn default() -> Self {
        Self::new()
    }
}

impl<R: Rng> VehicleGenerator<R> {
    #[must_use]
    pub const fn with_rng(rng: R) -> Self {
        Self { rng }
    }

    #[must_use]
    pub fn spawn_request_for_origin(&mut self, origin: Direction) -> SpawnRequest {
        SpawnRequest {
            origin,
            route: choose_random(&mut self.rng, &Route::ALL),
            speed_tier: SpeedTier::High,
        }
    }

    #[must_use]
    pub fn random_spawn_request(&mut self) -> SpawnRequest {
        let origin = choose_random(&mut self.rng, &Direction::ALL);
        self.spawn_request_for_origin(origin)
    }
}

fn choose_random<T: Copy, R: Rng, const N: usize>(rng: &mut R, items: &[T; N]) -> T {
    let idx = rng.gen_range(0..N);
    items[idx]
}

#[cfg(test)]
mod tests {
    use rand::SeedableRng;
    use rand::rngs::StdRng;

    use super::*;

    #[test]
    fn generator_produces_requests_from_known_enums() {
        let mut generator = VehicleGenerator::with_rng(StdRng::seed_from_u64(7));
        let request = generator.random_spawn_request();

        assert!(Direction::ALL.contains(&request.origin));
        assert!(Route::ALL.contains(&request.route));
        assert!(SpeedTier::ALL.contains(&request.speed_tier));
    }
}
