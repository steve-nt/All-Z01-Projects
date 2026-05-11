use std::collections::HashMap;

use crate::domain::intersection::{ConflictEntry, RouteId, build_conflict_map};

/// Manages intersection entry reservations for the smart intersection algorithm.
///
/// The manager holds a static conflict map (built once at startup) and a dynamic
/// table of active route holders.  Two routes may proceed simultaneously only if
/// they share no conflict tiles.  Right-turn routes are excluded from the conflict
/// map entirely — they use the outer corner tiles and are never registered here.
///
/// ### Fairness policy (Phase 6 + Phase 8c)
///
/// This type is order-agnostic: [`ReservationManager::request`] grants whenever
/// no conflicting route is occupied.  The **application** layer (`World`) calls
/// `request` in priority order (earliest detection time, plus a small aging bonus
/// for long waits) so FIFO-by-arrival is preserved even when the vehicle `Vec`
/// order differs from detection order.
///
/// Multiple vehicles on the *same* route may hold concurrent reservations (lane
/// safety handles their spacing).
pub struct ReservationManager {
    conflict_map: HashMap<RouteId, Vec<ConflictEntry>>,
    /// route → list of vehicle IDs currently holding a reservation on it.
    active: HashMap<RouteId, Vec<u64>>,
}

impl ReservationManager {
    #[must_use]
    pub fn new() -> Self {
        Self {
            conflict_map: build_conflict_map(),
            active: HashMap::new(),
        }
    }

    /// Attempt to grant an intersection reservation for `vehicle_id` on `route`.
    ///
    /// Returns `true` when:
    /// - the vehicle already holds this reservation, or
    /// - no conflicting route currently has active holders.
    ///
    /// Returns `false` when at least one conflicting route is occupied.
    pub fn request(&mut self, vehicle_id: u64, route: RouteId) -> bool {
        // Already granted — idempotent.
        if self
            .active
            .get(&route)
            .is_some_and(|v| v.contains(&vehicle_id))
        {
            return true;
        }

        // Check whether any conflicting route is currently occupied.
        let blocked = self
            .conflict_map
            .get(&route)
            .is_some_and(|entries| entries.iter().any(|e| self.route_occupied(e.other_route)));

        if blocked {
            return false;
        }

        self.active.entry(route).or_default().push(vehicle_id);
        true
    }

    /// Release the reservation held by `vehicle_id` on `route`.
    ///
    /// No-op if the vehicle did not hold the reservation.
    pub fn release(&mut self, vehicle_id: u64, route: RouteId) {
        if let Some(holders) = self.active.get_mut(&route) {
            holders.retain(|&id| id != vehicle_id);
            if holders.is_empty() {
                self.active.remove(&route);
            }
        }
    }

    /// Returns the conflict entries for `route` — the list of routes that share
    /// tiles with it and the specific shared tiles for each.  Empty slice when
    /// the route has no conflicts (e.g. right-turn routes are absent from the map).
    #[must_use]
    pub fn conflict_entries(
        &self,
        route: RouteId,
    ) -> &[crate::domain::intersection::ConflictEntry] {
        self.conflict_map.get(&route).map_or(&[], Vec::as_slice)
    }

    /// Returns true if no reservations are currently active.
    #[must_use]
    pub fn is_empty(&self) -> bool {
        self.active.is_empty()
    }

    fn route_occupied(&self, route: RouteId) -> bool {
        self.active.get(&route).is_some_and(|v| !v.is_empty())
    }
}

impl Default for ReservationManager {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::domain::intersection::{EntryDir, TurnDir};

    fn ns() -> RouteId {
        (EntryDir::North, TurnDir::Straight)
    }
    fn sn() -> RouteId {
        (EntryDir::South, TurnDir::Straight)
    }
    fn ew() -> RouteId {
        (EntryDir::East, TurnDir::Straight)
    }
    fn ne() -> RouteId {
        (EntryDir::North, TurnDir::Left)
    }

    #[test]
    fn first_request_on_free_route_is_granted() {
        let mut m = ReservationManager::new();
        assert!(m.request(1, ns()));
    }

    #[test]
    fn same_vehicle_requesting_same_route_is_idempotent() {
        let mut m = ReservationManager::new();
        assert!(m.request(1, ns()));
        assert!(m.request(1, ns())); // second call still returns true
    }

    #[test]
    fn conflicting_route_is_blocked_while_first_is_active() {
        let mut m = ReservationManager::new();
        assert!(m.request(1, ns()));
        assert!(!m.request(2, ew())); // NS and EW conflict
    }

    #[test]
    fn non_conflicting_routes_are_granted_simultaneously() {
        let mut m = ReservationManager::new();
        assert!(m.request(1, ns()));
        assert!(m.request(2, sn())); // NS and SN do not conflict
    }

    #[test]
    fn release_unblocks_conflicting_route() {
        let mut m = ReservationManager::new();
        m.request(1, ns());
        assert!(!m.request(2, ew()));
        m.release(1, ns());
        assert!(m.request(2, ew()));
    }

    #[test]
    fn multiple_vehicles_on_same_route_allowed() {
        let mut m = ReservationManager::new();
        assert!(m.request(1, ns()));
        assert!(m.request(2, ns())); // lane safety handles spacing
    }

    #[test]
    fn release_one_of_two_holders_does_not_unblock_conflict() {
        let mut m = ReservationManager::new();
        m.request(1, ns());
        m.request(2, ns());
        m.release(1, ns());
        // Vehicle 2 is still on NS → EW must remain blocked.
        assert!(!m.request(3, ew()));
    }

    #[test]
    fn manager_is_empty_after_all_releases() {
        let mut m = ReservationManager::new();
        m.request(1, ns());
        m.request(2, ne());
        m.release(1, ns());
        m.release(2, ne());
        assert!(m.is_empty());
    }
}
