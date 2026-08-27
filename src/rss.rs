use crate::legifrance::{JorfArticle, JorfContainerResult, JorfContainerSection};
use crate::persist::{JorfTextRow, latest_jorf_text};
use serde::Serialize;
use sqlx::{Pool, Postgres};

const ATOM_NS: &str = "http://www.w3.org/2005/Atom";
const JORF_BASE_URL: &str = "https://www.legifrance.gouv.fr/jorf/id/";

pub async fn latest(
    author: Option<String>,
    nature: Option<String>,
    keyword: Option<String>,
    pool: &Pool<Postgres>,
) -> String {
    let author = normalize_filter(author);
    let nature = normalize_filter(nature);
    let keyword = normalize_filter(keyword);

    let rows = match latest_jorf_text(
        author.as_deref(),
        nature.as_deref(),
        keyword.as_deref(),
        pool,
    )
    .await
    {
        Ok(rows) => rows,
        Err(err) => {
            eprintln!("latest query failed: {err}");
            return String::new();
        }
    };

    let feed = build_feed(rows, author.as_deref(), nature.as_deref(), keyword.as_deref());
    to_xml(&feed).unwrap_or_default()
}

fn normalize_filter(value: Option<String>) -> Option<String> {
    value.filter(|s| !s.trim().is_empty())
}

fn build_feed(
    rows: Vec<JorfTextRow>,
    author: Option<&str>,
    nature: Option<&str>,
    keyword: Option<&str>,
) -> Feed {
    let entries: Vec<Entry> = rows.into_iter().map(transform_row).collect();
    let updated = entries
        .iter()
        .map(|entry| entry.updated)
        .max()
        .unwrap_or_else(chrono::Utc::now);

    let suffix = [author, nature, keyword]
        .into_iter()
        .flatten()
        .collect::<Vec<_>>()
        .join(" ");

    Feed {
        xmlns: ATOM_NS.to_string(),
        title: format!("Legifrance RSS {suffix}").trim_end().to_string(),
        id: "https://legifrss.org/latest".to_string(),
        updated,
        logo: Some("https://www.legifrance.gouv.fr/favicon.ico".to_string()),
        subtitle: Some(
            "This is a non-official RSS feed for Legifrance's Official Law updates. \
             If you want to follow that topic, you can find more info at https://legifrss.org/."
                .to_string(),
        ),
        author: Author {
            name: "Luca Di Carlo".to_string(),
            email: Some("luca@di-carlo.fr".to_string()),
            uri: None,
        },
        links: vec![Link {
            href: "https://legifrss.org/latest".to_string(),
            rel: Some("self".to_string()),
            r#type: None,
        }],
        entries,
    }
}

fn transform_row(row: JorfTextRow) -> Entry {
    let date = row.date.and_utc();
    let url = format!("{JORF_BASE_URL}{}", row.id);
    let content = extract_content(&row.content.0);

    Entry {
        title: row.content.0.title.clone(),
        id: url.clone(),
        updated: date,
        published: Some(date),
        author: row.author.map(|name| Author {
            name,
            email: None,
            uri: None,
        }),
        links: vec![Link {
            href: url,
            rel: None,
            r#type: None,
        }],
        summary: None,
        content: (!content.is_empty()).then(|| Content {
            r#type: "html".to_string(),
            value: content,
        }),
        categories: row
            .nature
            .map(|term| vec![Category { term, label: None }])
            .unwrap_or_default(),
    }
}

fn extract_content(result: &JorfContainerResult) -> String {
    let mut parts = Vec::new();
    collect_parts(&result.articles, &result.sections, &mut parts);
    parts.join("\n")
}

fn collect_parts(
    articles: &[JorfArticle],
    sections: &[JorfContainerSection],
    parts: &mut Vec<String>,
) {
    let mut articles: Vec<&JorfArticle> = articles.iter().collect();
    articles.sort_by_key(|article| article_order(article));
    for article in articles {
        parts.push(article.content.clone());
    }

    let mut sections: Vec<&JorfContainerSection> = sections.iter().collect();
    sections.sort_by_key(|section| section_order(section));
    for section in sections {
        collect_parts(&section.articles, &section.sections, parts);
    }
}

fn article_order(article: &JorfArticle) -> i32 {
    article
        .order
        .as_deref()
        .and_then(|order| order.parse().ok())
        .unwrap_or(-1)
}

fn section_order(section: &JorfContainerSection) -> i32 {
    match section.articles.first() {
        Some(article) => article_order(article),
        None => section.sections.first().map_or(-1, section_order),
    }
}

