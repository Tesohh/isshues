use std::sync::Arc;

use tokio::sync::broadcast;
use tokio_stream::StreamMap;

use crate::{
    app::{
        event::ToClientEvent,
        event_broker::{ToClientEventBroker, ToClientTopic},
    },
    ssh::terminal::SshTerminal,
};

pub struct Client {
    pub id: usize,
    pub terminal: SshTerminal,
    pub app: App,

    pub bus: Arc<ToClientEventBroker>,
    pub subscriptions: StreamMap<ToClientTopic, broadcast::Receiver<ToClientEvent>>,
}

impl Client {
    pub fn new(id: usize, terminal: SshTerminal, bus: Arc<ToClientEventBroker>) -> Self {
        let mut map = StreamMap::new();
        let topic = ToClientTopic::Client(id);
        map.insert(topic, bus.subscribe(topic));

        Self {
            id,
            terminal,
            app: App::default(),
            bus,
            subscriptions: map,
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
