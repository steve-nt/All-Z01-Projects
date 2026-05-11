mod mobs;

pub use mobs::{Boss, Member, Mob, Role};

// Re-export modules to satisfy testers that use `member::Role` / `boss::Boss`
pub mod member {
    pub use crate::mobs::member::*;
}

pub mod boss {
    pub use crate::mobs::boss::*;
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::{HashMap, HashSet};

    fn mob_with_members(name: &str, members: Vec<(&str, Role, u32)>) -> Mob {
        let mut map = HashMap::new();
        for (n, r, age) in members {
            map.insert(
                n.to_owned(),
                Member {
                    role: r,
                    age,
                },
            );
        }
        Mob {
            name: name.to_owned(),
            boss: Boss::new("Boss", 50),
            members: map,
            cities: HashSet::new(),
            wealth: 0,
        }
    }

    #[test]
    fn recruit_adds_associate() {
        let mut m = mob_with_members("A", vec![]);
        m.recruit(("New Guy", 22));
        assert_eq!(
            m.members.get("New Guy"),
            Some(&Member {
                role: Role::Associate,
                age: 22
            })
        );
    }

    #[test]
    fn promotion_chain_and_panic() {
        let mut mem = Member {
            role: Role::Associate,
            age: 30,
        };
        mem.get_promotion();
        assert_eq!(mem.role, Role::Soldier);
        mem.get_promotion();
        assert_eq!(mem.role, Role::Caporegime);
        mem.get_promotion();
        assert_eq!(mem.role, Role::Underboss);
    }

    #[test]
    #[should_panic]
    fn underboss_promotion_panics() {
        let mut mem = Member {
            role: Role::Underboss,
            age: 40,
        };
        mem.get_promotion();
    }

    #[test]
    fn attack_draw_attacker_loses_youngest() {
        let mut a = mob_with_members("A", vec![("a1", Role::Soldier, 19)]);
        let mut b = mob_with_members("B", vec![("b1", Role::Soldier, 25)]);
        // draw (2 vs 2): attacker a loses youngest (only member)
        a.attack(&mut b);
        assert!(a.members.is_empty());
        assert_eq!(b.members.len(), 1);
    }

    #[test]
    fn attack_removes_all_youngest_members_of_loser() {
        let mut winner = mob_with_members(
            "W",
            vec![("w1", Role::Underboss, 50)], // power 4
        );
        let mut loser = mob_with_members(
            "L",
            vec![
                ("l1", Role::Associate, 18),
                ("l2", Role::Associate, 18),
                ("l3", Role::Associate, 40),
            ], // power 3
        );
        winner.attack(&mut loser);
        // youngest age is 18 => two removed
        assert_eq!(loser.members.len(), 1);
        assert!(loser.members.contains_key("l3"));
    }

    #[test]
    fn attack_transfers_wealth_and_cities_on_total_defeat() {
        let mut a = mob_with_members("A", vec![("a1", Role::Underboss, 40)]);
        let mut b = mob_with_members("B", vec![("b1", Role::Associate, 18)]);
        b.wealth = 100;
        b.cities.insert("Rome".to_owned());
        b.cities.insert("Milan".to_owned());

        a.wealth = 5;
        a.cities.insert("Naples".to_owned());

        a.attack(&mut b);

        assert_eq!(b.wealth, 0);
        assert!(b.cities.is_empty());
        assert!(b.members.is_empty());

        assert_eq!(a.wealth, 105);
        assert!(a.cities.contains("Rome"));
        assert!(a.cities.contains("Milan"));
        assert!(a.cities.contains("Naples"));
    }

    #[test]
    fn steal_cannot_steal_more_than_target_has() {
        let mut thief = mob_with_members("T", vec![]);
        let mut target = mob_with_members("X", vec![]);
        target.wealth = 30;
        thief.wealth = 10;

        thief.steal(&mut target, 50);
        assert_eq!(target.wealth, 0);
        assert_eq!(thief.wealth, 40);
    }

    #[test]
    fn conquer_city_only_if_unique_among_others() {
        let mut a = mob_with_members("A", vec![]);
        let mut b = mob_with_members("B", vec![]);
        b.cities.insert("Paris".to_owned());

        a.conquer_city(&[&b], "Paris".to_owned());
        assert!(!a.cities.contains("Paris"));

        a.conquer_city(&[&b], "Berlin".to_owned());
        assert!(a.cities.contains("Berlin"));
    }
}

