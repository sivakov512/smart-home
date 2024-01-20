use crate::config::Config;
use crate::state::{Mode, State};
use crate::thermal_sensor;
use futures::stream::StreamExt;
use paho_mqtt as mqtt;
use std::sync::Arc;
use tokio::sync::Mutex;

pub struct Device {
    state: Arc<Mutex<State>>,
    mqtt: mqtt::AsyncClient,
    relay: sonoff_minir3::Client,
    config: Config,
}

impl Device {
    pub fn new(config: Config) -> Self {
        Self {
            state: Arc::new(Mutex::new(State::default())),
            mqtt: mqtt::AsyncClient::new(config.mqtt_broker.as_str()).unwrap(),
            relay: sonoff_minir3::Client::new(&config.sonoff_relay_ip, 8081),
            config,
        }
    }

    async fn handle_message(&self, msg: mqtt::Message) -> anyhow::Result<()> {
        let mut state = self.state.lock().await;

        match msg.topic() {
            // Handle status topic to initialize and restore previous state
            t if t == self.config.status_topic => {
                *state = msg.payload().into();
                self.sync_state_and_relay(&mut state).await?;

                self.mqtt.unsubscribe(&self.config.status_topic).await?;
                self.mqtt
                    .subscribe_many(
                        &[&self.config.update_topic, &self.config.thermal_sensor_topic],
                        &[0, 0],
                    )
                    .await?;

                self.publish_state(&state).await?;
            }
            // Handle update request
            t if t == self.config.update_topic => {
                state.update(msg.payload().into());

                self.sync_state_and_relay(&mut state).await?;

                self.publish_state(&state).await?;
            }
            // Handle temperature changes
            t if t == self.config.thermal_sensor_topic => {
                let thermal_state = thermal_sensor::State::from(msg.payload());

                if let Some(temperature) = thermal_state.temperature {
                    state.current_temperature = temperature;

                    self.sync_state_and_relay(&mut state).await?;

                    self.publish_state(&state).await?;
                }
            }
            _ => (),
        }

        Ok(())
    }

    async fn sync_state_and_relay(&self, state: &mut State) -> anyhow::Result<()> {
        use sonoff_minir3::SwitchPosition::*;

        if state.current_temperature >= state.target_temperature {
            state.mode = Mode::Idle;
        } else {
            state.mode = Mode::Heat;
        }

        let mut switch_position = Off;
        if state.is_active && state.mode == Mode::Heat {
            switch_position = On;
        }

        self.relay.set_switch_position(switch_position).await?;

        Ok(())
    }

    pub async fn run(&mut self) -> anyhow::Result<()> {
        let mut stream = self.mqtt.get_stream(25);

        self.mqtt.connect(None).await?;

        // State initialization
        self.mqtt.subscribe(&self.config.status_topic, 0).await?;
        tokio::time::sleep(std::time::Duration::from_millis(100)).await;
        match stream.try_recv() {
            // Restore previous state if exists
            Ok(msg) => {
                if let Some(msg) = msg {
                    self.handle_message(msg).await?;
                }
            }
            // Publish current state for initialization
            Err(_) => {
                let state = self.state.lock().await;
                self.publish_state(&state).await?;
            }
        }

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

    async fn publish_state(&self, state: &State) -> anyhow::Result<()> {
        self.mqtt
            .publish(mqtt::Message::new_retained(
                &self.config.status_topic,
                Vec::from(state),
                0,
            ))
            .await?;
        Ok(())
    }
}
