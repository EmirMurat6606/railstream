#ifndef LCD_H
#define LCD_H

#include <LiquidCrystal_I2C.h>
#include <Wire.h>

class LCD
{

private:

    const uint8_t I2C_addr;

    const uint8_t cols;
    const uint8_t rows;

    const uint8_t SDA_PIN;
    const uint8_t SCL_PIN;

    LiquidCrystal_I2C lcd;

public:
    /**
     * @brief Construct a new LCD object and initialize the display.
     *
     * @param interface_addr The preferred I2C address of the LCD (commonly 0x27 or 0x3F)
     * @param cols Number of columns of the LCD (usually 16 or 20)
     * @param rows Number of rows of the LCD (usually 2 or 4)
     * @param sda_pin Pin number for SDA (ignored on Arduino Uno R4 and Freenove ESP32-S3)
     * @param scl_pin Pin number for SCL (ignored on Arduino Uno R4 and Freenove ESP32-S3)
     */
    LCD(uint8_t interface_addr, uint8_t cols, uint8_t rows, uint8_t sda_pin, uint8_t scl_pin);

    /**
     * @brief Does the hardware setup of the lcd display
     *
     *  This function will:
     * - Initialize I2C (Wire.begin)
     * - Turn on the backlight
     * - Print "Initialized" on the first line
     */
    void begin();

    /**
     * @brief Displays some text on the lcd
     *
     * @param data data to display
     */
    void display(const char* data);

    /**
     * @brief Clears the display
     */
    void clear();
};

#endif