pub mod app;
pub mod ssh;

use crate::ssh::server::Server;

#[tokio::main]
async fn main() {
    let mut server = Server::new();
    server.run().await.expect("Failed running server");
}
