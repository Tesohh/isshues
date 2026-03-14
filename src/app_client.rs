use crate::ssh_terminal_handle::SshTerminal;

pub struct App {
    pub counter: usize,
}

// App that runs per client.
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

pub struct AppClient {
    pub terminal: SshTerminal,
    pub app: App,
}
