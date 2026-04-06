#ifndef LCD_H
#define LCD_H

#include <LiquidCrystal_I2C.h>
#include <Wire.h>
#include <string>

class LCD {

private:

    const u_int8_t SDA_PIN;
    const u_int8_t SCL_PIN;

    
    LiquidCrystal_I2C lcd;

    /**
     * @brief Tests whether a device responds on a given I2C address.
     * 
     * @param addr The I2C address to test.
     * @return true If a device acknowledges at this address.
     * @return false Otherwise.
     */
    bool i2CAddrTest(uint8_t addr);

public:
    
    /**
     * @brief Construct a new LCD object and initialize the display.
     * 
     * The constructor will:
     * - Initialize I2C (Wire.begin)
     * - Detect the LCD's I2C address if the provided address does not respond
     * - Initialize the LiquidCrystal_I2C object
     * - Turn on the backlight
     * - Print "Initialized" on the first line
     * 
     * @param interface_addr The preferred I2C address of the LCD (commonly 0x27 or 0x3F)
     * @param cols Number of columns of the LCD (usually 16 or 20)
     * @param rows Number of rows of the LCD (usually 2 or 4)
     * @param sda_pin Pin number for SDA (ignored on Arduino Uno R4)
     * @param scl_pin Pin number for SCL (ignored on Arduino Uno R4)
     */
    LCD(u_int8_t interface_addr, u_int8_t cols, u_int8_t rows, uint8_t sda_pin, uint8_t scl_pin);

    /**
     * @brief Displays some text on the lcd
     * 
     * @param text the text to display
     */
    void display(const char * text);

    /**
     * @brief Clears the display 
     */
    void clear();

};

#endif