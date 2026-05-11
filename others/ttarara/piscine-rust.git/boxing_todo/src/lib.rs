mod err;

pub use err::{ParseErr, ReadErr};

use json;
use std::error::Error;
use std::fs;

#[derive(Debug, Eq, PartialEq)]
pub struct Task {
    pub id: u32,
    pub description: String,
    pub level: u32,
}

#[derive(Debug, Eq, PartialEq)]
pub struct TodoList {
    pub title: String,
    pub tasks: Vec<Task>,
}

impl TodoList {
    pub fn get_todo(path: &str) -> Result<TodoList, Box<dyn Error>> {
        // 1. Read file
        let contents = fs::read_to_string(path)
            .map_err(|e| Box::new(ReadErr::new(e)) as Box<dyn Error>)?;

        // 2. Parse JSON
        let parsed = json::parse(&contents)
            .map_err(|e| Box::new(ParseErr::Malformed(Box::new(e))) as Box<dyn Error>)?;

        // 3. Extract title
        let title = parsed["title"]
            .as_str()
            .unwrap_or_default()
            .to_string();

        // 4. Extract tasks
        let tasks_json = &parsed["tasks"];

        // If tasks is missing or empty → ParseErr::Empty
        if !tasks_json.is_array() || tasks_json.len() == 0 {
            return Err(Box::new(ParseErr::Empty));
        }

        let mut tasks = Vec::new();

        for t in tasks_json.members() {
            let id = t["id"].as_u32().unwrap_or(0);
            let description = t["description"].as_str().unwrap_or_default().to_string();
            let level = t["level"].as_u32().unwrap_or(0);

            tasks.push(Task {
                id,
                description,
                level,
            });
        }

        Ok(TodoList { title, tasks })
    }
}