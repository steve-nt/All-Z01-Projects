use diamond_creation::get_diamond;

fn main() {
    println!("{:?}", get_diamond('A'));
    println!("{:?}", get_diamond('C'));
    for line in get_diamond('C') {
        println!("{}", line);
    }
}

