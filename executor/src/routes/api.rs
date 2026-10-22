use crate::config::language_config::LanguageConfig;
use rocket::request::Form;
use rocket::State;
use rocket_contrib::json::JsonValue;
use std::fs;
use std::io::prelude::*;
use std::process::{Command, Stdio};

const TEMP_FILENAME: &'static str = "Submission";

#[derive(FromForm)]
pub struct ExecForm {
    code: String,
    language: String,
    timeout: Option<u8>,
    memory: Option<u32>,
    stdin: Option<String>,
}

fn compile(compiler: &str, flag: &Option<String>) -> Result<(), String> {
    println!("Compiling");
    let mut compiler = Command::new(&compiler);

    let output = if let Some(flag) = flag {
        compiler
            .args({ flag.to_string().split_whitespace() })
            .output()
            .expect("Failed to execute process")
    } else {
        compiler.output().expect("Failed to execute process")
    };

    if output.status.success() {
        Ok(())
    } else {
        Err(String::from_utf8_lossy(&output.stderr).to_string())
    }
}

fn run(command: &str, stdin: &Option<String>) -> Result<String, String> {
    println!("Running");
    let process = Command::new("firejail")
        .arg("--noprofile")
        .arg("--quiet")
        .arg("--seccomp")
        .arg("--timeout=0:0:5")
        // .arg("--read-only=\"$(pwd)/sandbox\"")
        .args({ command.to_string().split_whitespace() })
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .expect("Failed to execute process");
    let mut stdout = String::new();
    let mut stderr = String::new();

    if let Some(stdin) = stdin {
        process
            .stdin
            .unwrap()
            .write_all(stdin.as_bytes())
            .expect("Failed to write in stdin");
    }
    process
        .stdout
        .unwrap()
        .read_to_string(&mut stdout)
        .expect("Failed to write stdout to buffer");
    process
        .stderr
        .unwrap()
        .read_to_string(&mut stderr)
        .expect("Failed to write stderr to buffer");

    return if !stderr.is_empty() {
        Err(stderr)
    } else {
        Ok(stdout)
    };
}

fn to_file(filename: &str, extension: &str, code: &str) -> Result<(), &'static str> {
    match std::fs::File::create(&format!("sandbox/{}.{}", filename, extension)) {
        Ok(mut file) => match file.write_all(code.as_bytes()) {
            Ok(_) => Ok(()),
            Err(_) => Err("Error when trying to write to file"),
        },
        Err(_) => Err("Error when trying to create the file"),
    }
}

#[post("/exec", data = "<exec_form>")]
pub fn exec(exec_form: Form<ExecForm>, language_config: State<LanguageConfig>) -> JsonValue {
    if let Some(config) = language_config.of(&exec_form.language) {
        if let Ok(_) = to_file(TEMP_FILENAME, &config.extension, &exec_form.code) {
            if let Some(compiler) = &config.compiler {
                match compile(compiler, &config.flag) {
                    Err(error) => {
                        return json!({
                            "stderr": error,
                            "success": true,
                        })
                    }
                    Ok(_) => {}
                }
            }
            let process = run(&config.command, &exec_form.stdin);
            fs::remove_file(&format!("sandbox/{}.{}", TEMP_FILENAME, &config.extension)).expect("Can't remove file");
            return match process {
                Ok(result) => json!({
                    "stdout": result,
                    "success": true,
                }),
                Err(error) => json!({
                    "stderr": error,
                    "success": false,
                }),
            };
        } else {
            return json!({
                "message": "Error when writing the files",
                "success": false,
            });
        }
    } else {
        return json!({
            "message": "Invalid language",
            "success": false,
        });
    }
}
