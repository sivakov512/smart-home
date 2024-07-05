use embassy_executor::Spawner;
use embassy_time::{Duration, Timer};
use esp_hal::{rng::Rng, timer::systimer};
use esp_wifi::{
    wifi::{self, WifiController, WifiDevice, WifiEvent, WifiStaDevice, WifiState},
    EspWifiInitFor,
};

const SSID: &str = env!("SSID");
const PASSWORD: &str = env!("PASSWORD");

macro_rules! new_static {
    ($t:ty,$val:expr) => {{
        static STATIC_CELL: static_cell::StaticCell<$t> = static_cell::StaticCell::new();
        STATIC_CELL.uninit().write(($val))
    }};
}

pub fn setup(
    timer: systimer::Alarm<systimer::Target, esp_hal::Blocking, 0>,
    rng: Rng,
    radio_clocks: esp_hal::peripherals::RADIO_CLK,
    clocks: &esp_hal::clock::Clocks,
    wifi: esp_hal::peripherals::WIFI,
    spawner: &Spawner,
) -> &'static embassy_net::Stack<WifiDevice<'static, WifiStaDevice>> {
    let wifi_init =
        esp_wifi::initialize(EspWifiInitFor::Wifi, timer, rng, radio_clocks, &clocks).unwrap();

    let (wifi_iface, wifi_controller) =
        wifi::new_with_mode(&wifi_init, wifi, WifiStaDevice).unwrap();

    let dhcp_config = embassy_net::Config::dhcpv4(Default::default());
    let net_stack = &*new_static!(
        embassy_net::Stack<WifiDevice<'_, WifiStaDevice>>,
        embassy_net::Stack::new(
            wifi_iface,
            dhcp_config,
            new_static!(
                embassy_net::StackResources<3>,
                embassy_net::StackResources::<3>::new()
            ),
            1234,
        )
    );

    spawner.spawn(connect(wifi_controller)).unwrap();
    spawner.spawn(run(net_stack)).unwrap();

    net_stack
}

#[embassy_executor::task]
async fn run(net_stack: &'static embassy_net::Stack<WifiDevice<'static, WifiStaDevice>>) {
    net_stack.run().await;
}

#[embassy_executor::task]
async fn connect(mut controller: WifiController<'static>) {
    loop {
        match wifi::get_wifi_state() {
            WifiState::StaConnected => {
                controller.wait_for_event(WifiEvent::StaDisconnected).await;
                Timer::after(Duration::from_millis(5_000)).await;
            }
            _ => {}
        }

        if !matches!(controller.is_started(), Ok(true)) {
            let config = wifi::Configuration::Client(wifi::ClientConfiguration {
                ssid: SSID.try_into().unwrap(),
                password: PASSWORD.try_into().unwrap(),
                auth_method: wifi::AuthMethod::WPA2Personal,
                ..Default::default()
            });
            controller.set_configuration(&config).unwrap();
            controller.start().await.unwrap();
        }
        match controller.connect().await {
            Ok(_) => log::info!("Connected"),
            Err(e) => {
                log::error!("Failed to connect: {e:?}");
                Timer::after(Duration::from_millis(500)).await;
            }
        }
    }
}
