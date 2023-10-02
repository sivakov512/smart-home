use crate::config::Config;
use crate::state::State;
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
        let mut state = self.state.lock().await;
        *state = msg.payload().into();

        self.mqtt
            .publish(mqtt::Message::new_retained(
                &self.config.status_topic,
                Vec::from(&*state),
                0,
            ))
            .await?;

        Ok(())
    }

    pub async fn run(&mut self) -> anyhow::Result<()> {
        let mut stream = self.mqtt.get_stream(25);

        self.mqtt.connect(None).await?;
        self.mqtt.subscribe(&self.config.update_topic, 0).await?;

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
