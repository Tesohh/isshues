use std::hash::Hash;

use dashmap::DashMap;
use tokio::sync::broadcast;

use crate::app::event::ToClientEvent;

#[derive(Clone, Copy, Hash, PartialEq, Eq)]
pub enum ToClientTopic {
    /// Only subscribed to from the client with that id.
    Client(usize),
}

/// Meant to be used as Arc<EventBus>
pub struct EventBroker<Topic: Hash + PartialEq + Eq, Event: Clone> {
    topics: DashMap<Topic, broadcast::Sender<Event>>,
}

pub type ToClientEventBroker = EventBroker<ToClientTopic, ToClientEvent>;

impl<Topic: Hash + PartialEq + Eq, Event: Clone> EventBroker<Topic, Event> {
    pub fn new() -> Self {
        Self {
            topics: DashMap::new(),
        }
    }
    /// Get a receiver to the broadcast channel associated with a topic.
    /// In case it doesn't exist, it will create a new topic channel.
    pub fn subscribe(&self, topic: Topic) -> broadcast::Receiver<Event> {
        let sender = self
            .topics
            .entry(topic)
            .or_insert_with(|| broadcast::channel(16).0);
        sender.subscribe()
    }

    /// Checks if the broadcast channel associated with a topic is dead,
    /// and in that case removes it.
    ///
    /// Please call this function whenever you unsubscribe (aka drop a receiver) to avoid stale channels.
    pub fn unsubscribe(&self, topic: Topic) {
        self.topics
            .remove_if(&topic, |_, sender| sender.receiver_count() == 0);
    }

    /// Broadcasts a message to a specific topic.
    ///
    /// In case no clients are subscribed (aka. the topic doesn't exist), nothing is sent.
    pub fn publish(&self, topic: Topic, event: Event) {
        if let Some(sender) = self.topics.get(&topic) {
            let _ = sender.send(event); // WARNING: this may fail to update some receivers, and it may be ignored, consider logging
        };
    }
}

impl<Topic: Hash + PartialEq + Eq, Event: Clone> Default for EventBroker<Topic, Event> {
    fn default() -> Self {
        Self::new()
    }
}
