use futures::stream::StreamExt;
use paho_mqtt as mqtt;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::Mutex;

#[derive(Debug, PartialEq, Default, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
enum Mode {
    #[default]
    Cool,
    Heat,
}

#[derive(Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct State {
    is_active: bool,
    mode: Mode,
    current_temperature: f32,
    target_temperature: f32,
}

impl From<&State> for Vec<u8> {
    fn from(input: &State) -> Self {
        serde_json::to_vec(input).unwrap()
    }
}

impl From<&[u8]> for State {
    fn from(input: &[u8]) -> Self {
        serde_json::from_slice(input).unwrap()
    }
}

pub struct Device {
    state: Arc<Mutex<State>>,
    mqtt: mqtt::AsyncClient,
}

impl Device {
    pub fn new(mqtt_broker: &str) -> Self {
        Self {
            state: Arc::new(Mutex::new(State::default())),
            mqtt: mqtt::AsyncClient::new(mqtt_broker).unwrap(),
        }
    }

    async fn handle_message(&self, msg: mqtt::Message) -> anyhow::Result<()> {
        let mut state = self.state.lock().await;
        *state = msg.payload().into();

        self.mqtt
            .publish(mqtt::Message::new_retained(
                "ac/status/LivingRoom",
                Vec::from(&*state),
                0,
            ))
            .await?;

        Ok(())
    }

    pub async fn run(&mut self) -> anyhow::Result<()> {
        let mut stream = self.mqtt.get_stream(25);

        self.mqtt.connect(None).await?;
        self.mqtt.subscribe("ac/update/LivingRoom", 0).await?;

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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn default_state_looks_as_expected() {
        let got = State::default();

        assert_eq!(
            got,
            State {
                is_active: false,
                mode: Mode::Cool,
                current_temperature: 0.0,
                target_temperature: 0.0,
            }
        )
    }
}
