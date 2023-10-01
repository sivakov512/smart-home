mod device;
use device::Device;

#[tokio::main]
async fn main() {
    Device::new("tcp://localhost:1883").run().await.unwrap();
}
