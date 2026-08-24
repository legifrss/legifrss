use serde::{Deserialize, Serialize};
use sqlx::types::chrono::{self, Local};
#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]

pub struct Config {
    pub database_url: String,
    pub client_secret: String,
    pub client_id: String,
    pub oauth_url: String,
    pub api_url: String,
}
#[derive(Debug, Hash, PartialEq, Eq, Clone)]

pub struct DynamicState {
    pub oauth: Option<Oauth>,
}
#[derive(Debug, Hash, PartialEq, Eq, Clone)]
pub struct Oauth {
    pub expires: chrono::DateTime<Local>,
    pub token: Option<String>,
}
