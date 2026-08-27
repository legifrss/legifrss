use crate::legifrance::{Container, JoContainer, JoContainerResult, JorfContainerResult, LastNJo};
use sqlx::{FromRow, Pool, Postgres, types::Json};

#[derive(FromRow)]
pub struct JorfTextRow {
    pub id: String,
    pub date: chrono::NaiveDateTime,
    pub content: Json<JorfContainerResult>,
    pub nature: Option<String>,
    pub author: Option<String>,
}

pub async fn latest_jorf_text(
    author: Option<&str>,
    nature: Option<&str>,
    keyword: Option<&str>,
    pool: &Pool<Postgres>,
) -> Result<Vec<JorfTextRow>, sqlx::Error> {
    sqlx::query_as::<_, JorfTextRow>(
        "SELECT id, date, content, nature, author FROM jorf_text \
         WHERE ($1::text IS NULL OR author ILIKE '%' || $1 || '%') \
           AND ($2::text IS NULL OR nature ILIKE '%' || $2 || '%') \
           AND ($3::text IS NULL OR content::text ILIKE '%' || $3 || '%') \
         ORDER BY date DESC
         LIMIT 500",
    )
    .bind(author)
    .bind(nature)
    .bind(keyword)
    .fetch_all(pool)
    .await
}

pub async fn distinct_authors(pool: &Pool<Postgres>) -> Result<Vec<String>, sqlx::Error> {
    sqlx::query_scalar::<_, String>(
        "SELECT DISTINCT author FROM jorf_text \
         WHERE author IS NOT NULL AND author <> '' \
         ORDER BY author",
    )
    .fetch_all(pool)
    .await
}

pub async fn distinct_natures(pool: &Pool<Postgres>) -> Result<Vec<String>, sqlx::Error> {
    sqlx::query_scalar::<_, String>(
        "SELECT DISTINCT nature FROM jorf_text \
         WHERE nature IS NOT NULL AND nature <> '' \
         ORDER BY nature",
    )
    .fetch_all(pool)
    .await
}

pub async fn write_all_jorf(last_njo: LastNJo, pool: &Pool<Postgres>) {
    let futures = last_njo
        .containers
        .iter()
        .map(|e| persist_jorf(&e.id, e.clone(), pool));
    futures::future::join_all(futures).await;
}

pub async fn persist_jorf(id: &String, njo: Container, pool: &Pool<Postgres>) {
    sqlx::query_file!(
        "sql/insert_jorf.sql",
        id,
        njo.date.naive_utc(),
        Json(njo) as _
    )
    .execute(pool)
    .await
    .unwrap();
}

pub async fn jorf_cont_exists(id: &str, pool: &Pool<Postgres>) -> bool {
    sqlx::query_scalar::<_, bool>(
        "SELECT EXISTS(SELECT 1 FROM jorf WHERE id = $1 AND jorf_content IS NOT NULL)",
    )
    .bind(id)
    .fetch_one(pool)
    .await
    .unwrap_or(false)
}

pub async fn jorf_text_exists(id: &str, pool: &Pool<Postgres>) -> bool {
    sqlx::query_scalar::<_, bool>(
        "SELECT EXISTS(SELECT 1 FROM jorf_text WHERE id = $1 AND content IS NOT NULL)",
    )
    .bind(id)
    .fetch_one(pool)
    .await
    .unwrap_or(false)
}

pub async fn write_all_jorf_cont(result: JoContainerResult, pool: &Pool<Postgres>) {
    let futures = result
        .items
        .iter()
        .map(|e| persist_jorf_cont(&e.container.id, e.container.clone(), pool));
    futures::future::join_all(futures).await;
}

pub async fn persist_jorf_cont(id: &String, cont: JoContainer, pool: &Pool<Postgres>) {
    sqlx::query_file!("sql/insert_jorf_content.sql", id, Json(cont) as _)
        .execute(pool)
        .await
        .unwrap();
}

pub async fn persist_jorf_text(
    id: &String,
    date: chrono::DateTime<chrono::Utc>,
    text: JorfContainerResult,
    nature: Option<String>,
    author: Option<String>,
    pool: &Pool<Postgres>,
) {
    sqlx::query_file!(
        "sql/insert_jorf_text.sql",
        id,
        date.naive_utc(),
        Json(text) as _,
        nature,
        author
    )
    .execute(pool)
    .await
    .unwrap();
}
