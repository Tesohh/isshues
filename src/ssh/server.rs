use std::{collections::HashMap, sync::Arc};

use ratatui::widgets::{Block, Borders, Clear, Paragraph};
use russh::{
    keys::ssh_key::{self, rand_core::OsRng},
    server::Server as _,
};
use tokio::sync::Mutex;

use crate::{app::client::Client, ssh::handler::ClientHandler};

#[derive(Clone)]
pub struct Server {
    /// maps client ids to AppClients
    clients: Arc<Mutex<HashMap<usize, Client>>>,

    /// id to assign to the next client
    next_id: usize,
}

impl Server {
    pub fn new() -> Self {
        Self {
            clients: Arc::new(Mutex::new(HashMap::new())),
            next_id: 0,
        }
    }

    pub async fn run(&mut self) -> Result<(), anyhow::Error> {
        let clients = self.clients.clone();
        tokio::spawn(async move {
            loop {
                tokio::time::sleep(tokio::time::Duration::from_secs(1)).await;

                for (_, client) in clients.lock().await.iter_mut() {
                    client
                        .terminal
                        .draw(|f| {
                            let area = f.area();
                            f.render_widget(Clear, area);
                            let paragraph = Paragraph::new("from the server")
                                .alignment(ratatui::layout::Alignment::Center);
                            let block = Block::default().title("Hello world").borders(Borders::ALL);
                            f.render_widget(paragraph.block(block), area);
                        })
                        .unwrap();
                }
            }
        });

        let config = russh::server::Config {
            inactivity_timeout: Some(std::time::Duration::from_secs(3600)),
            auth_rejection_time: std::time::Duration::from_secs(3),
            auth_rejection_time_initial: Some(std::time::Duration::from_secs(0)),
            keys: vec![
                russh::keys::PrivateKey::random(&mut OsRng, ssh_key::Algorithm::Ed25519).unwrap(),
            ],
            nodelay: true,
            ..Default::default()
        };

        self.run_on_address(Arc::new(config), ("0.0.0.0", 2222))
            .await?;
        Ok(())
    }
}

/// Trait used to create new handlers when clients connect.
impl russh::server::Server for Server {
    type Handler = ClientHandler;
    fn new_client(&mut self, _: Option<std::net::SocketAddr>) -> ClientHandler {
        self.next_id += 1;
        ClientHandler::new(Arc::clone(&self.clients), self.next_id)
    }
}

impl Default for Server {
    fn default() -> Self {
        Self::new()
    }
}
