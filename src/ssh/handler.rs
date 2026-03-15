use std::{collections::HashMap, io, sync::Arc};

use ratatui::{Terminal, TerminalOptions, Viewport, layout::Rect, prelude::CrosstermBackend};
use russh::{
    Channel, ChannelId, Pty,
    keys::PublicKey,
    server::{Auth, Msg, Session},
};
use tokio::sync::Mutex;

use crate::{
    app::{
        client::Client,
        event::ToClientEvent,
        event_broker::{ToClientEventBroker, ToClientTopic},
    },
    ssh::terminal::SshTerminalHandle,
};

/// Handles SSH events for and manages a single client.
pub struct ClientHandler {
    clients: Arc<Mutex<HashMap<usize, Client>>>,
    bus: Arc<ToClientEventBroker>,
    id: usize,
    input_parser: termwiz::input::InputParser,
}

impl ClientHandler {
    pub fn new(
        clients: Arc<Mutex<HashMap<usize, Client>>>,
        bus: Arc<ToClientEventBroker>,
        id: usize,
    ) -> Self {
        Self {
            clients,
            bus,
            id,
            input_parser: termwiz::input::InputParser::new(),
        }
    }

    fn close(
        &mut self,
        channel: ChannelId,
        session: &mut Session,
        exit_message: &str,
    ) -> Result<(), io::Error> {
        let _ = session.data(channel, "\x1b[?1049l\x1b[?25h\x1b[0m".into());
        let _ = session.data(channel, exit_message.into());
        let _ = session.data(channel, "\n\r".into());

        let _ = session.exit_status_request(channel, 0);
        let _ = session.eof(channel);
        let _ = session.close(channel);

        Ok(())
    }
}

impl russh::server::Handler for ClientHandler {
    // TODO: proper error
    type Error = anyhow::Error;

    async fn channel_open_session(
        &mut self,
        channel: Channel<Msg>,
        session: &mut Session,
    ) -> Result<bool, Self::Error> {
        let terminal_handle = SshTerminalHandle::start(session.handle(), channel.id()).await;

        let backend = CrosstermBackend::new(terminal_handle);

        // the correct viewport area will be set when the client request a pty
        let options = TerminalOptions {
            viewport: Viewport::Fixed(Rect::default()),
        };

        let terminal = Terminal::with_options(backend, options)?;

        let mut clients = self.clients.lock().await;
        clients.insert(
            self.id,
            Client::new(self.id, terminal, Arc::clone(&self.bus)),
        );
        Ok(true)
    }

    // this is called after the PTY is ready.
    async fn shell_request(
        &mut self,
        channel: ChannelId,
        session: &mut Session,
    ) -> Result<(), Self::Error> {
        let _ = session.channel_success(channel);
        // go into alternate screen, move cursor to top left, hide cursor
        let _ = session.data(channel, "\x1b[?1049h\x1b[H\x1b[?25l".into());
        Ok(())
    }

    async fn auth_publickey(&mut self, _: &str, _: &PublicKey) -> Result<Auth, Self::Error> {
        Ok(Auth::Accept)
    }

    async fn data(
        &mut self,
        _channel: ChannelId,
        data: &[u8],
        _session: &mut Session,
    ) -> Result<(), Self::Error> {
        // TODO: termwiz parse this stuff
        // match data {
        //     // Pressing 'q' closes the connection.
        //     b"q" => {
        //         let _ = self.close(channel, session, "Harris");
        //         let mut clients_lock = self.clients.lock().await;
        //         clients_lock.remove(&self.id);
        //     }
        //     // Pressing 'c' resets the counter for the app.
        //     // Only the client with the id sees the counter reset.
        //     b"c" => {
        //         let mut clients = self.clients.lock().await;
        //         let client = clients.get_mut(&self.id).unwrap();
        //         client.app.counter = 0;
        //     }
        //     _ => {}
        // }
        println!("GOT SOME DATA");

        let events = self.input_parser.parse_as_vec(data, false);
        for event in events {
            println!("GOT {:?}", event);
            self.bus
                .publish(ToClientTopic::Client(self.id), ToClientEvent::Input(event));
        }

        Ok(())
    }

    /// The client's window size has changed.
    async fn window_change_request(
        &mut self,
        _: ChannelId,
        col_width: u32,
        row_height: u32,
        _: u32,
        _: u32,
        _: &mut Session,
    ) -> Result<(), Self::Error> {
        let rect = Rect {
            x: 0,
            y: 0,
            width: col_width as u16,
            height: row_height as u16,
        };

        let mut clients = self.clients.lock().await;
        let client = clients.get_mut(&self.id).unwrap();
        client.terminal.resize(rect)?;

        Ok(())
    }

    /// The client requests a pseudo-terminal with the given
    /// specifications.
    ///
    /// **Note:** Success or failure should be communicated to the client by calling
    /// `session.channel_success(channel)` or `session.channel_failure(channel)` respectively.
    async fn pty_request(
        &mut self,
        channel: ChannelId,
        _: &str,
        col_width: u32,
        row_height: u32,
        _: u32,
        _: u32,
        _: &[(Pty, u32)],
        session: &mut Session,
    ) -> Result<(), Self::Error> {
        let rect = Rect {
            x: 0,
            y: 0,
            width: col_width as u16,
            height: row_height as u16,
        };

        let mut clients = self.clients.lock().await;
        let client = clients.get_mut(&self.id).unwrap();
        client.terminal.resize(rect)?;

        session.channel_success(channel)?;

        Ok(())
    }
}

impl Drop for ClientHandler {
    fn drop(&mut self) {
        let id = self.id;
        let clients = self.clients.clone();
        tokio::spawn(async move {
            let mut clients = clients.lock().await;
            clients.remove(&id);
        });
    }
}
