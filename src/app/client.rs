use crate::ssh::terminal::SshTerminal;

pub struct Client {
    pub terminal: SshTerminal,
    pub app: App,
}

impl Client {
    pub fn new(terminal: SshTerminal) -> Self {
        Self {
            terminal,
            app: App::default(),
        }
    }
}

pub struct App {
    pub counter: usize,
}

impl App {
    pub fn new() -> App {
        Self { counter: 0 }
    }
}

impl Default for App {
    fn default() -> Self {
        Self::new()
    }
}
