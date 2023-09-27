use serde::Deserialize;
use std::str::FromStr;

#[derive(Deserialize, PartialEq, Debug, Clone)]
pub struct Characteristics {
    pub manufacturer: String,
    pub model: String,
    pub name: String,
    pub temperature: TemperatureCharacteristics,
}

impl Characteristics {
    pub fn from_config_file(file_path: &str) -> anyhow::Result<Self> {
        let config = std::fs::read_to_string(file_path)?;
        Ok(config.parse::<Self>()?)
    }
}

impl FromStr for Characteristics {
    type Err = toml::de::Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        toml::from_str(s)
    }
}

#[derive(Deserialize, PartialEq, Debug, Clone)]
pub struct TemperatureCharacteristics {
    pub min: f32,
    pub max: f32,
    pub step: f32,
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::*;

    #[test]
    fn characteristics_parsing_works_right() {
        let input = r#"
manufacturer = 'Pretty vendor'
model = 'Pretty model'
name = "Pretty name"

[temperature]
min = 21.5
max = 28.5
step = 0.5
"#;

        let got = input.parse::<Characteristics>();

        assert_eq!(
            got,
            Ok(Characteristics {
                manufacturer: "Pretty vendor".into(),
                model: "Pretty model".into(),
                name: "Pretty name".into(),
                temperature: TemperatureCharacteristics {
                    min: 21.5,
                    max: 28.5,
                    step: 0.5
                }
            })
        )
    }

    #[rstest]
    fn errored_if_config_is_wrong(
        #[values(
            r#"
manufacturer = 100
model = 'Pretty model'
name = "Pretty name"

[temperature]
min = 10.5
max = 20
step = 1.1
"#,
            r#"
model = 'Pretty model'
name = "Pretty name"

[temperature]
min = 10
max = 20
step = 1
"#,
            r#"
manufacturer = 'Pretty vendor'
model = 'Pretty model'
name = "Pretty name"

min = 10
max = 20
step = 1
"#,
            r#"
manufacturer = 'Pretty vendor'
model = 'Pretty model'
name = "Pretty name"

[temperature]
min = hundred
max = 20
step = 1
"#
        )]
        input: &str,
    ) {
        let got = input.parse::<Characteristics>();

        assert!(got.is_err())
    }
}
