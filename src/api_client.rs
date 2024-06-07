use cynic::{impl_scalar, QueryBuilder};
use rayon::prelude::*;
use reqwest::{blocking, header, Url};
use serde::de;
use std::fmt::Debug;
use std::fs::{self, File};
use std::io::copy;
use std::path::Path;
use time::{macros::format_description, PrimitiveDateTime, Time};

pub struct Client {
    http_client: blocking::Client,
    base_url: Url,
    assets_queue: Vec<Asset>,
}

struct Asset {
    id: String,
    extension: &'static str,
    key: Option<&'static str>,
}

impl Client {
    pub fn build(base_url: Url, api_token: &str) -> Self {
        let mut headers = header::HeaderMap::new();
        let mut auth_value = header::HeaderValue::from_str(&("Bearer ".to_string() + api_token))
            .expect("Unable to set authentication header");
        auth_value.set_sensitive(true);
        headers.insert(header::AUTHORIZATION, auth_value);
        let http_client = reqwest::blocking::Client::builder()
            .default_headers(headers)
            .connection_verbose(true)
            .build()
            .expect("Failed to create reqwest client");

        Client {
            http_client,
            base_url,
            assets_queue: vec![],
        }
    }

    pub fn get_general_settings(&self) -> GeneralSettings {
        use cynic::http::ReqwestBlockingExt;
        let graphql_url = self.base_url.join("/graphql").unwrap();
        let resp = self
            .http_client
            .post(graphql_url)
            .run_graphql(GetGeneralSettings::build(()))
            .expect("Failed to fetch general settings");
        if let Some(errors) = resp.errors {
            for e in errors {
                log::error!("GraphQL query GeneralSettingsQuery returned error(s): {e}")
            }
        };
        resp.data
            .expect("No general settings returned")
            .general_settings
            .unwrap()
    }

    pub fn get_carousel_images(&self) -> Vec<String> {
        use cynic::http::ReqwestBlockingExt;
        let graphql_url = self.base_url.join("/graphql").unwrap();
        let resp = self
            .http_client
            .post(graphql_url)
            .run_graphql(GetCarousel::build(()))
            .expect("Failed to fetch carousel");
        if let Some(errors) = resp.errors {
            for e in errors {
                log::error!("GraphQL query GetCarousel returned error(s): {e}")
            }
        };
        resp.data
            .expect("No carousel returned")
            .carousel
            .unwrap()
            .images
            .into_iter()
            .map(|i| i.directus_files_id.unwrap())
            .collect()
    }

    pub fn get_artists(&self) -> Vec<Artist> {
        use cynic::http::ReqwestBlockingExt;
        let graphql_url = self.base_url.join("/graphql").unwrap();
        let resp = self
            .http_client
            .post(graphql_url)
            .run_graphql(GetArtists::build(()))
            .expect("Failed to fetch artists");
        if let Some(errors) = resp.errors {
            for e in errors {
                log::error!("GraphQL query GetArtists returned error(s): {e}")
            }
        };
        resp.data.expect("No artists returned").artists
    }

    pub fn get_sponsors(&self) -> Sponsors {
        use cynic::http::ReqwestBlockingExt;
        let graphql_url = self.base_url.join("/graphql").unwrap();
        let resp = self
            .http_client
            .post(graphql_url)
            .run_graphql(GetSponsors::build(()))
            .expect("Failed to fetch sponsors");
        if let Some(errors) = resp.errors {
            for e in errors {
                log::error!("GraphQL query GetSponsors returned error(s): {e}")
            }
        };
        let mut main = vec![];
        let mut regular = vec![];
        for sponsor in resp
            .data
            .expect("No sponsors returned")
            .sponsors
            .into_iter()
        {
            if sponsor.is_main_sponsor {
                main.push(sponsor);
            } else {
                regular.push(sponsor);
            }
        }
        return Sponsors { main, regular };
    }

    pub fn get_stages(&self) -> Vec<Stage> {
        use cynic::http::ReqwestBlockingExt;
        let graphql_url = self.base_url.join("/graphql").unwrap();
        let resp = self
            .http_client
            .post(graphql_url)
            .run_graphql(GetStages::build(()))
            .expect("Failed to fetch stages");
        if let Some(errors) = resp.errors {
            for e in errors {
                log::error!("GraphQL query GetStages returned error(s): {e}")
            }
        };
        resp.data.expect("No stages returned").stages
    }

