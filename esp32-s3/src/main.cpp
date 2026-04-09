#include <Arduino.h>

#include "LCD.h"
#include "JoyStick.h"
#include "DisplayController.h"

// Define the LCD display
LCD lcd = LCD{0x27, 16, 2, A4, A5};

// Define the Joystick
JoyStick joyStick = JoyStick{A0, A1, 8};

// Define the DisplayController
DisplayController controller = DisplayController{lcd, joyStick};

// Define the Wifi-manager

void setup() {
  Serial.begin(11500);
  Serial.print("Hello");

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
}

void loop() {
  controller.update();
}

