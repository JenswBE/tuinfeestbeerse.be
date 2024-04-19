mod api_client;
mod renderer;
mod timetable;

use crate::api_client::Client;
use crate::renderer::Renderer;
use crate::timetable::Timetable;
use api_client::Sponsors;
use askama::Template;
use dotenv::dotenv;
use reqwest::Url;
use std::path::Path;
use std::process::{self, Command};
use std::{env, fs, io};
use time::macros::offset;
use time::OffsetDateTime;
use time::{
    format_description::{well_known::Rfc3339, BorrowedFormatItem},
    macros::format_description,
};

const LOCAL_BASE_URL: &'static str = "http://localhost:8000";
const LOCAL_API_BASE_URL: &'static str = "http://localhost:8055";
const LOCAL_API_KEY: &'static str = "YR9_OxK1MHywW71CZ9tG-SFnDD65qnTo";
const TIME_FORMAT_KITCHEN: &[BorrowedFormatItem<'_>] = format_description!("[hour]:[minute]");
const TIME_FORMAT_RFC3339: &Rfc3339 = &Rfc3339;

/// Fields present in each template with the same value.
struct TemplateBaseCommon {
    email_artists: Option<String>,
    email_info: String,
    facebook_link: Option<String>,
    instagram_link: Option<String>,
    now_year: i32,
}

/// Fields present in each template but with a different value.
struct TemplateBaseSpecific<'a> {
    canonical_url: Option<&'a str>,
    title: Option<String>,
}

#[derive(Template)]
#[template(path = "index.html.jinja2", ext = "html")]
struct TemplateIndex<'a> {
    base_common: &'a TemplateBaseCommon,
    base_specific: TemplateBaseSpecific<'a>,
    show_timetable: bool,
    show_artists: bool,
    show_sponsors: bool,
    saturday_start: OffsetDateTime,
    saturday_end: OffsetDateTime,
    sunday_start: OffsetDateTime,
    sunday_end: OffsetDateTime,
    carousel_images: Vec<String>,
    timetable: Timetable,
    artists: Vec<Artist>,
    sponsors: Sponsors,
}

pub struct Artist {
    pub name: String,
    pub artist_type: String,
    pub genre: Option<String>,
    pub description: Option<String>,
    pub picture_path: String,
    pub picture_alt: String,
    pub artist_links: Vec<api_client::ArtistLink>,
}

impl From<api_client::Artist> for Artist {
    fn from(a: api_client::Artist) -> Self {
        let picture_path;
        let picture_alt;
        match a.picture {
            Some(img_id) => {
                picture_path = format!("/assets/{}-artist.jpg", img_id);
                picture_alt = format!("Foto van artiest {}", &a.name);
            }
            None => match a.artist_type.as_str() {
                "band" => {
                    picture_path = "/assets/images/band.jpg".to_string();
                    picture_alt = "Stock-foto van een gitaar".to_string();
                }
                "dj" => {
                    picture_path = "/assets/images/dj.jpg".to_string();
                    picture_alt = "Stock-foto van een mengpaneel van een DJ".to_string();
                }
                artist_type => panic!("{} is not a valid artist type", artist_type),
            },
        }

        Self {
            name: a.name,
            artist_type: a.artist_type,
            genre: a.genre,
            description: a.description,
            artist_links: a.artist_links,
            picture_path,
            picture_alt,
        }
    }
}

#[derive(Template)]
#[template(path = "huisreglement.html.jinja2", ext = "html")]
struct TemplateHuisreglement<'a> {
    base_common: &'a TemplateBaseCommon,
    base_specific: TemplateBaseSpecific<'a>,
}

#[derive(Template)]
#[template(path = "vrijwilligers_privacy.html.jinja2", ext = "html")]
struct TemplateVrijwilligersPrivacy<'a> {
    base_common: &'a TemplateBaseCommon,
    base_specific: TemplateBaseSpecific<'a>,
    organisation_name: String,
    address_street: String,
    address_housenumber_box: String,
    address_postal_code: String,
    address_city: String,
    email: String,
    phone_number: String,
    complaint_link: String,
}

