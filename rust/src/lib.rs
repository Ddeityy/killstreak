extern crate libc;

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
