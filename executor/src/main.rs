#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use]
extern crate rocket;
#[macro_use]
extern crate rocket_contrib;
extern crate toml;

mod config;
mod routes;

use config::language_config::LanguageConfig;

fn start() -> rocket::Rocket {
    let language_config = LanguageConfig::load("Language.toml");

    rocket::ignite()
        .manage(language_config)
        .mount("/", routes![routes::index::index])
        .mount("/api/", routes![routes::api::exec])
}

fn main() {
    start().launch();
}
