#include <Arduino.h>

#include "LCD.h"


void setup() {
  LCD lcd = LCD(0x27, 16, 2, A4, A5);
}

void loop() {
 
}

