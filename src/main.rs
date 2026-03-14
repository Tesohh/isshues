pub mod app_client;
pub mod app_handler;
pub mod app_server;
pub mod ssh_terminal_handle;

use crate::app_server::AppServer;

#[tokio::main]
async fn main() {
    let mut server = AppServer::new();
    server.run().await.expect("Failed running server");
}
