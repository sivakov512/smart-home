use serde::Deserialize;

#[derive(Debug, PartialEq, Default, Deserialize)]
pub struct State {
    pub battery: Option<u32>,
    pub humidity: Option<f32>,
    #[serde(rename(deserialize = "linkquality"))]
    pub link_quality: Option<u32>,
    pub temperature: Option<f32>,
    pub voltage: Option<u32>,
}

impl From<&[u8]> for State {
    fn from(input: &[u8]) -> Self {
        serde_json::from_slice(input).unwrap()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::*;

    #[rstest]
    #[case(
        "{\"battery\":100,\"humidity\":53.66,\"linkquality\":188,\"temperature\":24.04,\"voltage\":3200}",
        State{
            battery: Some(100),
            humidity: Some(53.66),
            link_quality: Some(188),
            temperature: Some(24.04),
            voltage: Some(3200)
        }
    )]
    #[case(
        "{\"linkquality\":188,\"voltage\":3200}",
        State{
            battery: None,
            humidity: None,
            link_quality: Some(188),
            temperature: None,
            voltage: Some(3200)
        }
    )]
    #[case(
        "{}",
        State{
            battery: None,
            humidity: None,
            link_quality: None,
            temperature: None,
            voltage: None
        }
    )]
    fn state_deserialized_as_expected(#[case] input: &str, #[case] expected: State) {
        let got = State::from(input.as_bytes());

        assert_eq!(got, expected)
    }
}
