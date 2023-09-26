#![allow(dead_code)]
use crate::characteristics::Characteristics;

#[derive(Debug, PartialEq, Clone)]
pub enum Mode {
    Heating(f32),
    Cooling(f32),
}

#[derive(Default, Debug, PartialEq)]
pub struct State {
    mode: Option<Mode>,
    current_temperature: f32,
}

#[derive(Debug, PartialEq)]
pub struct AC {
    characteristics: Characteristics,
    state: State,
}

impl From<Characteristics> for AC {
    fn from(characteristics: Characteristics) -> Self {
        Self {
            characteristics,
            state: State::default(),
        }
    }
}

impl AC {
    pub fn set_mode(&mut self, mode: Option<Mode>) {
        self.state.mode = mode;
    }

    pub fn mode(&self) -> &Option<Mode> {
        &self.state.mode
    }

    pub fn current_temperature(&self) -> f32 {
        self.state.current_temperature
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::characteristics::TemperatureCharacteristics;
    use rstest::*;

    #[fixture]
    fn characteristics() -> Characteristics {
        Characteristics {
            manufacturer: "Pretty vendor".into(),
            model: "Pretty model".into(),
            name: "Pretty name".into(),
            cooling: TemperatureCharacteristics {
                min: 10.0,
                max: 20.0,
                step: 1.0,
            },
            heating: TemperatureCharacteristics {
                min: 21.5,
                max: 28.5,
                step: 0.5,
            },
        }
    }

    #[fixture]
    fn ac(characteristics: Characteristics) -> AC {
        AC {
            characteristics,
            state: State {
                mode: None,
                current_temperature: 0.0,
            },
        }
    }

    #[rstest]
    fn ac_created_from_characteristics_correctly(characteristics: Characteristics, ac: AC) {
        let got = AC::from(characteristics);

        assert_eq!(got, ac)
    }

    #[rstest]
    fn sets_mode_correctly(
        mut ac: AC,
        #[values(None, Some(Mode::Heating(26.5)), Some(Mode::Cooling(17.5)))] mode: Option<Mode>,
    ) {
        ac.set_mode(mode.clone());

        assert_eq!(ac.state.mode, mode)
    }

    #[rstest]
    fn returns_mode_from_state(
        mut ac: AC,
        #[values(None, Some(Mode::Heating(26.5)), Some(Mode::Cooling(17.5)))] mode: Option<Mode>,
    ) {
        ac.state.mode = mode.clone();

        let got = ac.mode();

        assert_eq!(got, &mode)
    }

    #[rstest]
    fn returns_current_temperature_from_state(mut ac: AC) {
        ac.state.current_temperature = 13.2;

        let got = ac.current_temperature();

        assert_eq!(got, 13.2);
    }
}
