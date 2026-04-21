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
  delay(2000);

  // Setup display
  lcd.begin();

  // WiFi connect (background task)
  wifiManager.connect();

  // Setup MQTT Client
  while (!mqttSubscriber.connect()){
    delay(500);
  }

  
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
    bool result = mqttSubscriber.connect();
  }

  // Check WiFi periodically
    if (millis() - lastWifiCheck >= wifiCheckInterval) {
        lastWifiCheck = millis();

        if (!wifiManager.status()) {
            wifiManager.tryConnect(5000);
        }
    }

    delay(10);
}
