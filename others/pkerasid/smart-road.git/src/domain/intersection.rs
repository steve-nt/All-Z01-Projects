use std::collections::HashMap;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum EntryDir {
    North,
    South,
    East,
    West,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum TurnDir {
    Straight,
    Left,
}

// RouteId = (where you enter, which way you turn)
pub type RouteId = (EntryDir, TurnDir);
pub type TileId = (u8, u8); // (col, row), 0-indexed, (0,0) = top-left

#[derive(Debug, Clone)]
pub struct ConflictEntry {
    pub other_route: RouteId,
    pub shared_tiles: Vec<TileId>,
}

/// Tile occupancy for each route in the inner 4×4 conflict zone.
///
/// Inner tile `(col, row)` with `(0,0)` = top-left.  Each tile maps to one
/// lane-width cell (40 px); the inner zone spans outer lanes 1–4 in both axes.
///
/// ```text
///  col:  0     1     2     3
/// row 0: (0,0) (1,0) (2,0) (3,0)
/// row 1: (0,1) (1,1) (2,1) (3,1)
/// row 2: (0,2) (1,2) (2,2) (3,2)
/// row 3: (0,3) (1,3) (2,3) (3,3)
/// ```
///
/// Right-turn routes use the outer corner tiles and are excluded from this map.
#[must_use]
pub fn build_route_tiles() -> HashMap<RouteId, Vec<TileId>> {
    [
        // NS  - North entering, going South (col 0)
        (
            (EntryDir::North, TurnDir::Straight),
            vec![(0, 0), (0, 1), (0, 2), (0, 3)],
        ),
        // SN  - South entering, going North (col 3)
        (
            (EntryDir::South, TurnDir::Straight),
            vec![(3, 3), (3, 2), (3, 1), (3, 0)],
        ),
        // EW  - East entering, going West (row 0)
        (
            (EntryDir::East, TurnDir::Straight),
            vec![(3, 0), (2, 0), (1, 0), (0, 0)],
        ),
        // WE  - West entering, going East (row 3)
        (
            (EntryDir::West, TurnDir::Straight),
            vec![(0, 3), (1, 3), (2, 3), (3, 3)],
        ),
        // NE  - North entering, turning Left → exits East (col 1 → row 2)
        (
            (EntryDir::North, TurnDir::Left),
            vec![(1, 0), (1, 1), (1, 2), (2, 2), (3, 2)],
        ),
        // SW  - South entering, turning Left → exits West (col 2 → row 1)
        (
            (EntryDir::South, TurnDir::Left),
            vec![(2, 3), (2, 2), (2, 1), (1, 1), (0, 1)],
        ),
        // ES  - East entering, turning Left → exits South (row 1 → col 1)
        (
            (EntryDir::East, TurnDir::Left),
            vec![(3, 1), (2, 1), (1, 1), (1, 2), (1, 3)],
        ),
        // WN  - West entering, turning Left → exits North (row 2 → col 2)
        (
            (EntryDir::West, TurnDir::Left),
            vec![(0, 2), (1, 2), (2, 2), (2, 1), (2, 0)],
        ),
    ]
    .into_iter()
    .collect()
}

/// Built once at startup. Never recomputed.
/// Encodes: "if I am on route X, which routes conflict with me, and on which tiles?"
/// Derived directly from the tile diagram.
#[must_use]
pub fn build_conflict_map() -> HashMap<RouteId, Vec<ConflictEntry>> {
    let route_tiles = build_route_tiles();

    // Two routes conflict if they share at least one tile.
    let routes: Vec<RouteId> = route_tiles.keys().copied().collect();
    let mut conflict_map: HashMap<RouteId, Vec<ConflictEntry>> = HashMap::new();

    for &route_a in &routes {
        let tiles_a = &route_tiles[&route_a];
        let mut entries = Vec::new();

        for &route_b in &routes {
            if route_a == route_b {
                continue;
            }

            let tiles_b = &route_tiles[&route_b];

            let shared: Vec<TileId> = tiles_a
                .iter()
                .filter(|t| tiles_b.contains(t))
                .copied()
                .collect();

            if !shared.is_empty() {
                entries.push(ConflictEntry {
                    other_route: route_b,
                    shared_tiles: shared,
                });
            }
        }

        conflict_map.insert(route_a, entries);
    }

    conflict_map
}

#[cfg(test)]
mod tests {
    use super::*;
    use EntryDir::*;
    use TurnDir::*;

    fn ns() -> RouteId {
        (North, Straight)
    }
    fn sn() -> RouteId {
        (South, Straight)
    }
    fn ew() -> RouteId {
        (East, Straight)
    }

    #[test]
    fn map_contains_all_eight_routes() {
        assert_eq!(build_conflict_map().len(), 8);
    }

    #[test]
    fn no_route_conflicts_with_itself() {
        let map = build_conflict_map();
        for (route, entries) in &map {
            assert!(!entries.iter().any(|e| e.other_route == *route));
        }
    }

    #[test]
    fn conflict_is_symmetric() {
        let map = build_conflict_map();
        for (route_a, entries) in &map {
            for entry in entries {
                let route_b = entry.other_route;
                assert!(
                    map[&route_b].iter().any(|e| e.other_route == *route_a),
                    "{route_a:?} lists {route_b:?} as a conflict but not vice versa",
                );
            }
        }
    }

    #[test]
    fn opposite_straights_do_not_conflict() {
        // NS occupies col 0, SN occupies col 3 — no shared tiles.
        let map = build_conflict_map();
        assert!(!map[&ns()].iter().any(|e| e.other_route == sn()));
        assert!(!map[&sn()].iter().any(|e| e.other_route == ns()));
    }

    #[test]
    fn crossing_straights_conflict_at_corner_tile() {
        // NS (col 0) and EW (row 0) share exactly tile (0,0).
        let map = build_conflict_map();
        let entry = map[&ns()].iter().find(|e| e.other_route == ew());
        assert!(entry.is_some(), "NS should conflict with EW");
        assert!(entry.unwrap().shared_tiles.contains(&(0, 0)));
    }

    #[test]
    fn ns_has_four_conflicts() {
        // NS (col 0) shares one tile with each of: EW, WE, SW (South,Left), WN (West,Left).
        // It does not share tiles with SN, NE, or ES.
        let map = build_conflict_map();
        assert_eq!(map[&ns()].len(), 4);
    }

    #[test]
    fn all_conflict_entries_have_nonempty_shared_tiles() {
        let map = build_conflict_map();
        for entries in map.values() {
            for entry in entries {
                assert!(!entry.shared_tiles.is_empty());
            }
        }
    }
}
