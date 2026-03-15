/// Events sent from the server, or the handler, to the clients
#[derive(Clone)]
pub enum ToClientEvent {
    Input(termwiz::input::InputEvent),
}

/// Events sent from the clients to the server
#[derive(Clone)]
pub enum ToServerEvent {}
