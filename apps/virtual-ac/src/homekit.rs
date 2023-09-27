use crate::ac::{Mode, AC};
use hap::accessory::{thermostat::ThermostatAccessory, AccessoryInformation};
use hap::characteristic::{CharacteristicCallbacks, HapCharacteristic};
use hap::service::thermostat::ThermostatService;

pub fn build_accessory(ac: &AC) -> ThermostatAccessory {
    let mut accessory = ThermostatAccessory::new(
        1,
        AccessoryInformation {
            manufacturer: ac.characteristics.manufacturer.clone(),
            model: ac.characteristics.model.clone(),
            name: ac.characteristics.name.clone(),
            ..Default::default()
        },
    )
    .unwrap();

    let service = &mut accessory.thermostat;

    service.target_relative_humidity = None;
    service.current_relative_humidity = None;
    service.cooling_threshold_temperature = None;
    service.heating_threshold_temperature = None;

    configure_target_temperature(service, ac).unwrap();
    configure_heating_cooling_state(service, ac);
    configure_current_temperature(service, ac);

    accessory
}

fn configure_target_temperature(service: &mut ThermostatService, ac: &AC) -> anyhow::Result<()> {
    let cloned_ac = ac.clone();
    service
        .target_temperature
        .on_read(Some(move || Ok(Some(cloned_ac.get_target_temperature()))));

    let cloned_ac = ac.clone();
    service
        .target_temperature
        .on_update(Some(move |_: &f32, new: &f32| {
            cloned_ac.set_target_temperature(*new);
            Ok(())
        }));

    service
        .target_temperature
        .set_min_value(Some(ac.characteristics.temperature.min.into()))?;
    service
        .target_temperature
        .set_max_value(Some(ac.characteristics.temperature.max.into()))?;
    service
        .target_temperature
        .set_step_value(Some(ac.characteristics.temperature.step.into()))?;

    Ok(())
}

fn configure_heating_cooling_state(service: &mut ThermostatService, ac: &AC) {
    let cloned_ac = ac.clone();
    let read_state = Some(move || Ok(Some(cloned_ac.get_mode().to_hap())));
    service
        .current_heating_cooling_state
        .on_read(read_state.clone());
    service.target_heating_cooling_state.on_read(read_state);

    let cloned_ac = ac.clone();
    let update_state = Some(move |_: &u8, new: &u8| {
        cloned_ac.set_mode(Mode::from_hap(*new));
        Ok(())
    });
    service
        .current_heating_cooling_state
        .on_update(update_state.clone());
    service.target_heating_cooling_state.on_update(update_state);
}

fn configure_current_temperature(service: &mut ThermostatService, ac: &AC) {
    let cloned_ac = ac.clone();
    service
        .current_temperature
        .on_read(Some(move || Ok(Some(cloned_ac.get_current_temperature()))));
}

impl Mode {
    fn from_hap(mode: u8) -> Mode {
        match mode {
            0 => Mode::Off,
            1 => Mode::Heat,
            2 => Mode::Cool,
            3 => Mode::Auto,
            _ => panic!("Unexpected mode from HAP {:?}", mode),
        }
    }

    fn to_hap(&self) -> u8 {
        match self {
            Self::Off => 0,
            Self::Heat => 1,
            Self::Cool => 2,
            Self::Auto => 3,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::ac::AC;
    use crate::characteristics::{Characteristics, TemperatureCharacteristics};
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
        AC::from(characteristics)
    }

    #[rstest]
    fn will_setup_target_temperature_correctly(ac: AC) {
        let got = build_accessory(&ac);

        let threshold = got.thermostat.target_temperature;
        assert_eq!(threshold.get_min_value(), Some(serde_json::json!(10.0)));
        assert_eq!(threshold.get_max_value(), Some(serde_json::json!(20.0)));
        assert_eq!(threshold.get_step_value(), Some(serde_json::json!(1.0)));
    }
}
