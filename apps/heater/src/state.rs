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

impl State {
    pub fn update(&mut self, updates: StateUpdate) {
        self.is_active = updates.is_active;
        self.mode = updates.mode;
        self.target_temperature = updates.target_temperature;
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
    use rstest::*;

    #[fixture]
    fn state() -> State {
        State {
            is_active: true,
            mode: Mode::Heat,
            current_temperature: 21.3,
            target_temperature: 23.5,
        }
    }

    #[fixture]
    fn state_update() -> StateUpdate {
        StateUpdate {
            is_active: false,
            mode: Mode::Idle,
            target_temperature: 28.7,
        }
    }

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

    #[rstest]
    fn state_encoded_as_expectedf(state: State) {
        let got = Vec::<u8>::from(&state);

        assert_eq!(
            String::from_utf8(got).unwrap(),
            r#"{"is_active":true,"mode":"heat","current_temperature":21.3,"target_temperature":23.5}"#
        );
    }

    #[rstest]
    fn state_decoded_as_expected(state: State) {
        let input = r#"{"is_active":true,"mode":"heat","current_temperature":21.3,"target_temperature":23.5}"#;

        let got = State::from(input.as_bytes());

        assert_eq!(got, state)
    }

    #[rstest]
    fn state_updates_as_expected(mut state: State, state_update: StateUpdate) {
        state.update(state_update);

        assert_eq!(
            state,
            State {
                is_active: false,
                mode: Mode::Idle,
                current_temperature: 21.3,
                target_temperature: 28.7
            }
        );
    }

    #[rstest]
    fn state_update_decoded_as_expected(state_update: StateUpdate) {
        let input = r#"{"is_active":false,"mode":"idle","target_temperature":28.7}"#;

        let got = StateUpdate::from(input.as_bytes());

        assert_eq!(got, state_update)
    }
}
