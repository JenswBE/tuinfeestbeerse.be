use time::{ext::NumericalDuration, PrimitiveDateTime, Time};

/// Number of pixels to represent 1 minute on the timetable
const PIXELS_PER_MINUTE: f64 = 1.2;
/// Number of minutes for each timeslot
const TIMESLOT_MINUTES: i64 = 30;

#[derive(PartialEq, Debug)]
pub struct Timetable {
    pub stages: Vec<Stage>,
    pub slots: Vec<String>,
    pub slot_height_pixels: u32,
}

#[derive(PartialEq, Debug)]
pub struct Stage {
    pub name: String,
    pub performances: Vec<Performance>,
}

#[derive(PartialEq, Debug)]
pub struct Performance {
    pub artist_name: String,
    pub start: Time,
    pub end: Time,
    pub start_pixels: u32,
    pub height_pixels: u32,
}

impl Timetable {
    pub fn build(
        event_start: PrimitiveDateTime,
        event_end: PrimitiveDateTime,
        api_stages: Vec<crate::api_client::Stage>,
    ) -> Self {
        // Calculate slots
        let slot_count = (event_end - event_start).whole_minutes() / TIMESLOT_MINUTES + 1;
        let mut slots = Vec::with_capacity(slot_count.try_into().unwrap());
        let mut slot_start = event_start;
        while slot_start <= event_end {
            slots.push(slot_start.format(crate::TIME_FORMAT_KITCHEN).unwrap());
            slot_start = slot_start.saturating_add(TIMESLOT_MINUTES.minutes());
        }
        let slot_height_pixels = TIMESLOT_MINUTES as f64 * PIXELS_PER_MINUTE;

        // Calculate performance starts and heights
        let mut stages = Vec::with_capacity(api_stages.len());
        for api_stage in api_stages {
            let api_performances = api_stage.performances;
            let mut stage = Stage {
                name: api_stage.name,
                performances: Vec::with_capacity(api_performances.len()),
            };
            for api_performance in api_performances.into_iter() {
                let api_performance = api_performance;
                let performance_start = calc_time(&event_start, &api_performance.start.0);
                if performance_start < event_start {
                    panic!("Start of performance cannot be before start of event");
                }
                let performance_end = calc_time(&event_start, &api_performance.end.0);
                if performance_end > event_end {
                    panic!("End of performance cannot be after end of event");
                }
                stage.performances.push(Performance {
                    artist_name: api_performance.artist.unwrap().name,
                    start: performance_start.time(),
                    end: performance_end.time(),
                    // Timeslot indicator is in the middle of the row instead of the top.
                    // Therefore, adding half of a timeslot row to the start pixels.
                    start_pixels: (((performance_start - event_start).whole_minutes() as f64
                        + TIMESLOT_MINUTES as f64 / 2.0)
                        * PIXELS_PER_MINUTE) as u32,
                    height_pixels: ((performance_end - performance_start).whole_minutes() as f64
                        * PIXELS_PER_MINUTE) as u32
                        - 3, // Add 3px spacing around shows
                })
            }
            stages.push(stage);
        }

        // Build timetable
        Timetable {
            stages,
            slots,
            slot_height_pixels: slot_height_pixels as u32,
        }
    }
}

/// calc_time calculates the time of the performance based on the start
/// date of the event and the start time of the performance.
/// Times between 00:00 and 07:00 are considered the next day.
fn calc_time(event_start: &PrimitiveDateTime, performance_time: &Time) -> PrimitiveDateTime {
    if performance_time.hour() >= 7 {
        // Performances after 07:00 are on the start day
        event_start.replace_time(performance_time.clone())
    } else {
        // Performances before 07:00 are on the next day
        event_start
            .saturating_add(1.days())
            .replace_time(performance_time.clone())
    }
}

#[cfg(test)]
mod tests {
    use crate::api_client;
    use pretty_assertions::assert_eq;
    use time::macros::{datetime, time};

    use super::*;

