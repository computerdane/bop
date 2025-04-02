use std::{env, process::Command};

use cliconf::{Flag, FlagValue, Flags};

fn init_flags<'a>() -> Flags<'a> {
    let mut flags = Flags::new();
    flags.add(Flag {
        name: "host",
        shorthand: Some('h'),
        default_value: FlagValue::String("localhost".to_string()),
        description: Some("The SSH host you want to connect to."),
    });
    flags.add(Flag {
        name: "dir",
        shorthand: Some('d'),
        default_value: FlagValue::String("/srv".to_string()),
        description: Some("The directory to scan for files."),
    });
    flags.add(Flag {
        name: "resume",
        shorthand: Some('r'),
        default_value: FlagValue::Bool(false),
        description: Some("If true, will enable resuming playback."),
    });
    flags.add(Flag {
        name: "background",
        shorthand: Some('b'),
        default_value: FlagValue::Bool(false),
        description: Some("If true, will run mpv as a background process."),
    });
    flags.add(Flag {
        name: "album",
        shorthand: Some('a'),
        default_value: FlagValue::Bool(false),
        description: Some("If true, will look for media folders and play files from the selected folder in alphabetical order."),
    });
    flags.add(Flag {
        name: "shuffle",
        shorthand: Some('s'),
        default_value: FlagValue::Bool(false),
        description: Some("If true, will shuffle an album. Only works when --album is true."),
    });
    flags.add_home_config_file(".config/bop/config.json");
    flags
        .load(
            &env::vars().collect(),
            &env::args().collect::<Vec<String>>()[1..].to_vec(),
        )
        .expect("Failed to load flags");
    flags
}

fn main() {
    let flags = init_flags();

    let album = flags.get_bool("album");
    let command = match album {
        true => "fd . --type directory | fzf",
        false => "fzf",
    };

    let host = flags.get_string("host");
    let dir = flags.get_string("dir");

    // Use fzf to locate a file
    Command::new("ssh")
        .args(["-t", host, &format!("cd {dir} && {command} > /tmp/bop")])
        .status()
        .expect("SSH command failed");

    // Retrieve the located file
    let output = Command::new("ssh")
        .args([host, "cat /tmp/bop"])
        .output()
        .expect("SSH command failed");
    let path = match output.status.success() {
        true => String::from_utf8_lossy(&output.stdout).trim().to_owned(),
        false => panic!("Failed to get result!"),
    };

    let mut paths = vec![];
    if *album {
        // List the file paths in an album
        let output = Command::new("ssh")
            .args([host, &format!("ls \"{dir}/{path}\"")])
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
    let paths = paths;

    let mut mpv_args: Vec<String> = paths
        .iter()
        .map(|path| format!("sftp://{host}:{dir}/{path}"))
        .collect();
    mpv_args.push(match flags.get_bool("resume") {
        true => "--save-position-on-quit".to_string(),
        false => "--no-resume-playback".to_string(),
    });
    let mpv_args = mpv_args;

    if *flags.get_bool("background") {
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
