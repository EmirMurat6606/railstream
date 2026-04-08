#ifndef LCD_H
#define LCD_H

#include <LiquidCrystal_I2C.h>
#include <Wire.h>
#include <vector>

class LCD
{

private:
    const uint8_t I2C_addr;

    const uint8_t cols;
    const uint8_t rows;

    const uint8_t SDA_PIN;
    const uint8_t SCL_PIN;

    std::vector<String> buffer;

    LiquidCrystal_I2C lcd;

public:
    /**
     * @brief Construct a new LCD object and initialize the display.
     *
     * @param interface_addr The preferred I2C address of the LCD (commonly 0x27 or 0x3F)
     * @param cols Number of columns of the LCD (usually 16 or 20)
     * @param rows Number of rows of the LCD (usually 2 or 4)
     * @param sda_pin Pin number for SDA (ignored on Arduino Uno R4)
     * @param scl_pin Pin number for SCL (ignored on Arduino Uno R4)
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
     * @param pos position of the data in the buffer to display
     */
    void display(uint8_t pos);

    /**
     * @brief Clears the display
     */
    void clear();

    /**
     * @brief Returns the length of the buffer (amount of data to display)
     * 
     * @return The length of the internal data buffer
     */
    int getBufferLength() const;

    /**
     * @brief Clears the data in the buffer
     */
    void clearBuffer();

    /**
     * @brief Adds data to the data buffer
     */
    void addBufferData(String data);
};

#endif