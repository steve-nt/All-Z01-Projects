#[derive(Debug, PartialEq)]
pub enum Role {
    CEO,
    Manager,
    Worker,
}

impl From<&str> for Role {
    fn from(s: &str) -> Self {
        match s {
            "CEO" => Role::CEO,
            "Manager" => Role::Manager,
            "Worker" | "Normal Worker" => Role::Worker,
            _ => Role::Worker,
        }
    }
}

pub type Link = Option<Box<Worker>>;

#[derive(Debug)]
pub struct WorkEnvironment {
    pub grade: Link,
}

#[derive(Debug)]
pub struct Worker {
    pub role: Role,
    pub name: String,
    pub next: Link,
}

impl WorkEnvironment {
    pub fn new() -> Self {
        WorkEnvironment { grade: None }
    }

    pub fn add_worker(&mut self, name: &str, role: &str) {
        let worker = Worker {
            role: Role::from(role),
            name: name.to_string(),
            next: self.grade.take(),
        };
        self.grade = Some(Box::new(worker));
    }

    pub fn remove_worker(&mut self) -> Option<String> {
        let head = self.grade.take()?;
        let name = head.name;
        self.grade = head.next;
        Some(name)
    }

    pub fn last_worker(&self) -> Option<(String, Role)> {
        self.grade.as_ref().map(|w| {
            let role = match &w.role {
                Role::CEO => Role::CEO,
                Role::Manager => Role::Manager,
                Role::Worker => Role::Worker,
            };
            (w.name.clone(), role)
        })
    }
}
