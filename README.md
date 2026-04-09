# RailStream – Smart Train Display using ESP32, MQTT & AWS

## 📌 Project Motivation

Every morning before leaving for school, the same frustration:

You already put on your gloves.  
Your jacket is zipped.  
Your backpack is on.  

And then you suddenly need to check your train.

You grab your phone. Fingerprint doesn't work because you have your gloves on. You open the NMBS app. It loads slowly. You manually enter your route.

All of this just to check if your train is delayed.

This project solves that (my) problem.

RailStream is a dedicated embedded device that continuously displays live train data for my commute. It stays up-to-date automatically and allows navigation using a joystick — no phone required.

It is designed to be:
- Lightweight
- Battery-efficient
- Real-world scalable
- IoT-architected using MQTT
- Cloud-connected via AWS

---

# 🏗️ System Architecture Overview

The system consists of three main components:

1. **Freenove ESP32-S3 Device (Client / Subscriber)**
2. **MQTT Broker (HiveMQ Cloud – Serverless)**
3. **AWS EC2 Instance (Publisher / Data Fetcher)**

Belgian Mobility API  -> AWS EC2 (Golang Publisher) -> MQTT Broker (HiveMQ) -> ESP32-S3 (Subscriber) -> 16x2 LCD Display

The ESP32 never directly queries the NMBS API. All heavy processing is handled by the EC2 instance.

---

# 🧠 Why This Architecture?

## Why NOT fetch train data directly on the ESP32?

The ESP32 is:
- Battery-powered (via powerbank)
- Limited in RAM and CPU
- Not designed for heavy JSON parsing and HTTP polling

Fetching and processing API data directly on the ESP32 would:

- Increase power consumption
- Increase network overhead
- Increase latency
- Complicate firmware
- Reduce battery life

Instead, computation and API logic are offloaded to an AWS EC2 instance.

---

## Why Use Golang on the Server?

Golang was chosen because:

- Excellent concurrency model (goroutines)
- Efficient memory usage
- Fast execution
- Lightweight

The Go application currently:

- Periodically fetches NMBS train data
- Parses the API response
- Extracts relevant information
- Publishes structured data to MQTT topics
- Maintains connection to MQTT broker

⚠️ Note: The implementation is still a work in progress and will be extended with additional filtering and suggestion logic.

---

# 📡 Why MQTT Instead of HTTP?

MQTT is ideal for IoT devices.

Compared to HTTP:

| Feature | HTTP | MQTT |
|----------|-------|-------|
| Overhead | High | Very Low |
| Persistent Connection | No | Yes |
| Packet Size | Large | Small |
| Battery Efficiency | Lower | Higher |
| Push-based | No | Yes |

Because this device runs on battery, minimizing:

- Data usage
- CPU usage
- Network overhead

is essential.

MQTT allows:

- Lightweight communication
- Small data packets
- Persistent TCP connection
- Event-driven updates

This significantly increases battery life.

---

# 🔌 MQTT Configuration

Broker: HiveMQ Cloud (Serverless – 10GB/month free tier)

The ESP32:

- Connects to WiFi
- Connects to MQTT broker
- Subscribes to multiple topics
- Receives published train data
- Displays it on the LCD

### QoS Level

The project uses:

- **QoS 1 (At least once delivery)**

Reason:
Train data must reliably arrive, but minor duplicates are acceptable.

QoS 1 guarantees message delivery without the overhead of QoS 2.

---

### Clean Session = false

The ESP32 MQTT client is configured with: cleanSession = false

This ensures:

- The broker stores subscriptions
- Messages published while ESP is offline are queued (if persistent session)
- Device can reconnect and continue receiving data

This is important for battery-powered devices that may temporarily disconnect.

---

# 🖥️ Hardware Setup

## Board

- Freenove ESP32-S3 V5

---

## LCD Display

- 16 columns
- 2 rows
- Connected via I2C

| LCD Pin | ESP32 Pin |
|----------|-----------|
| SCL | A4 |
| SDA | A5 |

Used to display:

- Train departure times
- Delay information
- Platform information
- Possible suggestions

---

## Joystick

Used to scroll through available train data.

| Function | ESP32 Pin |
|----------|-----------|
| X-Axis | A0 |
| Y-Axis | A1 |
| Z-Axis (Button) | Digital Pin 8 |

The Z-axis pin is configured as: INPUT_PULLUP

You get:
Dedicated Device → Always Updated → Scroll → Go

It combines:

- Embedded programming
- Cloud computing
- MQTT communication
- Energy efficiency
- Real-world IoT principles

And most importantly:

No more struggling with your phone while wearing gloves.

---

Author: Emir  
Project: RailStream  
Year: 2026





