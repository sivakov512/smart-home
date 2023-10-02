use serde::{Deserialize, Serialize};

#[derive(Debug, PartialEq, Default, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
enum Mode {
    #[default]
    Cool,
    Heat,
}

impl std::fmt::Display for Mode {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let stringified = match self {
            Mode::Cool => "cool",
            Mode::Heat => "heat",
        };
        write!(f, "{}", stringified)
    }
}

#[derive(Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct State {
    is_active: bool,
    mode: Mode,
    current_temperature: f32,
    target_temperature: f32,
}

impl State {
    pub fn as_broadlink_command(&self, prefix: &str) -> String {
        format!(
            "{prefix}/is_active/{is_active}/mode/{mode}/target_temperature/{target_temperature:.1}",
            is_active = self.is_active,
            mode = self.mode.to_string(),
            target_temperature = self.target_temperature,
        )
    }
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

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::*;

    #[test]
    fn default_state_looks_as_expected() {
        let got = State::default();

        assert_eq!(
            got,
            State {
                is_active: false,
                mode: Mode::Cool,
                current_temperature: 0.0,
                target_temperature: 0.0,
            }
        )
    }

    #[rstest]
    #[case(
        State{
            is_active: false,
            mode: Mode::Cool,
            current_temperature: 0.0,
            target_temperature: 0.0,
        },
        "broadlink/ac/LivingRoom/is_active/false/mode/cool/target_temperature/0.0"
    )]
    #[case(
        State{
            is_active: true,
            mode: Mode::Heat,
            current_temperature: 20.0,
            target_temperature: 25.0,
        },
        "broadlink/ac/LivingRoom/is_active/true/mode/heat/target_temperature/25.0"
    )]
    #[case(
        State{
            is_active: true,
            mode: Mode::Cool,
            current_temperature: 20.0,
            target_temperature: 17.5,
        },
        "broadlink/ac/LivingRoom/is_active/true/mode/cool/target_temperature/17.5"
    )]
    fn creates_expected_broadlink_command(#[case] state: State, #[case] expected: &str) {
        let got = state.as_broadlink_command("broadlink/ac/LivingRoom");

        assert_eq!(got, expected)
    }
}
