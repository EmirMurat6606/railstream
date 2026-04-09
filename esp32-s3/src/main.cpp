#include <Arduino.h>

#include "LCD.h"
#include "JoyStick.h"
#include "DisplayController.h"

// Define the LCD display
LCD lcd = LCD{0x27, 16, 2, A4, A5};

// Define the Joystick
JoyStick joyStick = JoyStick{A0, A1, 8};

// Define the Wifi-manager

void setup() {
  Serial.begin(11500);
  Serial.print("Hello");
}

void loop() {
  delay(1000);
  joyStick.update();
  String coordinates = joyStick.getState().print();
  Serial.print(coordinates);
  Serial.println("");
}

