mod ac;
mod characteristics;
mod homekit;

use crate::ac::AC;
use crate::characteristics::Characteristics;
use crate::homekit::build_accessory;
use hap::{
    accessory::AccessoryCategory,
    server::{IpServer, Server},
    storage::{FileStorage, Storage},
    Config, Result,
};

const AC_NAME: &str = "Midea AC";

#[tokio::main]
async fn main() -> Result<()> {
    let characteristics = Characteristics::from_config_file("./config.toml").unwrap();

    let ac = AC::from(characteristics);

    let accessory = build_accessory(&ac);

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

    let server = IpServer::new(config, storage).await?;
    server.add_accessory(accessory).await?;

    let handle = server.run_handle();

    env_logger::init();

    handle.await
}
