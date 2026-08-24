use std::sync::Mutex;

use futures::stream::StreamExt;
use sqlx::{Pool, Postgres};

const JORF_CONCURRENCY: usize = 1;

use crate::{
    legifrance::{
        HierarchyStep, Summary, get_consult_jorf, get_consult_jorf_cont, get_consult_last_n_jo,
    },
    model::{DynamicState, Oauth},
};

pub async fn batch_jorf(
    http_client: reqwest::Client,
    client_id: String,
    client_secret: String,
    api_url: String,
    oauth_url: String,
    existing_oauth: &Mutex<DynamicState>,
    pool: &Pool<Postgres>,
) {
    let maybe_token: Option<Oauth> = crate::oauth::get_token(
        http_client.clone(),
        client_id,
        client_secret,
        oauth_url,
        existing_oauth,
    )
    .await;

    let token = match maybe_token.and_then(|oauth| oauth.token) {
        Some(token) => token,
        None => {
            eprintln!("batch_jorf: no OAuth token available, skipping run");
            return;
        }
    };
    let last_njo =
        match get_consult_last_n_jo(http_client.clone(), api_url.clone(), token.clone()).await {
            Some(last_njo) => last_njo,
            None => {
                eprintln!("batch_jorf: consult/lastNJo returned no result, skipping run");
                return;
            }
        };

    crate::persist::write_all_jorf(last_njo.clone(), pool).await;

    let container_futures = last_njo.containers.iter().map(async |container| {
        if crate::persist::jorf_cont_exists(&container.id, pool).await {
            return;
        }
        let cont_result = match get_consult_jorf_cont(
            http_client.clone(),
            api_url.clone(),
            token.clone(),
            container.id.clone(),
        )
        .await
        {
            Some(cont_result) => cont_result,
            None => {
                eprintln!(
                    "skipping jorf container {}: consult/jorfCont returned no usable result",
                    container.id
                );
                return;
            }
        };
        crate::persist::write_all_jorf_cont(cont_result.clone(), pool).await;

        let text_futures = cont_result
            .items
            .iter()
            .flat_map(|item| collect_summaries(&item.container.structure.contents))
            .map(async |summary| {
                if crate::persist::jorf_text_exists(&summary.id, pool).await {
                    return;
                }
                let text = get_consult_jorf(
                    http_client.clone(),
                    api_url.clone(),
                    token.clone(),
                    summary.id.clone(),
                )
                .await;
                match text {
                    Some(text) => {
                        let author = summary.emitter.clone().or_else(|| summary.minister.clone());
                        crate::persist::persist_jorf_text(
                            &summary.id,
                            container.date,
                            text,
                            summary.nature.clone(),
                            author,
                            pool,
                        )
                        .await
                    }
                    None => eprintln!(
                        "skipping jorf text {}: consult/jorf returned no usable result",
                        summary.id
                    ),
                }
            });
        futures::stream::iter(text_futures)
            .buffer_unordered(JORF_CONCURRENCY)
            .collect::<Vec<()>>()
            .await;
    });
    futures::stream::iter(container_futures)
        .buffer_unordered(1)
        .collect::<Vec<()>>()
        .await;
}

fn collect_summaries(steps: &[HierarchyStep]) -> Vec<&Summary> {
    steps
        .iter()
        .flat_map(|step| {
            let mut summaries: Vec<&Summary> = step.summaries.iter().collect();
            summaries.extend(collect_summaries(&step.step));
            summaries
        })
        .collect()
}