#[derive(Template)]
#[template(path = "404.html.jinja2", ext = "html")]
struct Template404<'a> {
    base_common: &'a TemplateBaseCommon,
    base_specific: TemplateBaseSpecific<'a>,
}

fn main() {
    // Load .env file
    dotenv().ok();

    // Setup logger
    env_logger::init_from_env(
        env_logger::Env::default().filter_or(env_logger::DEFAULT_FILTER_ENV, "info"),
    );

    // Collect env vars
    let base_url = env_var_with_default("TUINFEEST_BASE_URL", LOCAL_BASE_URL);
    let api_base_url = env_var_with_default("TUINFEEST_API_BASE_URL", LOCAL_API_BASE_URL);
    let api_key = env_var_with_default("TUINFEEST_API_KEY", LOCAL_API_KEY);
    let path_cache = env::var("TUINFEEST_CACHE_DIR").ok();

    // Create HTTP client
    let api_base_url = Url::parse(&api_base_url).unwrap();
    let mut client = Client::build(api_base_url, &api_key);

    // Fetch remote data
    let general_settings = client.get_general_settings();
    let carousel_images = client.get_carousel_images();
    let artists = client.get_artists();
    let sponsors = client.get_sponsors();
    let stages = client.get_stages();

    // Prepare output dir
    let path_output = Path::new("output");
    let path_static = Path::new("static");
    ensure_empty_dir(path_output).expect("Unable to ensure empty output directory");
    copy_static(&path_static.join("."), path_output).expect("Unable to copy statics");

    // Create renderer
    let base_url = Url::parse(&base_url).unwrap();
    let mut renderer = Renderer::new(base_url.clone(), path_output);

    // Base template
    let base_template_common = TemplateBaseCommon {
        email_artists: general_settings.email_artists,
        email_info: general_settings.email_info,
        facebook_link: general_settings.facebook_link,
        instagram_link: general_settings.instagram_link,
        now_year: time::OffsetDateTime::now_utc().year(),
    };

    // Generate "Index" page
    let sitemap_url = derive_sitemap_url(&base_url, "");
    for img in carousel_images.iter() {
        client.queue_asset(img.clone(), "jpg", Some("carousel"));
    }
    for artist in artists.iter() {
        if let Some(picture) = artist.picture.as_ref() {
            client.queue_asset(picture.clone(), "jpg", Some("artist"));
        };
    }
    for sponsor in sponsors.main.iter() {
        client.queue_asset(
            sponsor
                .logo
                .as_ref()
                .expect("Sponsors must have a logo")
                .clone(),
            "jpg",
            Some("sponsor-main"),
        );
    }
    for sponsor in sponsors.regular.iter() {
        client.queue_asset(
            sponsor
                .logo
                .as_ref()
                .expect("Sponsors must have a logo")
                .clone(),
            "jpg",
            Some("sponsor-regular"),
        );
    }
    renderer.render_page(
        "index.html",
        &TemplateIndex {
            base_common: &base_template_common,
            base_specific: TemplateBaseSpecific {
                canonical_url: Some(&sitemap_url),
                title: None,
            },
            show_timetable: general_settings
                .is_section_timetable_visible
                .unwrap_or_default(),
            show_artists: general_settings
                .is_section_artists_visible
                .unwrap_or_default(),
            show_sponsors: general_settings
                .is_section_sponsors_visible
                .unwrap_or_default(),
            saturday_start: general_settings.saturday_start.0.assume_offset(offset!(+2)), // Tuinfeest is always in the summer
            saturday_end: general_settings.saturday_end.0.assume_offset(offset!(+2)), // Tuinfeest is always in the summer
            sunday_start: general_settings.sunday_start.0.assume_offset(offset!(+2)), // Tuinfeest is always in the summer
            sunday_end: general_settings.sunday_end.0.assume_offset(offset!(+2)), // Tuinfeest is always in the summer
            carousel_images,
            timetable: Timetable::build(
                general_settings.saturday_start.0,
                general_settings.saturday_end.0,
                stages,
            ),
            artists: artists.into_iter().map(Artist::from).collect(),
            sponsors,
        }
        .render()
        .expect("Unable to render index template"),
        Some(sitemap_url),
    );

    // Generate "Huisreglement" page
    let sitemap_url = derive_sitemap_url(&base_url, "/huisreglement/");
    renderer.render_page(
        "huisreglement/index.html",
        &TemplateHuisreglement {
            base_common: &base_template_common,
            base_specific: TemplateBaseSpecific {
                canonical_url: Some(&sitemap_url),
                title: Some("Huisreglement".to_string()),
            },
        }
        .render()
        .expect("Unable to render huisreglement template"),
        Some(sitemap_url),
    );

    // Generate "Vrijwilligers privacy" page
    let sitemap_url = derive_sitemap_url(&base_url, "/vrijwilligers/privacy/");
    renderer.render_page(
        "vrijwilligers/privacy/index.html",
        &TemplateVrijwilligersPrivacy {
            base_common: &base_template_common,
            base_specific: TemplateBaseSpecific {
                canonical_url: Some(&sitemap_url),
                title: Some("Privacy Policy voor Vrijwilligers".to_string()),
            },
            organisation_name: general_settings.organisation_name,
            address_street: general_settings.address_street,
            address_housenumber_box: general_settings.address_housenumber_box,
            address_postal_code: general_settings.address_postal_code,
            address_city: general_settings.address_city,
            email: general_settings.email_privacy,
            phone_number: general_settings.phone_number,
            complaint_link: general_settings.privacy_complaint_link,
        }
        .render()
        .expect("Unable to render vrijwilligers privacy template"),
        Some(sitemap_url),
    );

    // Generate "404" page
    renderer.render_page(
        "404.html",
        &Template404 {
            base_common: &base_template_common,
            base_specific: TemplateBaseSpecific {
                canonical_url: None,
                title: Some("Pagina niet gevonden".to_string()),
            },
        }
        .render()
        .expect("Unable to render 404 template"),
        None,
    );

    // Write robots and sitemap
    renderer.render_robots_txt();
    renderer.render_sitemap_xml();

    // Prepare asset cache dir and download queue
    let path_assets = path_output.join("assets");
    let path_cache = path_cache.as_ref().map(Path::new);
    let path_cache_assets = path_cache.map(|p| p.join("assets"));
    client.download_assets_queue(&path_assets, path_cache_assets.as_ref())
}

