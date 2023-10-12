use crate::config::Config;
use crate::state::State;
use crate::thermal_sensor;
use futures::stream::StreamExt;
use paho_mqtt as mqtt;
use std::sync::Arc;
use tokio::sync::Mutex;

pub struct Device {
    state: Arc<Mutex<State>>,
    mqtt: mqtt::AsyncClient,
    config: Config,
}

impl Device {
    pub fn new(config: Config) -> Self {
        Self {
            state: Arc::new(Mutex::new(State::default())),
            mqtt: mqtt::AsyncClient::new(config.mqtt_broker.as_str()).unwrap(),
            config,
        }
    }

    async fn handle_message(&self, msg: mqtt::Message) -> anyhow::Result<()> {
        match msg.topic() {
            t if t == self.config.update_topic => {
                let mut state = self.state.lock().await;

                *state = msg.payload().into();

                // Change realy state here
                dbg!("Handle state update here", &state);

                // Notify about changed state
                self.mqtt
                    .publish(mqtt::Message::new_retained(
                        &self.config.status_topic,
                        Vec::from(&*state),
                        0,
                    ))
                    .await?;
            }
            t if t == self.config.thermal_sensor_topic => {
                let thermal_state = thermal_sensor::State::from(msg.payload());

                if let Some(temperature) = thermal_state.temperature {
                    let mut state = self.state.lock().await;

                    state.current_temperature = temperature;

                    // Notify about changed state
                    self.mqtt
                        .publish(mqtt::Message::new_retained(
                            &self.config.status_topic,
                            Vec::from(&*state),
                            0,
                        ))
                        .await?;
                }
            }
            _ => (),
        }

        Ok(())
    }

    pub async fn run(&mut self) -> anyhow::Result<()> {
        let mut stream = self.mqtt.get_stream(25);

        self.mqtt.connect(None).await?;
        self.mqtt.subscribe(&self.config.update_topic, 0).await?;
        self.mqtt
            .subscribe(&self.config.thermal_sensor_topic, 0)
            .await?;

        while let Some(msg_opts) = stream.next().await {
            if let Some(msg) = msg_opts {
                self.handle_message(msg).await?;
            } else {
                // Got disconnected, trying to reconnect
                while let Err(_err) = self.mqtt.connect(None).await {
                    tokio::time::sleep(std::time::Duration::from_secs(1)).await;
                }
            }
        }

        Ok(())
    }
}
