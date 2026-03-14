use ratatui::Terminal;
use ratatui::prelude::CrosstermBackend;
use russh::ChannelId;
use russh::server::*;
use tokio::sync::mpsc::{UnboundedSender, unbounded_channel};

pub type SshTerminal = Terminal<CrosstermBackend<SshTerminalHandle>>;

pub struct SshTerminalHandle {
    sender: UnboundedSender<Vec<u8>>,
    /// collects the data which is finally sent to sender
    sink: Vec<u8>,
}

impl SshTerminalHandle {
    pub async fn start(handle: Handle, channel_id: ChannelId) -> Self {
        let (sender, mut receiver) = unbounded_channel::<Vec<u8>>();
        tokio::spawn(async move {
            while let Some(data) = receiver.recv().await {
                let result = handle.data(channel_id, data.into()).await;
                if result.is_err() {
                    eprintln!("Failed to send data: {result:?}");
                }
            }
        });
        Self {
            sender,
            sink: Vec::new(),
        }
    }
}

// The crossterm backend writes to the terminal handle.
impl std::io::Write for SshTerminalHandle {
    fn write(&mut self, buf: &[u8]) -> std::io::Result<usize> {
        self.sink.extend_from_slice(buf);
        Ok(buf.len())
    }

    fn flush(&mut self) -> std::io::Result<()> {
        self.sender
            .send(self.sink.clone())
            .map_err(|err| std::io::Error::new(std::io::ErrorKind::BrokenPipe, err))?;

        self.sink.clear();
        Ok(())
    }
}
