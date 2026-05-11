use events::*;
use std::time::Duration;

fn main() {
    println!("{}", Event::Remainder("Go to the doctor").notify());
    println!(
        "{}",
        Event::Registration(Duration::from_secs(49094)).notify()
    );
    println!("{}", Event::Appointment("Go to the doctor").notify());
    println!("{}", Event::Holiday.notify());
}
