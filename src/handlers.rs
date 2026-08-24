use crate::model::DynamicState;
use actix_web::{HttpResponse, post, web};
use sqlx::{Pool, Postgres};
use std::sync::Mutex;

const INDEX_HTML: &str = include_str!("../assets/index.html");

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

    HttpResponse::Ok()
        .content_type("application/json")
        .body("ok")
}

#[derive(serde::Deserialize)]
pub struct LatestQuery {
    pub nature: Option<String>,
    pub author: Option<String>,
}

pub async fn stream(
    query: web::Query<LatestQuery>,
    pool: web::Data<Pool<Postgres>>,
) -> HttpResponse {
    let query = query.into_inner();
    HttpResponse::Ok()
        .content_type("application/atom+xml")
        .body(crate::rss::latest(query.author, query.nature, &pool).await)
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
