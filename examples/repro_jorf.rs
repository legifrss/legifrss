use std::time::Duration;

#[tokio::main]
async fn main() {
    let token = std::env::var("TOK").unwrap();
    let api_url = "https://api.piste.gouv.fr/dila/legifrance/lf-engine-app".to_string();

    let http_client = reqwest::Client::builder()
        .pool_max_idle_per_host(0)
        .read_timeout(Duration::from_secs(20))
        .build()
        .unwrap();

    // A handful of real JORFTEXT ids, repeated to build volume, fired concurrently
    let base = ["JORFTEXT000054427175", "JORFTEXT000054443420"];
    let ids: Vec<String> = (0..40).map(|i| base[i % base.len()].to_string()).collect();

    let futs = ids.iter().enumerate().map(|(i, id)| {
        let http_client = http_client.clone();
        let api_url = api_url.clone();
        let token = token.clone();
        let id = id.clone();
        async move {
            let result = http_client
                .post(format!("{api_url}/consult/jorf"))
                .header("Authorization", format!("Bearer {token}"))
                .json(&serde_json::json!({ "textCid": id }))
                .send()
                .await
                .unwrap();
            let status = result.status();
            let contents = result.text().await.unwrap();
            let ok_cid = contents.contains("\"cid\"");
            let echoed = contents.trim_start().starts_with("{\"textCid\"");
            println!(
                "[{i:02}] status={status} len={} cid={ok_cid} echoed={echoed} {}",
                contents.len(),
                if ok_cid {
                    "".to_string()
                } else {
                    contents.chars().take(120).collect::<String>()
                }
            );
        }
    });
    futures::future::join_all(futs).await;
}
