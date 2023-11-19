use serde;
use serde::{Deserialize, Serialize};
use serde_json;
use std::env;
use std::fs;
use tf_demo_parser::demo::header::Header;
use tf_demo_parser::demo::parser::analyser::MatchState;
use tf_demo_parser::Demo;
use tf_demo_parser::DemoParser;

#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
struct JsonDemo {
    header: Header,
    state: MatchState,
}
fn main() {
    let args: Vec<_> = env::args().collect();

    let path = args[1].clone();
    let file = fs::read(path).unwrap();
    let demo = Demo::new(&file);

    let parser = DemoParser::new(demo.get_stream());
    let (header, state) = parser.parse().unwrap();
    let demo = JsonDemo { header, state };
    println!("{}", serde_json::to_string(&demo).unwrap());
}
