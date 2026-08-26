mod batch;
mod env;
mod handlers;
mod legifrance;
mod model;
mod oauth;
mod persist;
mod rss;
use crate::model::DynamicState;
use actix_web::{App, HttpServer, web};
use sqlx::{migrate::Migrator, postgres::PgPoolOptions};
use std::{collections::HashMap, path::Path, sync::Mutex, time::Duration};

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let password_file = std::env::var("CREDENTIALS_FILE")
        .expect("You should define a file with at least the secrets inside. (Not just in a .env)");

    let config: model::Config = env::load_env(password_file);
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .acquire_timeout(Duration::from_secs(5))
        .connect(&config.database_url)
        .await
        .unwrap();

    Migrator::new(Path::new("./migrations"))
        .await
        .unwrap()
        .run(&pool)
        .await
        .unwrap();

    let client = reqwest::Client::builder()
        .user_agent("legifrss")
        .read_timeout(Duration::from_secs(20))
        .build()
        .unwrap();

    let state = web::Data::new(Mutex::new(DynamicState { oauth: None }));
    let cache: web::Data<handlers::FeedCache> = web::Data::new(Mutex::new(HashMap::new()));

    HttpServer::new(move || {
        App::new()
            .service(web::resource("/").route(web::get().to(handlers::index)))
            .service(handlers::batch)
            .service(web::resource("/latest").route(web::get().to(handlers::stream)))
            .service(web::resource("/latest.xml").route(web::get().to(handlers::stream)))
            .service(web::resource("/authors").route(web::get().to(handlers::authors)))
            .service(web::resource("/natures").route(web::get().to(handlers::natures)))
            .app_data(web::Data::new(pool.clone()))
            .app_data(web::Data::new(client.clone()))
            .app_data(state.clone())
            .app_data(cache.clone())
            .app_data(web::Data::new(config.clone()))
    })
    .bind((
        "0.0.0.0",
        std::env::var("PORT")
            .ok()
            .map(|p| p.parse::<u16>().unwrap())
            .unwrap_or(8080),
    ))?
    .run()
    .await
}
