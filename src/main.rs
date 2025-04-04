use cliconf::Parse;
use rand::{rng, seq::SliceRandom};
use serde::{de::DeserializeOwned, Deserialize, Serialize};
use serde_json::Value;
use std::{env, fs::File, io::Read, process::Command};

fn merge(a: &mut Value, b: &Value) {
    match (a, b) {
        (Value::Object(a), Value::Object(b)) => {
            for (k, v) in b {
                merge(a.entry(k.clone()).or_insert(Value::Null), v);
            }
        }
        (a, b) => *a = b.clone(),
    }
}

trait Merge: DeserializeOwned + Serialize {
    fn merge(&mut self, value: &Value) -> Result<(), serde_json::Error> {
        let mut a = serde_json::to_value(&self)?;
        merge(&mut a, &value);
        *self = serde_json::from_value(a)?;
        Ok(())
    }
}

impl Merge for Flags {}

#[derive(Parse, Deserialize, Serialize)]
#[serde(default)]
struct Flags {
    #[cliconf(shorthand = 'h')]
    help: bool,

    #[cliconf(shorthand = 'v')]
    version: bool,

    #[cliconf(shorthand = 'H')]
    host: String,

    #[cliconf(shorthand = 'd')]
    dir: String,

    #[cliconf(shorthand = 'r')]
    resume: bool,

    #[cliconf(shorthand = 'D')]
    disown: bool,

    #[cliconf(shorthand = 'a')]
    album: bool,

    #[cliconf(shorthand = 's')]
    shuffle: bool,
}

impl Default for Flags {
    fn default() -> Self {
        Flags {
            help: false,
            version: false,
            host: "localhost".into(),
            dir: "/srv".into(),
            resume: false,
            disown: false,
            album: false,
            shuffle: false,
        }
    }
}

fn main() {
    let mut flags = Flags::default();

    if let Some(home_dir) = dirs::home_dir() {
        let path = home_dir.join(".config/bop/config.json");
        if let Ok(mut file) = File::open(path) {
            let mut data = String::new();
            if let Ok(_) = file.read_to_string(&mut data) {
                flags.merge(&serde_json::from_str(&data).unwrap()).unwrap();
            }
        }
    }

    flags.parse_env(env::vars().collect());
    flags.parse_args(env::args().skip(1).collect());

    if flags.version {
        let version = env!("CARGO_PKG_VERSION");
        println!("bop {version}");
        return;
    }

    if flags.help {
        // TODO: re-implement usage
        // let width = match term_size::dimensions() {
        //     Some((w, _)) => w,
        //     None => 75,
        // };
        // cliconf::usage::generate(&flags, width, &mut io::stdout()).expect("Failed to print usage");
        return;
    }

    let command = match flags.album {
        true => "fd . --type directory | fzf",
        false => "fzf",
    };

    // Use fzf to locate a file
    Command::new("ssh")
        .args([
            "-t",
            &flags.host,
            &format!("cd {} && {command} > /tmp/bop", flags.dir),
        ])
        .status()
        .expect("SSH command failed");

    // Retrieve the located file
    let output = Command::new("ssh")
        .args([&flags.host, "cat /tmp/bop"])
        .output()
        .expect("SSH command failed");
    let path = match output.status.success() {
        true => String::from_utf8_lossy(&output.stdout).trim().to_owned(),
        false => panic!("Failed to get result!"),
    };

    let mut paths = vec![];
    if flags.album {
        // List the file paths in an album
        let output = Command::new("ssh")
            .args([&flags.host, &format!("ls \"{}/{path}\"", flags.dir)])
            .output()
            .expect("SSH command failed");
        let path_list = match output.status.success() {
            true => String::from_utf8_lossy(&output.stdout).trim().to_owned(),
            false => panic!("Failed to get file list!"),
        };
        for subpath in path_list.split("\n") {
            paths.push(format!("{path}/{subpath}"));
        }
    } else {
        paths.push(path.to_string());
    }
    if flags.shuffle {
        paths.shuffle(&mut rng());
    }
    let paths = paths;

    let mut mpv_args: Vec<String> = paths
        .iter()
        .map(|path| format!("sftp://{}:{}/{path}", flags.host, flags.dir))
        .collect();
    mpv_args.push(match flags.resume {
        true => "--save-position-on-quit".to_string(),
        false => "--no-resume-playback".to_string(),
    });
    let mpv_args = mpv_args;

    if flags.disown {
        Command::new("mpv")
            .args(mpv_args)
            .spawn()
            .expect("Failed to start mpv");
    } else {
        Command::new("mpv")
            .args(mpv_args)
            .status()
            .expect("Failed to start mpv");
    }
}
