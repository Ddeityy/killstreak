extern crate libc;

use edit::{edit, EditOptions, TickRange};
use serde::{Deserialize, Serialize};
use serde_json;
use std::ffi::{CStr, CString};
use std::fs;
use std::path::PathBuf;
use tf_demo_parser::demo::header::Header;
use tf_demo_parser::demo::parser::analyser::MatchState;
use tf_demo_parser::Demo;
use tf_demo_parser::DemoParser;

#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
#[repr(C)]
struct JsonDemo {
    header: Header,
    state: MatchState,
}

#[no_mangle]
pub extern "C" fn parse_demo(demo_path: *const libc::c_char) -> *const libc::c_char {
    let file_path = unsafe { CStr::from_ptr(demo_path) };
    let file_path_str = file_path.to_str().unwrap();
    let not_found_error: Vec<u8> = String::from("File not found").into();
    let file = fs::read(PathBuf::from(file_path_str)).unwrap_or(not_found_error.clone());
    if file == not_found_error {
        return CString::new(file).unwrap().into_raw();
    }
    let demo = Demo::new(&file);

    let parser = DemoParser::new(demo.get_stream());
    let incomplete_demo_error: Vec<u8> = String::from("Incomplete demo").into();
    let (header, state) = match parser.parse() {
        Ok((header, state)) => (header, state),
        Err(_) => return CString::new(incomplete_demo_error).unwrap().into_raw(),
    };
    let demo = JsonDemo { header, state };
    let result = serde_json::to_string(&demo).unwrap().as_bytes().to_owned();
    CString::new(result).unwrap().into_raw()
}

#[no_mangle]
pub extern "C" fn cut_demo(c_demo_path: *const libc::c_char, c_start_tick: *const libc::c_char) {
    let file_path = unsafe { CStr::from_ptr(c_demo_path) };
    let file_path_str = file_path.to_str().unwrap();
    let mut demo_path = PathBuf::from(file_path_str);

    println!("{:?}", &demo_path);

    let start_tick = unsafe { CStr::from_ptr(c_start_tick) };
    let start_tick_str = start_tick.to_str().unwrap();
    let start_tick_i32: u32 = start_tick_str.parse().unwrap();

    println!("{:?}", &start_tick_i32);

    let options = EditOptions {
        unlock_pov: false,
        cut: Some(TickRange {
            from: (start_tick_i32 - 500).into(),
            to: (start_tick_i32 + 3000).into(),
        }),
        ..EditOptions::default()
    };
    let input = fs::read(&demo_path).unwrap();
    let output = edit(&input, options);
    let demo_file_name = demo_path.file_name().unwrap();
    demo_path.set_file_name(format!("cut_{}", &demo_file_name.to_str().unwrap()));
    fs::write(PathBuf::from(&demo_path), output).unwrap();
    ()
}
