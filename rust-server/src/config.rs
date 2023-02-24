use serde::{Deserialize, Serialize};
use serde_json::Value;
use serde_repr::{Deserialize_repr, Serialize_repr};
use std::{collections::HashMap, sync::Mutex};
use tokio::sync::broadcast;

pub static MAX_SIZE: usize = 500;
pub static X_MIN: u64 = 0;
pub static X_MAX: u64 = 20;
pub static Y_MIN: u64 = 0;
pub static Y_MAX: u64 = 20;

#[derive(Serialize_repr, Deserialize_repr, Debug)]
#[repr(u8)]
pub enum MessageType {
    JOIN = 0,
    LEAVE = 1,
    INIT = 2,
    MOVE = 3,
    ATTACK = 4,
    DIE = 5,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct Message {
    pub message_type: MessageType,
    pub payload: Value,
}

#[derive(Serialize, Clone)]
pub struct Player {
    pub id: String,
    pub hp: u64,
    pub x: u64,
    pub y: u64,
}

pub struct AppState {
    pub session: Mutex<HashMap<String, Player>>,
    pub tx: broadcast::Sender<String>,
}
