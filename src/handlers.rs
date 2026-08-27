use crate::model::DynamicState;
use actix_web::{HttpResponse, post, web};
use sqlx::{Pool, Postgres};
use std::{
    collections::HashMap,
    sync::Mutex,
    time::{Duration, Instant},
};

const INDEX_HTML: &str = include_str!("../assets/index.html");

/// Feeds only change when a batch run imports new texts, so the TTL is just a
/// safety net: `/batch` clears the cache itself.
const CACHE_TTL: Duration = Duration::from_mins(90);

pub type FeedCache = Mutex<HashMap<LatestQuery, (Instant, String)>>;

pub async fn index() -> HttpResponse {
    HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(INDEX_HTML)
}

#[post("/batch")]
async fn batch(
    http_client: web::Data<reqwest::Client>,
    oauth: web::Data<Mutex<DynamicState>>,
    config: web::Data<crate::model::Config>,
    pool: web::Data<Pool<Postgres>>,
    cache: web::Data<FeedCache>,
) -> HttpResponse {
    crate::batch::batch_jorf(
        http_client.get_ref().clone(),
        config.client_id.clone(),
        config.client_secret.clone(),
        config.api_url.clone(),
        config.oauth_url.clone(),
        oauth.get_ref(),
        &pool,
    )
    .await;

    cache.lock().unwrap().clear();

    HttpResponse::Ok()
        .content_type("application/json")
        .body("ok")
}

#[derive(serde::Deserialize, Hash, PartialEq, Eq, Clone)]
pub struct LatestQuery {
    pub nature: Option<String>,
    pub author: Option<String>,
    pub q: Option<String>,
}

pub async fn stream(
    query: web::Query<LatestQuery>,
    pool: web::Data<Pool<Postgres>>,
    cache: web::Data<FeedCache>,
) -> HttpResponse {
    let query = query.into_inner();

    let cached = cache
        .lock()
        .unwrap()
        .get(&query)
        .filter(|(stored_at, _)| stored_at.elapsed() < CACHE_TTL)
        .map(|(_, feed)| feed.clone());

    let feed = match cached {
        Some(feed) => feed,
        None => {
            let feed = crate::rss::latest(
                query.author.clone(),
                query.nature.clone(),
                query.q.clone(),
                &pool,
            )
            .await;
            cache
                .lock()
                .unwrap()
                .insert(query, (Instant::now(), feed.clone()));
            feed
        }
    };

    HttpResponse::Ok()
        .content_type("application/atom+xml")
        .body(feed)
}

pub async fn authors(pool: web::Data<Pool<Postgres>>) -> HttpResponse {
    match crate::persist::distinct_authors(&pool).await {
        Ok(authors) => HttpResponse::Ok().json(authors),
        Err(err) => {
            eprintln!("authors query failed: {err}");
            HttpResponse::InternalServerError().finish()
        }
    }
}

pub async fn natures(pool: web::Data<Pool<Postgres>>) -> HttpResponse {
    match crate::persist::distinct_natures(&pool).await {
        Ok(natures) => HttpResponse::Ok().json(natures),
        Err(err) => {
            eprintln!("natures query failed: {err}");
            HttpResponse::InternalServerError().finish()
        }
    }
}
