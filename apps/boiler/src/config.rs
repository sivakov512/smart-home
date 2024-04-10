use serde::Deserialize;
use std::str::FromStr;

#[derive(Deserialize, PartialEq, Debug, Clone)]
pub struct Config {
    pub mqtt_broker: String,
    pub status_topic: String,
    pub update_topic: String,
}

impl Config {
    pub fn from_config_file(fpath: &str) -> anyhow::Result<Self> {
        let input = std::fs::read_to_string(fpath)?;
        Ok(input.parse::<Self>()?)
    }
}

impl FromStr for Config {
    type Err = toml::de::Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        toml::from_str(s)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::*;

    #[test]
    fn parsing_works_right() {
        let input = r#"
mqtt_broker = 'tcp://localhost:1883'
status_topic = 'home/bathroom/boiler'
update_topic = 'home/bathroom/boiler/set'
"#;

        let got = input.parse::<Config>();

        assert!(got.is_ok());
        assert_eq!(
            got.unwrap(),
            Config {
                mqtt_broker: "tcp://localhost:1883".into(),
                status_topic: "home/bathroom/boiler".into(),
                update_topic: "home/bathroom/boiler/set".into(),
            }
        )
    }

    #[rstest]
    fn parsing_errors_for_wrong_input(
        #[values(
            r#"
qtt_broker = 'tcp://localhost:1883'
status_topic = 'home/bathroom/boiler'
update_topic = 'home/bathroom/boiler/set'
"#,
            r#"
mqtt_broker = 100
status_topic = 'home/bathroom/boiler'
update_topic = 'home/bathroom/boiler/set'
"#
        )]
        input: &str,
    ) {
        let got = input.parse::<Config>();

        assert!(got.is_err())
    }

    #[test]
    fn loads_from_file_correctly() {
        let got = Config::from_config_file("./testing_fixtures/valid_config.toml");

        assert!(got.is_ok());
        assert_eq!(
            got.unwrap(),
            Config {
                mqtt_broker: "tcp://localhost:1883".into(),
                status_topic: "home/bathroom/boiler".into(),
                update_topic: "home/bathroom/boiler/set".into(),
            }
        )
    }

    #[test]
    fn errored_when_loading_bad_config() {
        let got = Config::from_config_file("./testing_fixtures/invalid_config.toml");

        assert!(got.is_err())
    }
}
