use edit::{edit, EditOptions, TickRange};
use std::env;
use std::fs;
use std::path::PathBuf;

fn main() {
    let args: Vec<_> = env::args().collect();

    let mut demo_in: PathBuf = args[1].clone().into();
    let tick_from: u32;
    let tick_to: u32;
    if args.len() == 2 {
        let d = demo_in.file_stem().unwrap().to_str().unwrap().split("_");
        let c: Vec<_> = d.collect();
        tick_from = c[1].parse().unwrap();
        tick_to = tick_from.clone()
    } else {
        tick_from = args[2].clone().parse().unwrap();
        tick_to = args[2].clone().parse().unwrap();
    }

    let options = EditOptions {
        unlock_pov: false,
        cut: Some(TickRange {
            from: (tick_from - 500).into(),
            to: (tick_to + 3000).into(),
        }),
        ..EditOptions::default()
    };
    let input = fs::read(&demo_in).unwrap();
    let output = edit(&input, options);
    let demo_file_name = demo_in.file_name().unwrap();
    demo_in.set_file_name(format!("cut_{}", &demo_file_name.to_str().unwrap()));
    fs::write(PathBuf::from(&demo_in), output).unwrap();
}
