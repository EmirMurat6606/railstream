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
    delay(1000);
    lcd.clear();
}

void LCD::clear()
{
    lcd.clear();
}

void LCD::display(uint8_t pos)
{
    lcd.setCursor(0, 0);

    if (pos >= this->bufferSize){
        pos = this->bufferSize - 1;
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
    if (bufferSize >= MAX_ITEMS)
        return; 

    strncpy(buffer[bufferSize], text, MAX_LENGTH - 1);
    buffer[bufferSize][MAX_LENGTH - 1] = '\0';
    bufferSize++;
}
