use std::{env, io, process::Command};

use cliconf::{usage, Flag, FlagValue, Flags};
use rand::{rng, seq::SliceRandom};

fn init_flags() -> Flags {
    let mut flags = Flags::new();
    flags.add(Flag::new("help", FlagValue::Bool(false)).exclude_from_usage());
    flags.add(
        Flag::new("host", FlagValue::String("localhost".into()))
            .shorthand('h')
            .description("The SSH host you want to connect to."),
    );
    flags.add(
        Flag::new("dir", FlagValue::String("/srv".into()))
            .shorthand('d')
            .description("The directory to scan for files."),
    );
    flags.add(
        Flag::new("resume", FlagValue::Bool(false))
            .shorthand('r')
            .description("If true, will enable resuming playback."),
    );
    flags.add(
        Flag::new("background", FlagValue::Bool(false))
            .shorthand('b')
            .description("If true, will run mpv as a background process."),
    );
    flags.add(
        Flag::new("album", FlagValue::Bool(false))
            .shorthand('a')
            .description("If true, will look for media folders and play files from the selected folder in alphabetical order."),
    );
    flags.add(
        Flag::new("shuffle", FlagValue::Bool(false))
            .shorthand('s')
            .description("If true, will shuffle an album. Only works with the --album flag."),
    );
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

    if flags.get_bool("help") {
        let width = match term_size::dimensions() {
            Some((w, _)) => w,
            None => 75,
        };
        usage::generate(&flags, width, &mut io::stdout()).expect("Failed to print usage");
        return;
    }

    let album = flags.get_bool("album");
    let command = match album {
        true => "fd . --type directory | fzf",
        false => "fzf",
    };

    let host = flags.get_string("host");
    let dir = flags.get_string("dir");

    // Use fzf to locate a file
    Command::new("ssh")
        .args(["-t", &host, &format!("cd {dir} && {command} > /tmp/bop")])
        .status()
        .expect("SSH command failed");

    // Retrieve the located file
    let output = Command::new("ssh")
        .args([&host, "cat /tmp/bop"])
        .output()
        .expect("SSH command failed");
    let path = match output.status.success() {
        true => String::from_utf8_lossy(&output.stdout).trim().to_owned(),
        false => panic!("Failed to get result!"),
    };

    let mut paths = vec![];
    if album {
        // List the file paths in an album
        let output = Command::new("ssh")
            .args([&host, &format!("ls \"{dir}/{path}\"")])
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
    if flags.get_bool("shuffle") {
        paths.shuffle(&mut rng());
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

    if flags.get_bool("background") {
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
