#include "LCD.h"


LCD::LCD(const u_int8_t interface_addr, const u_int8_t cols, const u_int8_t rows, const uint8_t sda_pin, const uint8_t scl_pin) : I2C_addr(interface_addr), cols(cols), rows(rows), SDA_PIN(sda_pin), SCL_PIN(scl_pin), lcd(I2C_addr, cols, rows) {}


void LCD::begin()
{
    // Start I2C
    Wire.begin();

    // Setup the lcd
    lcd.init();
    lcd.backlight();

    // Test
    lcd.print("Initialized");
}

void LCD::clear()
{
    lcd.clear();
}

void LCD::display(const char *text)
{
    lcd.setCursor(0, 0);
    lcd.print(text);
}
