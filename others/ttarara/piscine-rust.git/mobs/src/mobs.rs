use std::collections::{HashMap, HashSet};

pub mod boss;
pub mod member;

pub use boss::Boss;
pub use member::{Member, Role};

#[derive(Debug, Clone, PartialEq)]
pub struct Mob {
    pub name: String,
    pub boss: Boss,
    pub members: HashMap<String, Member>,
    pub cities: HashSet<String>,
    pub wealth: u64,
}

impl Mob {
    pub fn recruit(&mut self, info: (&str, u32)) {
        let (name, age) = info;
        self.members.insert(
            name.to_owned(),
            Member {
                role: Role::Associate,
                age,
            },
        );
    }

    fn combat_power(&self) -> u32 {
        self.members.values().map(|m| m.role.power()).sum()
    }

    fn remove_youngest_members(&mut self) {
        let Some(min_age) = self.members.values().map(|m| m.age).min() else {
            return;
        };
        self.members.retain(|_, m| m.age != min_age);
    }

    pub fn attack(&mut self, other: &mut Mob) {
        let self_power = self.combat_power();
        let other_power = other.combat_power();

        // Loser is the one with least power; on draw attacker loses.
        let (winner, loser) = if self_power > other_power {
            (self, other)
        } else {
            (other, self)
        };

        loser.remove_youngest_members();

        if loser.members.is_empty() {
            winner.wealth = winner.wealth.saturating_add(loser.wealth);
            loser.wealth = 0;

            for city in loser.cities.drain() {
                winner.cities.insert(city);
            }
        }
    }

    pub fn steal(&mut self, target: &mut Mob, amount: u64) {
        let stolen = amount.min(target.wealth);
        target.wealth -= stolen;
        self.wealth += stolen;
    }

    pub fn conquer_city(&mut self, mobs: &[&Mob], city: String) {
        if mobs.iter().any(|m| m.cities.contains(&city)) {
            return;
        }
        self.cities.insert(city);
    }
}

