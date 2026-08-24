use crate::model::{DynamicState, Oauth};
use serde::{Deserialize, Serialize};
use sqlx::types::chrono::Local;
use std::{sync::Mutex, time::Duration};

pub async fn get_token(
    http_client: reqwest::Client,
    client_id: String,
    client_secret: String,
    oauth_url: String,
    state: &Mutex<DynamicState>,
) -> Option<Oauth> {
    let token: Option<Oauth> = state.lock().unwrap().oauth.clone();
    let new_token =
        refresh_token_if_needed(http_client, client_id, client_secret, oauth_url, token).await;
    state.lock().unwrap().oauth = new_token.clone();
    println!("{new_token:?}");
    new_token
}

pub async fn refresh_token_if_needed(
    http_client: reqwest::Client,
    client_id: String,
    client_secret: String,
    oauth_url: String,
    existing_oauth: Option<Oauth>,
) -> Option<Oauth> {
    match existing_oauth.filter(|oauth| oauth.expires > Local::now()) {
        Some(oauth) => Some(oauth),
        None => {
            let params = [
                ("grant_type", "client_credentials".to_string()),
                ("client_id", client_id),
                ("client_secret", client_secret),
                ("scope", "openid".to_string()),
            ];
            let res = http_client
                .post(format!("{oauth_url}/oauth/token"))
                .form(&params)
                .timeout(Duration::from_secs(60))
                .send()
                .await
                .unwrap();
            let result = res
                .error_for_status()
                .expect("")
                .json::<OauthResult>()
                .await
                .unwrap();
            Some(Oauth {
                expires: Local::now() + Duration::from_secs(result.expires_in),
                token: Some(result.access_token),
            })
        }
    }
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
struct OauthResult {
    access_token: String,
    token_type: String,
    expires_in: u64,
    scope: String,
}
