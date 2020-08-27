use std::io;
use std::{error::Error, thread};
use std::collections::HashMap;

use colored::*;
use serde_json as json;
use signal_hook::{iterator::Signals, SIGINT, SIGTERM};

const DEFAULT_COLOR: Color = Color::White;

fn main() -> Result<(), Box<dyn Error>> {
    handle_signals()?;

    loop {
        let mut input = String::new();
        if io::stdin().read_line(&mut input)? == 0 {
            return Ok(())
        }
        display(input.trim());
    }
}

fn display(line: &str) {
    let result: Result<HashMap<String, json::Value>, json::error::Error> = json::from_str(line);
    if result.is_err() {
        println!("{}\n", line);
        return;
    }

    let parsed = result.unwrap();

    let level = get_level(&parsed);
    let (key_color, txt_color) = get_colors(level);
    let fields = get_fields(&parsed);

    for key in fields {
        // Unpack string to remove quotes
        let value = match parsed.get(&key) {
            Some(v) => match v {
                json::Value::String(s) => format!("{}", s),
                _ => format!("{}", v),
            },
            None => continue, // this should never happen
        };
        println!("{}{}{}",
            key.color(key_color),
            ": ".color(key_color),
            format!("{}", value).color(txt_color),
        );
    }
    println!("")
}

fn get_level(m: &HashMap<String, json::Value>) -> Option<String> {
    for key in vec!["level", "lvl", "lev", "l", "type"] {
        match m.get(key) {
            Some(value) =>
                match value.as_str() {
                    Some(v) => return Some(String::from(v)),
                    None => return None
                }
            None => continue,
        }
    }
    return None;
}

fn get_colors(level: Option<String>) -> (Color, Color) {
    if level.is_none() {
        return (DEFAULT_COLOR, DEFAULT_COLOR);
    }
    match level.unwrap().to_lowercase().as_str() {
        "debug" | "dbg" | "d" =>
            return (Color::Magenta, DEFAULT_COLOR),
        "info" | "inf" | "i" =>
            return (Color::Cyan, DEFAULT_COLOR),
        "warning" | "warn" | "wrn" | "w" =>
            return (Color::Yellow, DEFAULT_COLOR),
        "error" | "err" | "e" =>
            return (Color::Red, DEFAULT_COLOR),
        "fatal" | "f" =>
            return (Color::Red, DEFAULT_COLOR),
        _ =>
            return (DEFAULT_COLOR, DEFAULT_COLOR),
    }
}

fn get_fields(m: &HashMap<String, json::Value>) -> Vec<String> {
    let first_fields = ["time"];
    let last_fields = ["message"];
    let blacklisted_fields = ["level", "type", "lineno", "function", "env", "tag"];

    let mut fields: Vec<String> = vec![];

    for key in first_fields.iter() {
        let k = key.to_string();
        if m.contains_key(&k) {
            fields.push(k);
        }
    }

    for key in m.keys() {
        if blacklisted_fields.contains(&key.as_str()) ||
            first_fields.contains(&key.as_str()) ||
            last_fields.contains(&key.as_str()) {
                continue
        }
        fields.push(key.to_string());
    }

    for key in last_fields.iter() {
        let k = key.to_string();
        if m.contains_key(&k) {
            fields.push(k);
        }
    }

    return fields;
}

// Catch signals and do nothing - program should quit only on
// the end of user input.
fn handle_signals() -> Result<(), Box<dyn Error>> {
    let signals = Signals::new(&[SIGINT, SIGTERM])?;
    thread::spawn(move || {
        for _ in signals.forever() {}
    });
    return Ok(());
}
