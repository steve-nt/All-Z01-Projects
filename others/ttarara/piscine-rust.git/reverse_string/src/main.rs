use reverse_string::*;

fn main() {
    println!("{}", rev_str("Hello, world!"));
    println!("{}", rev_str("Hello, my name is Roman"));
    println!("{}", rev_str("I have a nice car!"));
    println!("{}", rev_str("How old are You"));
    println!("{}", rev_str("ex: this is an example água"));
}

/*
And its output:

$ cargo run
!dlrow ,olleH
namoR si eman ym ,olleH
!rac ecin a evah I
uoY era dlo woH
augá elpmaxe na si siht :xe
$ */