#[derive(Debug, Serialize, PartialEq, Eq, Clone)]
pub struct Feed {
    #[serde(rename = "@xmlns")]
    pub xmlns: String,
    pub title: String,
    pub id: String,
    pub updated: chrono::DateTime<chrono::Utc>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub logo: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub subtitle: Option<String>,
    pub author: Author,
    #[serde(rename = "link")]
    pub links: Vec<Link>,
    #[serde(rename = "entry")]
    pub entries: Vec<Entry>,
}

fn to_xml(feed: &Feed) -> Result<String, quick_xml::SeError> {
    let body = quick_xml::se::to_string_with_root("feed", feed)?;
    Ok(format!(
        "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n{body}"
    ))
}

#[derive(Debug, Serialize, PartialEq, Eq, Clone)]
pub struct Entry {
    pub title: String,
    pub id: String,
    pub updated: chrono::DateTime<chrono::Utc>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub published: Option<chrono::DateTime<chrono::Utc>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub author: Option<Author>,
    #[serde(rename = "link")]
    pub links: Vec<Link>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub summary: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub content: Option<Content>,
    #[serde(rename = "category")]
    pub categories: Vec<Category>,
}

#[derive(Debug, Serialize, PartialEq, Eq, Clone)]
pub struct Author {
    pub name: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub email: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub uri: Option<String>,
}

#[derive(Debug, Serialize, PartialEq, Eq, Clone)]
pub struct Link {
    #[serde(rename = "@href")]
    pub href: String,
    #[serde(rename = "@rel", skip_serializing_if = "Option::is_none")]
    pub rel: Option<String>,
    #[serde(rename = "@type", skip_serializing_if = "Option::is_none")]
    pub r#type: Option<String>,
}

#[derive(Debug, Serialize, PartialEq, Eq, Clone)]
pub struct Category {
    #[serde(rename = "@term")]
    pub term: String,
    #[serde(rename = "@label", skip_serializing_if = "Option::is_none")]
    pub label: Option<String>,
}

#[derive(Debug, Serialize, PartialEq, Eq, Clone)]
pub struct Content {
    #[serde(rename = "@type")]
    pub r#type: String,
    #[serde(rename = "$text")]
    pub value: String,
}

#[cfg(test)]
mod tests {
    use super::*;

    fn articles(starting_at: i32) -> Vec<JorfArticle> {
        [2, 1, 3, 4]
            .into_iter()
            .map(|offset| JorfArticle {
                content: (starting_at + offset).to_string(),
                order: Some((starting_at + offset).to_string()),
            })
            .collect()
    }

    fn section(
        title: &str,
        articles: Vec<JorfArticle>,
        sections: Vec<JorfContainerSection>,
    ) -> JorfContainerSection {
        JorfContainerSection {
            title: title.to_string(),
            articles,
            sections,
        }
    }

    fn result(
        articles: Vec<JorfArticle>,
        sections: Vec<JorfContainerSection>,
    ) -> JorfContainerResult {
        JorfContainerResult {
            id: "cid".to_string(),
            title: "title".to_string(),
            sections,
            articles,
        }
    }

    #[test]
    fn articles_are_emitted_in_numeric_order() {
        let extracted = extract_content(&result(articles(0), vec![]));
        assert_eq!(extracted, "1\n2\n3\n4");
    }

    #[test]
    fn sections_are_emitted_by_their_first_article() {
        let input = result(
            vec![],
            vec![
                section("3", articles(10), vec![]),
                section("2", vec![], vec![section("nested", articles(5), vec![])]),
                section("4", articles(15), vec![]),
                section("1", vec![], vec![]),
            ],
        );

        let extracted = extract_content(&input);
        assert_eq!(extracted, "6\n7\n8\n9\n11\n12\n13\n14\n16\n17\n18\n19");
    }

    #[test]
    fn top_level_articles_precede_sections() {
        let input = result(articles(0), vec![section("annexe", articles(10), vec![])]);

        let extracted = extract_content(&input);
        assert_eq!(extracted, "1\n2\n3\n4\n11\n12\n13\n14");
    }

    #[test]
    fn non_numeric_order_sorts_first_and_keeps_input_order() {
        let input = result(
            vec![
                JorfArticle {
                    content: "second".to_string(),
                    order: Some("2".to_string()),
                },
                JorfArticle {
                    content: "annexe".to_string(),
                    order: Some("Annexe".to_string()),
                },
                JorfArticle {
                    content: "none".to_string(),
                    order: None,
                },
                JorfArticle {
                    content: "first".to_string(),
                    order: Some("1".to_string()),
                },
            ],
            vec![],
        );

        let extracted = extract_content(&input);
        assert_eq!(extracted, "annexe\nnone\nfirst\nsecond");
    }
}
