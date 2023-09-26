mod characteristics;
mod ac;

use hap::{
    accessory::{heater_cooler::HeaterCoolerAccessory, AccessoryCategory, AccessoryInformation},
    characteristic::{CharacteristicCallbacks, HapCharacteristic},
    server::{IpServer, Server},
    storage::{FileStorage, Storage},
    Config, Result,
};

const AC_NAME: &str = "Midea AC";

#[tokio::main]
async fn main() -> Result<()> {
    let mut heater_cooler = HeaterCoolerAccessory::new(
        1,
        AccessoryInformation {
            manufacturer: "Midea".into(),
            model: "unknown".into(),
            name: AC_NAME.into(),
            ..Default::default()
        },
    )?;

    let mut storage = FileStorage::current_dir().await?;

    let config = match storage.load_config().await {
        Ok(mut config) => {
            config.redetermine_local_ip();
            storage.save_config(&config).await?;
            config
        }
        Err(_) => {
            let config = Config {
                name: AC_NAME.into(),
                category: AccessoryCategory::AirConditioner,
                ..Default::default()
            };
            storage.save_config(&config).await?;
            config
        }
    };

    let service = &mut heater_cooler.heater_cooler;

    service.lock_physical_controls = None;
    service.name = None;
    service.swing_mode = None;
    service.rotation_speed = None;
    service.temperature_display_units = None;

    service.active.on_read(Some(|| {
        println!("Reading: active");
        Ok(Some(1))
    }));
    service.active.set_value(1.into()).await?;

    service.current_heater_cooler_state.on_read(Some(|| {
        println!("Reading: current_heater_cooler_state");
        Ok(Some(3))
    }));
    service
        .current_heater_cooler_state
        .set_value(3.into())
        .await?;

    service.target_heater_cooler_state.on_read(Some(|| {
        println!("Reading: target_heater_cooler_state");
        Ok(Some(2))
    }));
    service
        .target_heater_cooler_state
        .set_value(2.into())
        .await?;

    service.current_temperature.on_read(Some(|| {
        println!("Reading: current_temperature");
        Ok(None)
    }));
    // service.current_temperature.set_value(18.into()).await?;

    if let Some(cooling_threshold_temperature) = service.cooling_threshold_temperature.as_mut() {
        cooling_threshold_temperature.on_read(Some(|| {
            println!("Reading: cooling_threshold_temperature");
            Ok(Some(15_f32))
        }));
        cooling_threshold_temperature.set_value(15.into()).await?;
        cooling_threshold_temperature.set_step_value(Some(1.into()))?;
        cooling_threshold_temperature.set_min_value(Some(12.into()))?;
        cooling_threshold_temperature.set_max_value(Some(25.into()))?;
    }

    if let Some(heating_threshold_temperature) = service.heating_threshold_temperature.as_mut() {
        heating_threshold_temperature.on_read(Some(|| {
            println!("Reading: heating_threshold_temperature");
            Ok(Some(15_f32))
        }));
        heating_threshold_temperature.set_value(15.into()).await?;
        heating_threshold_temperature.set_step_value(Some(1.into()))?;
        heating_threshold_temperature.set_min_value(Some(12.into()))?;
        heating_threshold_temperature.set_max_value(Some(25.into()))?;
    }

    let server = IpServer::new(config, storage).await?;
    server.add_accessory(heater_cooler).await?;

    let handle = server.run_handle();

    std::env::set_var("RUST_LOG", "hap=debug");
    env_logger::init();

    handle.await
}
