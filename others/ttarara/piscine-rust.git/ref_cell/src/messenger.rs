use std::cell::RefCell;
use std::rc::Rc;

pub struct Tracker {
    pub messages: RefCell<Vec<String>>,
    pub value: RefCell<usize>,
    pub max: usize,
}

impl Tracker {
    pub fn new(max: usize) -> Self {
        Tracker {
            messages: RefCell::new(Vec::new()),
            value: RefCell::new(0),
            max,
        }
    }

    pub fn set_value<T>(&self, rc: &Rc<T>) {
        let count = Rc::strong_count(rc);
        if count > self.max {
            self.messages
                .borrow_mut()
                .push("Error: You can't go over your quota!".to_string());
            return;
        }

        let pct = count * 100 / self.max;
        if pct > 70 {
            self.messages.borrow_mut().push(format!(
                "Warning: You have used up over {}% of your quota!",
                pct
            ));
        }

        *self.value.borrow_mut() = count;
    }

    pub fn peek<T>(&self, rc: &Rc<T>) {
        let count = Rc::strong_count(rc);
        let pct = count * 100 / self.max;
        self.messages.borrow_mut().push(format!(
            "Info: This value would use {}% of your quota",
            pct
        ));
    }
}