    #[test]
    fn build_timetable() {
        let event_start = datetime!(2024-06-29 20:00);
        let event_end = datetime!(2024-06-30 03:00);
        let api_stages = vec![
            api_client::Stage {
                name: "Stage A".to_string(),
                performances: vec![
                    api_client::Performance {
                        start: api_client::GraphTime(time!(21:00)),
                        end: api_client::GraphTime(time!(22:30)),
                        artist: Some(api_client::PerformanceArtist {
                            name: "Artist A1".to_string(),
                        }),
                    },
                    api_client::Performance {
                        start: api_client::GraphTime(time!(23:00)),
                        end: api_client::GraphTime(time!(23:40)),
                        artist: Some(api_client::PerformanceArtist {
                            name: "Artist A2".to_string(),
                        }),
                    },
                    api_client::Performance {
                        start: api_client::GraphTime(time!(00:10)),
                        end: api_client::GraphTime(time!(01:40)),
                        artist: Some(api_client::PerformanceArtist {
                            name: "Artist A3".to_string(),
                        }),
                    },
                ],
            },
            api_client::Stage {
                name: "Stage B".to_string(),
                performances: vec![
                    api_client::Performance {
                        start: api_client::GraphTime(time!(20:00)),
                        end: api_client::GraphTime(time!(21:00)),
                        artist: Some(api_client::PerformanceArtist {
                            name: "Artist B1".to_string(),
                        }),
                    },
                    api_client::Performance {
                        start: api_client::GraphTime(time!(21:00)),
                        end: api_client::GraphTime(time!(22:00)),
                        artist: Some(api_client::PerformanceArtist {
                            name: "Artist B2".to_string(),
                        }),
                    },
                ],
            },
        ];

        let timetable = Timetable::build(event_start, event_end, api_stages);

        // Expected value is based on following constants:
        //   PIXELS_PER_MINUTE (ppm) = 1.2
        //   TIMESLOT_MINUTES (tsm)  = 30
        let expected = Timetable {
            stages: vec![
                Stage {
                    name: "Stage A".to_string(),
                    performances: vec![
                        Performance {
                            artist_name: "Artist A1".to_string(),
                            start: time!(21:00),
                            end: time!(22:30),
                            start_pixels: 90,   // ( 60 mins + 30 tsm / 2 ) * 1.2 ppm
                            height_pixels: 105, // 90 mins * 1.2 ppm - 3
                        },
                        Performance {
                            artist_name: "Artist A2".to_string(),
                            start: time!(23:00),
                            end: time!(23:40),
                            start_pixels: 234, // ( 180 mins + 30 tsm / 2 ) * 1.2 ppm
                            height_pixels: 45, // 40 mins * 1.2 ppm - 3
                        },
                        Performance {
                            artist_name: "Artist A3".to_string(),
                            start: time!(00:10),
                            end: time!(01:40),
                            start_pixels: 318, // ( 250 mins + 30 tsm / 2 ) * 1.2 ppm
                            height_pixels: 105, // 90 mins * 1.2 ppm - 3
                        },
                    ],
                },
                Stage {
                    name: "Stage B".to_string(),
                    performances: vec![
                        Performance {
                            artist_name: "Artist B1".to_string(),
                            start: time!(20:00),
                            end: time!(21:00),
                            start_pixels: 18,  // ( 0 mins + 30 tsm / 2 ) * 1.2 ppm
                            height_pixels: 69, // 60 mins * 1.2 ppm - 3
                        },
                        Performance {
                            artist_name: "Artist B2".to_string(),
                            start: time!(21:00),
                            end: time!(22:00),
                            start_pixels: 90,  // ( 60 mins + 30 tsm / 2 ) * 1.2 ppm
                            height_pixels: 69, // 60 mins * 1.2 ppm - 3
                        },
                    ],
                },
            ],
            slots: vec![
                "20:00".to_string(),
                "20:30".to_string(),
                "21:00".to_string(),
                "21:30".to_string(),
                "22:00".to_string(),
                "22:30".to_string(),
                "23:00".to_string(),
                "23:30".to_string(),
                "00:00".to_string(),
                "00:30".to_string(),
                "01:00".to_string(),
                "01:30".to_string(),
                "02:00".to_string(),
                "02:30".to_string(),
                "03:00".to_string(),
            ],
            slot_height_pixels: 36, // 30 tsm * 1.2 ppm
        };

        assert_eq!(timetable, expected);
    }
}