    pub fn queue_asset(&mut self, id: String, extension: &'static str, key: Option<&'static str>) {
        self.assets_queue.push(Asset { id, extension, key })
    }

    pub fn download_assets_queue<P: AsRef<Path> + Sync>(
        &self,
        output_dir: P,
        cache_dir: Option<P>,
    ) {
        fs::create_dir_all(&output_dir).expect("Failed to create output assets dir");
        if let Some(path_cache_assets) = &cache_dir {
            log::info!(
                "Caching enabled to folder \"{}\"",
                cache_dir.as_ref().unwrap().as_ref().display()
            );
            fs::create_dir_all(&path_cache_assets).expect("Failed to create assets cache dir");
        }

        self.assets_queue.par_iter().for_each(|asset| {
            self.download_asset(
                output_dir.as_ref(),
                &asset.id,
                &asset.extension,
                asset.key,
                cache_dir.as_ref().map(|d| d.as_ref()),
            )
        });
    }

    fn download_asset<P: AsRef<Path>>(
        &self,
        output_dir: P,
        id: &str,
        extension: &str,
        key: Option<&str>,
        cache_dir: Option<P>,
    ) {
        // Derive filename and paths
        let filename = format!(
            "{}{}.{}",
            id,
            key.map_or("".to_string(), |k| "-".to_string() + k),
            extension
        );
        let output_path = output_dir.as_ref().join(&filename);
        let cache_path = cache_dir.map(|p| p.as_ref().join(&filename));

        // Try get asset from cache if enabled
        if let Some(cache_path) = &cache_path {
            if cache_path
                .try_exists()
                .expect("Unable to check if validated file exists")
            {
                log::info!("Reusing asset {} from cache ...", filename);
                fs::copy(cache_path, output_path).expect("Failed to copy cached version");
                return;
            }
        }

        // Fetch asset
        let asset_url = self
            .base_url
            .join(&format!("/assets/{}/{}", id, filename))
            .unwrap();
        log::info!("Downloading asset {} ...", filename);
        let mut req = self.http_client.get(asset_url.as_ref());
        if let Some(key) = key {
            req = req.query(&[("key", key), ("download", "true")]);
        }
        let mut resp = req.send().expect("Unable to get asset");
        if !resp.status().is_success() {
            log::error!(
                "Failed to fetch asset from {}: {}",
                asset_url,
                resp.text().unwrap()
            );
            panic!("Failed to fetch asset")
        }
        {
            let mut file = File::create(&output_path)
                .expect(&format!("Unable to create file: {}", output_path.display()));
            let file_size = copy(&mut resp, &mut file).expect("Unable to write asset to file");
            if file_size == 0 {
                log::error!(
                    "Fetching asset {} (status {}) returned an empty body",
                    asset_url,
                    resp.status().as_u16(),
                );
                panic!("Fetched asset is empty")
            }
            file.sync_all().expect("Failed to flush asset file");
        }

        // Feed cache if enabled
        if let Some(cache_path) = &cache_path {
            fs::copy(output_path, cache_path).expect("Failed to copy asset to cache");
        }
    }
}

// Generated with https://generator.cynic-rs.dev/
#[cynic::schema("directus")]
mod schema {}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "Query")]
pub struct GetGeneralSettings {
    #[cynic(rename = "general_settings")]
    pub general_settings: Option<GeneralSettings>,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "general_settings")]
