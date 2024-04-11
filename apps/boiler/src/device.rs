use futures::stream::StreamExt;
use paho_mqtt as mqtt;
use std::sync::Arc;
use tokio::sync::Mutex;

use crate::config::Config;
use crate::state::State;

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

    pub async fn run(&mut self) -> anyhow::Result<()> {
        let mut stream = self.mqtt.get_stream(25);

        self.mqtt.connect(None).await?;

        while let Some(msg_opts) = stream.next().await {
            if let Some(msg) = msg_opts {
                self.handle_message(msg).await?;
            } else {
                log::error!("Got disconnected, trying to reconnect");

                while let Err(_err) = self.mqtt.connect(None).await {
                    tokio::time::sleep(std::time::Duration::from_secs(1)).await;
                }
            }
        }

        Ok(())
    }

    async fn handle_message(&self, msg: mqtt::Message) -> anyhow::Result<()> {
        let mut state = self.state.lock().await;

        match msg.topic() {
            // Handle device status topic
            t if t == self.config.device_status_topic => {
                *state = msg.payload().into();

                self.mqtt
                    .publish(mqtt::Message::new_retained(
                        &self.config.status_topic,
                        Vec::from(&*state),
                        0,
                    ))
                    .await?;
            }
            // Handle update request
            t if t == self.config.update_topic => {
                let state_update = msg.payload().into();

                state.update(&state_update);

                self.mqtt
                    .publish(mqtt::Message::new_retained(
                        &self.config.device_update_topic,
                        Vec::from(&state_update),
                        0,
                    ))
                    .await?;
            }
            _ => (),
        }

        Ok(())
    }
}
