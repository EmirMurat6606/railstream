#include <Arduino.h>

#include "LCD.h"
#include "JoyStick.h"
#include "DisplayController.h"
#include "WifiManager.h"
#include "MQTTClient.h"

// Define the LCD display
LCD lcd = LCD{0x27, 16, 2, A4, A5};

// Define the Joystick
JoyStick joyStick = JoyStick{A0, A1, 8};

// Define the DisplayController
DisplayController controller = DisplayController{lcd, joyStick};

// Define the Wifi-manager
Wifi::Manager wifiManager = Wifi::Manager{};
unsigned long lastWifiCheck = 0;
const unsigned long wifiCheckInterval = 2000;

// Define the MQTT Client (random client id)
MQTT::Client mqttClient = MQTT::Client("esp32_s3_mini1_emir_ncvkxmqeuifuazoeifd");

void setup() {
  // Setup serial
  Serial.begin(115200);
  delay(2000);

  Serial.print("Serial ready!");

  // Setup display
  lcd.begin();

  lcd.addBufferData("1");
  lcd.addBufferData("2");
  lcd.addBufferData("3");
  lcd.addBufferData("4");
  lcd.addBufferData("5");
  lcd.addBufferData("6");
  lcd.addBufferData("7");
  lcd.addBufferData("8");
  lcd.addBufferData("9");
  lcd.addBufferData("10");

  // WiFi connect (background task)
  Serial.println("Start WiFi setup");
  wifiManager.connect();
  Serial.println("WiFi connection established");

  // Setup MQTT Client
  while (!mqttClient.connect()){
    delay(500);
  }
  
  mqttClient.subscribe("topic/test/test1", 1);
  mqttClient.subscribe("topic/test/test2", 1);
 
}

void loop() {
  // Update the joystick & lcd
  controller.update();

  // The PubSubClient needs to update to get messages
  mqttClient.getClient().loop();

  // Check for MQTT Connection loss
  if (!mqttClient.getClient().connected() && wifiManager.status()) {
    Serial.println("MQTT lost connection!");
    mqttClient.connect();
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
