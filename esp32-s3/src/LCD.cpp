#include "LCD.h"

LCD::LCD(const uint8_t interface_addr, const uint8_t cols, const uint8_t rows, const uint8_t sda_pin, const uint8_t scl_pin) : I2C_addr(interface_addr), cols(cols), rows(rows), SDA_PIN(sda_pin), SCL_PIN(scl_pin), lcd(I2C_addr, cols, rows) {}

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

void LCD::display(uint8_t pos)
{
    lcd.setCursor(0, 0);

    uint8_t size = this->bufferSize;

    if (pos >= size){
        pos = size;
    }
    lcd.print(buffer[pos]);
}

int LCD::getBufferLength() const
{
    return this->bufferSize;
}

void LCD::clearBuffer()
{
    this->bufferSize = 0;
}

void LCD::addBufferData(const char* text)
{
    strncpy(buffer[bufferSize], text, MAX_LENGTH - 1);
    buffer[bufferSize][MAX_LENGTH - 1] = '\0';
    bufferSize++;
}
