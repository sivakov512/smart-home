mod config;
mod device;
mod state;
mod thermal_sensor;

use config::Config;
use device::Device;

const CONFIG_ENV_VAR: &str = "HEATERCONFIG";
const CONFIG_DEFAULT: &str = "./config.toml";

#[tokio::main]
async fn main() {
    env_logger::init();

    let config_fpath = std::env::var(CONFIG_ENV_VAR).unwrap_or(CONFIG_DEFAULT.to_owned());
    let config = Config::from_config_file(&config_fpath).unwrap();

    Device::new(config).run().await.unwrap();
}