pub struct GeneralSettings {
    #[cynic(rename = "address_city")]
    pub address_city: String,
    #[cynic(rename = "address_housenumber_box")]
    pub address_housenumber_box: String,
    #[cynic(rename = "address_postal_code")]
    pub address_postal_code: String,
    #[cynic(rename = "address_street")]
    pub address_street: String,
    #[cynic(rename = "email_artists")]
    pub email_artists: Option<String>,
    #[cynic(rename = "email_info")]
    pub email_info: String,
    #[cynic(rename = "email_privacy")]
    pub email_privacy: String,
    #[cynic(rename = "facebook_link")]
    pub facebook_link: Option<String>,
    #[cynic(rename = "instagram_link")]
    pub instagram_link: Option<String>,
    #[cynic(rename = "is_section_artists_visible")]
    pub is_section_artists_visible: Option<bool>,
    #[cynic(rename = "is_section_sponsors_visible")]
    pub is_section_sponsors_visible: Option<bool>,
    #[cynic(rename = "is_section_timetable_visible")]
    pub is_section_timetable_visible: Option<bool>,
    #[cynic(rename = "organisation_name")]
    pub organisation_name: String,
    #[cynic(rename = "phone_number")]
    pub phone_number: String,
    #[cynic(rename = "privacy_complaint_link")]
    pub privacy_complaint_link: String,
    #[cynic(rename = "saturday_end")]
    pub saturday_end: GraphDateTime,
    #[cynic(rename = "saturday_start")]
    pub saturday_start: GraphDateTime,
    #[cynic(rename = "sunday_end")]
    pub sunday_end: GraphDateTime,
    #[cynic(rename = "sunday_start")]
    pub sunday_start: GraphDateTime,
}

impl_scalar!(GraphDateTime, schema::Date);

#[derive(Debug)]
pub struct GraphDateTime(pub PrimitiveDateTime);

impl<'de> de::Deserialize<'de> for GraphDateTime {
    fn deserialize<D>(deserializer: D) -> Result<GraphDateTime, D::Error>
    where
        D: de::Deserializer<'de>,
    {
        let value = String::deserialize(deserializer)?;
        let format = format_description!("[year]-[month]-[day]T[hour]:[minute]:[second]");
        PrimitiveDateTime::parse(&value, format)
            .map(|d| GraphDateTime(d))
            .map_err(|e| {
                print!("Converting {:?} into DateTime returned {:?}", &value, e);
                de::Error::custom(e)
            })
    }
}

pub struct Sponsors {
    pub main: Vec<Sponsor>,
    pub regular: Vec<Sponsor>,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "Query")]
pub struct GetSponsors {
    #[arguments(sort: "sort")]
    pub sponsors: Vec<Sponsor>,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "sponsors")]
pub struct Sponsor {
    #[cynic(rename = "is_main_sponsor")]
    pub is_main_sponsor: bool,
    pub link: String,
    pub logo: Option<String>,
    pub name: String,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "Query")]
pub struct GetStages {
    pub stages: Vec<Stage>,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "stages")]
pub struct Stage {
    #[cynic(flatten)]
    pub performances: Vec<Performance>,
    pub name: String,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "performances")]
pub struct Performance {
    pub start: GraphTime,
    pub end: GraphTime,
    pub artist: Option<PerformanceArtist>,
}

impl_scalar!(GraphTime, schema::String);

#[derive(Debug)]
pub struct GraphTime(pub Time);

impl<'de> de::Deserialize<'de> for GraphTime {
    fn deserialize<D>(deserializer: D) -> Result<GraphTime, D::Error>
    where
        D: de::Deserializer<'de>,
    {
        let value = String::deserialize(deserializer)?;
        let format = format_description!("[hour]:[minute]:[second]");
        Time::parse(&value, format)
            .map(|t| GraphTime(t))
            .map_err(|e| {
                print!("Converting {:?} into Time returned {:?}", &value, e);
                de::Error::custom(e)
            })
    }
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "artists")]
pub struct PerformanceArtist {
    pub name: String,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "Query")]
pub struct GetArtists {
    pub artists: Vec<Artist>,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "artists")]
pub struct Artist {
    #[cynic(rename = "artist_links", flatten)]
    pub artist_links: Vec<ArtistLink>,
    #[cynic(rename = "artist_type")]
    pub artist_type: String,
    pub genre: Option<String>,
    pub description: Option<String>,
    pub name: String,
    pub picture: Option<String>,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "artist_links")]
pub struct ArtistLink {
    pub link: String,
    pub icon: String,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "Query")]
pub struct GetCarousel {
    pub carousel: Option<Carousel>,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "carousel")]
pub struct Carousel {
    #[cynic(flatten)]
    pub images: Vec<CarouselFiles>,
}

#[derive(cynic::QueryFragment, Debug)]
#[cynic(graphql_type = "carousel_files")]
pub struct CarouselFiles {
    #[cynic(rename = "directus_files_id")]
    pub directus_files_id: Option<String>,
}
