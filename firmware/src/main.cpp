#include <Arduino.h>

#include "LCD.h"
#include "JoyStick.h"
#include "DisplayController.h"
#include "WifiManager.h"
#include "MQTTSubscriber.h"

// Define the LCD display
LCD lcd = LCD{0x27, 16, 2, A4, A5};

// Define the Joystick
JoyStick joyStick = JoyStick{A0, A1, 8};

// Define the Wifi-manager
Wifi::Manager wifiManager = Wifi::Manager{};
unsigned long lastWifiCheck = 0;
const unsigned long wifiCheckInterval = 2000;

// Define the MQTT Subscriber (random client id)
MQTT::Subscriber mqttSubscriber = MQTT::Subscriber("esp32_s3_mini1_emir");

// Define the DisplayController
DisplayController controller = DisplayController{lcd, joyStick, mqttSubscriber};


void setup() {
  // Setup serial
  Serial.begin(115200);
  delay(2000);

  Serial.print("Serial ready!");

  // Setup display
  lcd.begin();

  // WiFi connect (background task)
  Serial.println("Start WiFi setup");
  wifiManager.connect();
  Serial.println("WiFi connection established");

  // Setup MQTT Client
  while (!mqttSubscriber.connect()){
    delay(500);
  }

  Serial.println("MQTT connection established");
  
  mqttSubscriber.subscribe("rail/tm/sn", 1);
  mqttSubscriber.subscribe("rail/tm/pu", 1);
  mqttSubscriber.subscribe("rail/tm/me", 1);
 
}

void loop() {
  // Update the joystick & lcd
  controller.update();

  // The PubSubClient needs to update to get messages
  mqttSubscriber.getClient().loop();

  // Check for MQTT Connection loss
  if (!mqttSubscriber.getClient().connected() && wifiManager.status()) {
    Serial.println("MQTT lost connection!");
    Serial.println("Trying to reconnect MQTT...");
    bool result = mqttSubscriber.connect();
    if (result) {
      Serial.println("MQTT reconnected!");
    } else {
      Serial.println("MQTT reconnect failed.");
    }
  }

  // Check WiFi periodically
    if (millis() - lastWifiCheck >= wifiCheckInterval) {
        lastWifiCheck = millis();

        if (!wifiManager.status()) {
            Serial.println("WiFi lost, reconnecting...");
            wifiManager.tryConnect(5000);
            if (wifiManager.status()) {
                Serial.println("WiFi reconnected!");
            } else {
                Serial.println("WiFi reconnect failed.");
            }
        }
    }

    delay(10);
}
