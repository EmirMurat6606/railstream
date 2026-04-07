#include <Arduino.h>

#include "LCD.h"

// Define the LCD display
LCD lcd = LCD(0x27, 16, 2, A4, A5);
int counter = 0;

// Define the Joystick


// Define the Wifi-manager

void setup() {
  lcd.begin();
}

void loop() {
  delay(2000);
  lcd.clear();

  char buffer[21];
  sprintf(buffer, "%d", counter);
  
  lcd.display(buffer);
  counter++;

}

