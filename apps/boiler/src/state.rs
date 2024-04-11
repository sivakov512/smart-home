use serde::{self, Deserialize, Serialize};

const STATE_ON: &str = "ON";
const STATE_OFF: &str = "OFF";

#[derive(Debug, PartialEq, Default, Deserialize, Serialize)]
pub struct State {
    #[serde(deserialize_with = "deserialize_state")]
    pub state: bool,
    #[serde(rename(deserialize = "linkquality"))]
    pub link_quality: u32,
}

fn deserialize_state<'de, D>(deserializer: D) -> Result<bool, D::Error>
where
    D: serde::Deserializer<'de>,
{
    match serde::Deserialize::deserialize(deserializer)? {
        STATE_ON => Ok(true),
        STATE_OFF => Ok(false),
        v => panic!("Unknown option: {}", v),
    }
}

impl From<&[u8]> for State {
    fn from(value: &[u8]) -> Self {
        serde_json::from_slice(value).unwrap()
    }
}

impl From<&State> for Vec<u8> {
    fn from(value: &State) -> Self {
        serde_json::to_vec(value).unwrap()
    }
}

impl State {
    pub fn update(&mut self, updates: &StateUpdate) {
        self.state = updates.state;
    }
}

#[derive(Debug, PartialEq, Default, Deserialize, Serialize)]
pub struct StateUpdate {
    #[serde(serialize_with = "serialize_state")]
    pub state: bool,
}

fn serialize_state<S>(v: &bool, serializer: S) -> Result<S::Ok, S::Error>
where
    S: serde::Serializer,
{
    match v {
        true => Ok(serializer.serialize_str(STATE_ON)?),
        false => Ok(serializer.serialize_str(STATE_OFF)?),
    }
}

impl From<&[u8]> for StateUpdate {
    fn from(value: &[u8]) -> Self {
        serde_json::from_slice(value).unwrap()
    }
}

impl From<&StateUpdate> for Vec<u8> {
    fn from(value: &StateUpdate) -> Self {
        serde_json::to_vec(value).unwrap()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::*;

    mod state_ {
        use super::*;

        #[test]
        fn default_looks_as_expected() {
            let got = State::default();

            assert_eq!(
                got,
                State {
                    state: false,
                    link_quality: 0
                }
            )
        }

        #[test]
        fn updates_as_expected() {
            let mut state = State {
                state: true,
                link_quality: 100,
            };

            state.update(&StateUpdate { state: false });

            assert_eq!(
                state,
                State {
                    state: false,
                    link_quality: 100
                }
            )
        }

        #[rstest]
        #[case(
            "{\"linkquality\":196,\"state\":\"OFF\"}",
            State{
                state: false,
                link_quality: 196,
            }
        )]
        #[case(
            "{\"linkquality\":0,\"state\":\"ON\"}",
            State{
                state: true,
                link_quality: 0,
            }
        )]
        fn deserialized_as_expected(#[case] input: &str, #[case] expected: State) {
            let got = State::from(input.as_bytes());

            assert_eq!(got, expected)
        }

        #[rstest]
        #[case(
            State{
                state: false,
                link_quality: 196,
            },
            "{\"state\":false,\"link_quality\":196}",
        )]
        #[case(
            State{
                state: true,
                link_quality: 0,
            },
            "{\"state\":true,\"link_quality\":0}",
        )]
        fn serialized_as_expected(#[case] input: State, #[case] expected: &str) {
            let got = Vec::from(&input);

            assert_eq!(String::from_utf8(got).unwrap(), expected)
        }
    }

    mod state_update_ {
        use super::*;

        #[rstest]
        #[case("{\"state\":false}", StateUpdate{state: false})]
        #[case("{\"state\":true}", StateUpdate{state: true})]
        fn deserialized_as_expected(#[case] input: &str, #[case] expected: StateUpdate) {
            let got = StateUpdate::from(input.as_bytes());

            assert_eq!(got, expected)
        }

        #[rstest]
        #[case(StateUpdate{state: false}, "{\"state\":\"OFF\"}")]
        #[case(StateUpdate{state: true}, "{\"state\":\"ON\"}")]
        fn serialized_as_expected(#[case] input: StateUpdate, #[case] expected: &str) {
            let got = Vec::from(&input);

            assert_eq!(String::from_utf8(got).unwrap(), expected)
        }
    }
}
