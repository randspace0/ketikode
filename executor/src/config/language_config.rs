use serde::Deserialize;
use std::collections::HashMap;
use std::io::prelude::*;

#[derive(Deserialize, Debug)]
pub struct Language {
    pub extension: String,
    pub compiler: Option<String>,
    pub flag: Option<String>,
    pub command: String,
}

pub struct LanguageConfig(HashMap<String, Language>);

impl LanguageConfig {
    pub fn load(file_name: &str) -> LanguageConfig {
        let mut language_config_file =
            std::fs::File::open(file_name).expect("Error when trying to open Language.toml");

        let mut file_content_buffer = String::new();

        language_config_file
            .read_to_string(&mut file_content_buffer)
            .expect("Error when copying language config file content to the string buffer");

        let config: HashMap<String, Language> =
            toml::from_str(&file_content_buffer).expect("Configuration error");

        LanguageConfig(config)
    }

    pub fn of(&self, language_name: &str) -> Option<&Language> {
        self.0.get(language_name)
    }
}