fn env_var_with_default(name: &'static str, default: &'static str) -> String {
    env::var(name).unwrap_or_else(|_| {
        log::info!("Unable to read {name}. Using default: {default}");
        default.to_string()
    })
}

fn ensure_empty_dir(path: &Path) -> io::Result<()> {
    fs::create_dir_all(path)?;
    for entry in fs::read_dir(path)? {
        let entry = entry?;
        let path = entry.path();
        if entry.file_type()?.is_dir() {
            fs::remove_dir_all(path)?;
        } else {
            fs::remove_file(path)?;
        }
    }
    Ok(())
}

fn copy_static(source: &Path, target: &Path) -> io::Result<process::Output> {
    Command::new("cp")
        .args([
            "--recursive",
            "--dereference",
            "--preserve=all",
            source.to_str().unwrap(),
            target.to_str().unwrap(),
        ])
        .output()
}

fn derive_sitemap_url<'a>(base_url: &Url, path: &'a str) -> String {
    if path == "" {
        // Seems reqwest.Url always forces at least the root path "/".
        // So, trimming the trailing slash in case provided sitemap URL is empty.
        base_url
            .to_string()
            .strip_suffix(base_url.path())
            .unwrap()
            .to_string()
    } else {
        base_url
            .join(path)
            .expect("Unable to join sitemap URL with base URL")
            .to_string()
    }
}
