use chrono::{self, Utc};
use serde::{Deserialize, Serialize};

pub async fn get_consult_last_n_jo(
    http_client: reqwest::Client,
    api_url: String,
    token: String,
) -> Option<LastNJo> {
    let result = http_client
        .post(format!("{api_url}/consult/lastNJo"))
        .header("Content-Type", "application/json")
        .header("Authorization", format!("Bearer {token}"))
        .body("{\"nbElement\": 5}")
        .send()
        .await
        .inspect_err(|err| println!("consult/lastNJo request failed: {err}"))
        .ok()?;
    let status = result.status();
    if !status.is_success() {
        println!("consult/lastNJo returned HTTP {status}");
        return None;
    }
    let contents = result
        .text()
        .await
        .inspect_err(|err| println!("consult/lastNJo read failed: {err}"))
        .ok()?;
    serde_json::from_str(&contents)
        .inspect_err(|err| println!("{err}, {contents}"))
        .ok()
}

pub async fn get_consult_jorf_cont(
    http_client: reqwest::Client,
    api_url: String,
    token: String,
    jorf_cont_id: String,
) -> Option<JoContainerResult> {
    println!("Call consult/jorf/cont {jorf_cont_id}");
    let result = http_client
        .post(format!("{api_url}/consult/jorfCont"))
        .header("Authorization", format!("Bearer {token}"))
        .json(&serde_json::json!({
            "id": jorf_cont_id,
            "pageNumber": 1,
            "pageSize": 10,
        }))
        .send()
        .await
        .inspect_err(|err| println!("consult/jorfCont request failed: {err}"))
        .ok()?;
    let status = result.status();
    if !status.is_success() {
        println!("consult/jorfCont returned HTTP {status}");
        return None;
    }
    let contents = result
        .text()
        .await
        .inspect_err(|err| println!("consult/jorfCont read failed: {err}"))
        .ok()?;
    serde_json::from_str(&contents)
        .inspect_err(|err| println!("{err}, {contents}"))
        .ok()
}

pub async fn get_consult_jorf(
    http_client: reqwest::Client,
    api_url: String,
    token: String,
    text_cid: String,
) -> Option<JorfContainerResult> {
    println!("Call consult/jorf {text_cid}");
    let result = http_client
        .post(format!("{api_url}/consult/jorf"))
        .header("Authorization", format!("Bearer {token}"))
        .json(&serde_json::json!({ "textCid": text_cid }))
        .send()
        .await
        .inspect_err(|err| println!("consult/jorf request failed: {err}"))
        .ok()?;

    let status = result.status();
    if !status.is_success() {
        println!("consult/jorf returned HTTP {status}");
        return None;
    }
    let contents = result
        .text()
        .await
        .inspect_err(|err| println!("consult/jorf read failed: {err}"))
        .ok()?;
    serde_json::from_str(&contents)
        .inspect_err(|err| println!("{err}, {contents}"))
        .ok()
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct LastNJo {
    pub containers: Vec<Container>,
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct Container {
    pub id: String,
    #[serde(rename = "titre")]
    pub title: Option<String>,
    #[serde(rename = "idEli")]
    pub id_eli: String,
    #[serde(rename = "datePubli", with = "chrono::serde::ts_milliseconds")]
    pub date: chrono::DateTime<Utc>,
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct Summary {
    pub id: String,
    #[serde(rename = "titre")]
    pub title: Option<String>,
    pub nature: Option<String>,
    #[serde(rename = "ministere")]
    pub minister: Option<String>,
    #[serde(rename = "emetteur")]
    pub emitter: Option<String>,
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct HierarchyStep {
    #[serde(rename = "titre")]
    pub title: String,
    #[serde(rename = "niv")]
    pub level: i32,
    #[serde(rename = "tms", default)]
    pub step: Vec<HierarchyStep>,
    #[serde(rename = "liensTxt", default)]
    pub summaries: Vec<Summary>,
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct Structure {
    #[serde(rename = "tms", default)]
    pub contents: Vec<HierarchyStep>,
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct JoContainer {
    pub id: String,
    pub structure: Structure,
    #[serde(rename = "datePubli", with = "chrono::serde::ts_milliseconds")]
    pub timestamp: chrono::DateTime<Utc>,
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct Item {
    #[serde(rename = "joCont")]
    pub container: JoContainer,
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct JoContainerResult {
    pub items: Vec<Item>,
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct JorfContainerResult {
    #[serde(rename = "cid")]
    pub id: String,
    pub title: String,
    #[serde(default)]
    pub sections: Vec<JorfContainerSection>,
    #[serde(default)]
    pub articles: Vec<JorfArticle>,
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct JorfContainerSection {
    pub title: String,
    #[serde(default)]
    pub articles: Vec<JorfArticle>,
    #[serde(default)]
    pub sections: Vec<JorfContainerSection>,
}

#[derive(Debug, Hash, Deserialize, Serialize, PartialEq, Eq, Clone)]
pub struct JorfArticle {
    pub content: String,
    #[serde(rename = "num")]
    pub order: Option<String>,
}
