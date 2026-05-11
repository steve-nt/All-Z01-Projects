use std::collections::HashMap;


#[inline]
fn coerce_map<V>(m: HashMap<impl Into<String>, V>) -> HashMap<String, V> {
    m.into_iter().map(|(k, v)| (k.into(), v)).collect()
}

#[derive(Debug, Clone, PartialEq)]
pub struct Mall {
    pub name: String,
    pub guards: HashMap<String, Guard>,
    pub floors: HashMap<String, Floor>,
}

impl Mall {
    pub fn new(
        name: impl Into<String>,
        guards: HashMap<impl Into<String>, Guard>,
        floors: HashMap<impl Into<String>, Floor>,
    ) -> Self {
        Self {
            name: name.into(),
            guards: coerce_map(guards),
            floors: coerce_map(floors),
        }
    }

    pub fn change_name(&mut self, new_name: impl Into<String>) {
        self.name = new_name.into();
    }

    pub fn hire_guard(&mut self, name: impl Into<String>, guard: Guard) {
        self.guards.insert(name.into(), guard);
    }

    pub fn fire_guard(&mut self, name: impl Into<String>) {
        self.guards.remove(&name.into());
    }
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub struct Guard {
    pub age: u32,
    pub years_experience: u32,
}

#[derive(Debug, Clone, PartialEq)]
pub struct Floor {
    pub stores: HashMap<String, Store>,
    pub size_limit: u64,
}

impl Floor {
    pub fn new(stores: HashMap<impl Into<String>, Store>, size_limit: u64) -> Self {
        Self {
            stores: coerce_map(stores),
            size_limit,
        }
    }

    pub fn replace_store(&mut self, store: impl Into<String>, with: Store) {
        self.stores.entry(store.into()).and_modify(|v| *v = with);
    }

    pub fn add_store(&mut self, name: impl Into<String>, store: Store) -> Result<(), ()> {
        let has_space = self.size_limit
            >= self.stores.values().map(|s| s.square_meters).sum::<u64>() + store.square_meters;

        if has_space {
            self.stores.insert(name.into(), store);
            Ok(())
        } else {
            Err(())
        }
    }

    pub fn remove_store(&mut self, name: impl Into<String>) {
        self.stores.remove(&name.into());
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct Store {
    pub employees: HashMap<String, Employee>,
    pub square_meters: u64,
}

impl Store {
    pub fn new(employees: HashMap<impl Into<String>, Employee>, square_meters: u64) -> Self {
        Self {
            employees: coerce_map(employees),
            square_meters,
        }
    }

    pub fn hire_employee(&mut self, name: impl Into<String>, employee: Employee) {
        self.employees.insert(name.into(), employee);
    }

    pub fn fire_employee(&mut self, name: impl Into<String>) {
        self.employees.remove(&name.into());
    }

    pub fn expand(&mut self, by: u64) {
        self.square_meters += by;
    }
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub struct Employee {
    pub age: u32,
    pub working_hours: (u32, u32),
    pub salary: f64,
}

impl Employee {
    pub fn birthday(&mut self) {
        self.age += 1;
    }

    pub fn change_workload(&mut self, from: u32, to: u32) {
        self.working_hours = (from, to);
    }

    pub fn raise(&mut self, amount: f64) {
        self.salary += amount;
    }

    pub fn cut(&mut self, amount: f64) {
        self.salary -= amount;
    }
}

// ------------ required functions ------------

pub fn biggest_store(mall: &Mall) -> (&String, &Store) {
    mall.floors
        .values()
        .flat_map(|floor| floor.stores.iter())
        .max_by_key(|(_, store)| store.square_meters)
        .expect("mall has at least one store")
}

pub fn highest_paid_employee(mall: &Mall) -> Vec<(&String, &Employee)> {
    let mut max_salary: Option<f64> = None;

    for floor in mall.floors.values() {
        for store in floor.stores.values() {
            for employee in store.employees.values() {
                max_salary = Some(match max_salary {
                    Some(current_max) if current_max >= employee.salary => current_max,
                    _ => employee.salary,
                });
            }
        }
    }

    let Some(max_salary) = max_salary else {
        return Vec::new();
    };

    let mut result = Vec::new();
    for floor in mall.floors.values() {
        for store in floor.stores.values() {
            for (name, employee) in store.employees.iter() {
                if (employee.salary - max_salary).abs() < f64::EPSILON {
                    result.push((name, employee));
                }
            }
        }
    }
    result
}

pub fn nbr_of_employees(mall: &Mall) -> usize {
    let employee_count: usize = mall
        .floors
        .values()
        .map(|floor| {
            floor
                .stores
                .values()
                .map(|store| store.employees.len())
                .sum::<usize>()
        })
        .sum();

    employee_count + mall.guards.len()
}

pub fn check_for_securities(mall: &mut Mall, mut new_guards: HashMap<String, Guard>) {
    // total floor size as sum of size_limit, per subject text
    let total_floor_size: u64 = mall.floors.values().map(|floor| floor.size_limit).sum();
    let required_guards = ((total_floor_size + 199) / 200) as usize;

    while mall.guards.len() < 4 {
        if let Some((name, guard)) = new_guards.drain().next() {
            mall.hire_guard(name, guard);
        } else {
            break; // no more guards available
        }
    }
}

pub fn cut_or_raise(mall: &mut Mall) {
    for floor in mall.floors.values_mut() {
        for store in floor.stores.values_mut() {
            for employee in store.employees.values_mut() {
                let hours = employee.working_hours.1.saturating_sub(employee.working_hours.0);
                let delta = employee.salary * 0.10;
                if hours >= 10 {
                    employee.raise(delta);
                } else {
                    employee.cut(delta);
                }
            }
        }
    }
}