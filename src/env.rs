use crate::model::Config;

pub fn load_env(config_file: String) -> Config {
    dotenvy::from_path(config_file).ok();
    Config {
        database_url: get_or_fail("DATABASE_URL"),
        oauth_url: get_or_fail("OAUTH_URL"),
        api_url: get_or_fail("API_URL"),
        client_secret: get_or_fail("CLIENT_SECRET"),
        client_id: get_or_fail("CLIENT_ID"),
    }
}

fn get_or_fail(variable: &str) -> String {
    dotenvy::var(variable).unwrap_or_else(|_| panic!("Fail to get {variable}"))
}
