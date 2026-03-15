use std::{collections::HashMap, io, sync::Arc};

use ratatui::widgets::{Block, Borders, Clear, Paragraph};
use russh::{
    keys::{
        Algorithm, PrivateKey,
        ssh_key::{LineEnding, rand_core::OsRng},
    },
    server::Server as _,
};
use tokio::sync::Mutex;

use crate::{
    app::{client::Client, event_broker::ToClientEventBroker},
    ssh::handler::ClientHandler,
};

#[derive(Clone)]
pub struct Server {
    /// maps client ids to AppClients
    clients: Arc<Mutex<HashMap<usize, Client>>>,

    bus: Arc<ToClientEventBroker>,

    /// id to assign to the next client
    next_id: usize,
}

// TODO: fix unwraps
// TODO: move this in a file
fn load_or_generate_key() -> Result<PrivateKey, io::Error> {
    let path = dirs::config_local_dir()
        .unwrap()
        // .ok_or_eyre("Failed to get config local dir")?
        .join("isshues")
        .join("host_key");
    let key = if path.exists() {
        // info!("Loading host key from {}", path.display());
        PrivateKey::read_openssh_file(&path).unwrap() //.wrap_err("Failed to read host key from file")?
    } else {
        // info!(
        //     "Host key not found at {}. Generating new host key",
        //     path.display()
        // );
        let key = PrivateKey::random(&mut OsRng, Algorithm::Ed25519).unwrap();
        // .wrap_err("Failed to generate host key")?;
        std::fs::create_dir_all(path.parent().unwrap()).unwrap();
        // .wrap_err("Failed to create directory for host key")?;
        key.write_openssh_file(&path, LineEnding::LF).unwrap();
        // .wrap_err("Failed to write host key to file")?;
        key
    };
    Ok(key)
}

impl Server {
    pub fn new() -> Self {
        Self {
            clients: Arc::new(Mutex::new(HashMap::new())),
            next_id: 0,
            bus: Arc::new(ToClientEventBroker::new()),
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
            keys: vec![load_or_generate_key().unwrap()], // TODO: remove unwrap
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
        ClientHandler::new(
            Arc::clone(&self.clients),
            Arc::clone(&self.bus),
            self.next_id,
        )
    }
}

impl Default for Server {
    fn default() -> Self {
        Self::new()
    }
}
