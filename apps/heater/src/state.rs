use serde::{Deserialize, Serialize};

#[derive(Debug, PartialEq, Default, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum Mode {
    #[default]
    Idle,
    Heat,
}

#[derive(Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct State {
    pub is_active: bool,
    pub mode: Mode,
    pub current_temperature: f32,
    pub target_temperature: f32,
}

impl From<&State> for Vec<u8> {
    fn from(input: &State) -> Self {
        serde_json::to_vec(input).unwrap()
    }
}

impl From<&[u8]> for State {
    fn from(input: &[u8]) -> Self {
        serde_json::from_slice(input).unwrap()
    }
}

#[derive(Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct StateUpdate {
    pub is_active: bool,
    pub mode: Mode,
    pub target_temperature: f32,
}

impl From<&[u8]> for StateUpdate {
    fn from(input: &[u8]) -> Self {
        serde_json::from_slice(input).unwrap()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn default_state_looks_as_expected() {
        let got = State::default();

        assert_eq!(
            got,
            State {
                is_active: false,
                mode: Mode::Idle,
                current_temperature: 0.0,
                target_temperature: 0.0,
            }
        )
    }

    #[test]
    fn state_update_decoded_as_expected() {
        let input = r#"{"is_active": true, "mode": "heat", "target_temperature": 23.5}"#;

        let got = StateUpdate::from(input.as_bytes());

        assert_eq!(
            got,
            StateUpdate {
                is_active: true,
                mode: Mode::Heat,
                target_temperature: 23.5
            }
        )
    }
}
