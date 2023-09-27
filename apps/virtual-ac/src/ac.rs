use crate::characteristics::Characteristics;
use std::sync::{Arc, Mutex};

#[derive(Debug, PartialEq, Clone, Default)]
pub enum Mode {
    #[default]
    Off,
    Heat,
    Cool,
    Auto,
}

#[derive(Default, Debug, PartialEq)]
struct State {
    mode: Mode,
    target_temperature: f32,
    current_temperature: f32,
}

#[derive(Debug, Clone)]
pub struct AC {
    pub characteristics: Characteristics,
    state: Arc<Mutex<State>>,
}

impl From<Characteristics> for AC {
    fn from(characteristics: Characteristics) -> Self {
        Self {
            characteristics,
            state: Arc::new(Mutex::new(State::default())),
        }
    }
}

impl AC {
    fn lock_state(&self) -> std::sync::MutexGuard<State> {
        self.state.lock().unwrap()
    }

    pub fn set_mode(&self, mode: Mode) {
        self.lock_state().mode = mode;
        dbg!(self.lock_state());
    }

    pub fn get_mode(&self) -> Mode {
        self.lock_state().mode.clone()
    }

    pub fn get_current_temperature(&self) -> f32 {
        self.lock_state().current_temperature
    }

    pub fn get_target_temperature(&self) -> f32 {
        self.lock_state().target_temperature
    }

    pub fn set_target_temperature(&self, target_temperature: f32) {
        self.lock_state().target_temperature = target_temperature;
        dbg!(self.lock_state());
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
            temperature: TemperatureCharacteristics {
                min: 10.0,
                max: 20.0,
                step: 1.0,
            },
        }
    }

    #[fixture]
    fn ac(characteristics: Characteristics) -> AC {
        AC {
            characteristics,
            state: Arc::new(Mutex::new(State {
                mode: Mode::Off,
                target_temperature: 0.0,
                current_temperature: 0.0,
            })),
        }
    }

    #[rstest]
    fn ac_created_from_characteristics_correctly(characteristics: Characteristics, ac: AC) {
        let got = AC::from(characteristics);

        assert_eq!(got.characteristics, ac.characteristics);
        let lhs_state = got.state.lock().unwrap();
        assert_eq!(*lhs_state, *ac.state.lock().unwrap())
    }

    #[rstest]
    fn sets_mode_correctly(ac: AC, #[values(Mode::Off, Mode::Heat, Mode::Cool)] mode: Mode) {
        ac.set_mode(mode.clone());

        assert!(ac.state.lock().is_ok());
        assert_eq!(ac.state.lock().unwrap().mode, mode)
    }

    #[rstest]
    fn returns_mode_from_state(ac: AC, #[values(Mode::Off, Mode::Heat, Mode::Cool)] mode: Mode) {
        ac.state.lock().unwrap().mode = mode.clone();

        let got = ac.get_mode();

        assert_eq!(got, mode)
    }

    #[rstest]
    fn returns_current_temperature_from_state(ac: AC) {
        ac.state.lock().unwrap().current_temperature = 13.2;

        let got = ac.get_current_temperature();

        assert_eq!(got, 13.2);
    }

    #[rstest]
    fn returns_target_temperature_from_state(ac: AC) {
        ac.state.lock().unwrap().target_temperature = 13.2;

        let got = ac.get_target_temperature();

        assert_eq!(got, 13.2);
    }

    #[rstest]
    fn sets_target_temperature_correctly(ac: AC) {
        ac.set_target_temperature(25.5);

        assert_eq!(ac.state.lock().unwrap().target_temperature, 25.5);
    }
}
