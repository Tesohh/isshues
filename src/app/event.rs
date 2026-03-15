/// Events sent from the server, or the handler, to the clients
pub enum ToClientEvent {
    Input(termwiz::input::InputEvent),
}

/// Events sent from the clients to the server
pub enum ToServerEvent {}
