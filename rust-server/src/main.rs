use std::collections::HashMap;
use std::env;
use std::net::SocketAddr;
use std::sync::{Arc, Mutex};

use axum::extract::ws::Message;
use axum::extract::WebSocketUpgrade;
use axum::response::IntoResponse;
use axum::routing::get;
use axum::{Extension, Router};
use dotenv::dotenv;
use futures::{SinkExt, StreamExt};

use rand::Rng;
use serde_json::json;
use tokio::sync::broadcast;
use tower_http::catch_panic::CatchPanicLayer;
use tower_http::cors::CorsLayer;
use tower_http::trace::TraceLayer;
use tracing_subscriber::layer::SubscriberExt;
use tracing_subscriber::util::SubscriberInitExt;
use uuid::Uuid;

mod config;

#[tokio::main]
async fn main() {
    env::set_var("RUST_BACKTRACE", "1");
    dotenv().ok();
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env().unwrap_or_else(|_| {
                "axum=debug,rust_server=debug,tower_http=debug,mongodb=debug".into()
            }),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    let port = std::env::var("SERVER_PORT")
        .unwrap_or(String::new())
        .parse()
        .unwrap_or(8888);

    let app_state = Arc::new(config::AppState {
        session: Mutex::new(HashMap::new()),
        tx: broadcast::channel(config::MAX_SIZE).0,
    });

    let app = Router::new()
        .route("/world", get(world_handler))
        .layer(TraceLayer::new_for_http())
        .layer(CorsLayer::permissive())
        .layer(CatchPanicLayer::new())
        .layer(Extension(app_state));

    axum::Server::bind(&SocketAddr::from(([0, 0, 0, 0], port)))
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn world_handler(
    ws: WebSocketUpgrade,
    Extension(state): Extension<Arc<config::AppState>>,
) -> impl IntoResponse {
    ws.on_upgrade(move |socket| async move {
        let (mut sender, mut receiver) = socket.split();

        let (msg_init, msg_join) = join_world(&state).await;
        let msg_init = serde_json::to_string(&msg_init).unwrap();
        let msg_join = serde_json::to_string(&msg_join).unwrap();

        if sender.send(Message::Text(msg_init)).await.is_err() {
            return;
        }

        let _ = state.tx.send(msg_join);

        let recv_state = state.clone();
        let mut recv_task = tokio::spawn(async move {
            while let Some(Ok(message)) = receiver.next().await {
                let message =
                    serde_json::from_str::<config::Message>(message.to_text().unwrap()).unwrap();

                match message.message_type {
                    config::MessageType::MOVE => {
                        let payload = message.payload.as_object().unwrap();
                        let mut session = recv_state.session.lock().unwrap();
                        let id = payload.get("id").unwrap().as_str().unwrap();

                        if let Some(player) = session.get_mut(id) {
                            player.x = payload.get("x").unwrap().as_u64().unwrap_or(player.x);
                            player.y = payload.get("y").unwrap().as_u64().unwrap_or(player.y);
                            let _ = recv_state.tx.send(serde_json::to_string(&message).unwrap());
                        }
                    }
                    config::MessageType::ATTACK => {
                        let payload = message.payload.as_object().unwrap();
                        let mut session = recv_state.session.lock().unwrap();
                        let id = payload.get("id").unwrap().as_str().unwrap();

                        if let Some(player) = session.get_mut(id) {
                            player.hp += 1;
                            let _ = recv_state.tx.send(serde_json::to_string(&message).unwrap());
                        }
                    }
                    _ => {}
                }
            }
        });

        let mut send_task = tokio::spawn(async move {
            while let Ok(msg) = state.tx.subscribe().recv().await {
                if sender.send(Message::Text(msg)).await.is_err() {
                    break;
                }
            }
        });

        tokio::select! {
            _ = (&mut send_task) => recv_task.abort(),
            _ = (&mut recv_task) => send_task.abort(),
        }
    })
}

async fn join_world(state: &Arc<config::AppState>) -> (config::Message, config::Message) {
    let id = Uuid::new_v4().to_string().replace("-", "");
    let mut rng = rand::thread_rng();
    let player = config::Player {
        id: id.to_string(),
        hp: 5,
        x: rng.gen_range(config::X_MIN..config::X_MAX),
        y: rng.gen_range(config::Y_MIN..config::Y_MAX),
    };

    let mut session = state.session.lock().unwrap();
    session.insert(id, player.clone());

    (
        config::Message {
            message_type: config::MessageType::INIT,
            payload: json!({
                "player":  player,
                "players": session.clone()
            }),
        },
        config::Message {
            message_type: config::MessageType::JOIN,
            payload: json!({
               "id":player.id,
               "hp":player.hp,
               "x":player.x,
               "y":player.y,
            }),
        },
    )
}
