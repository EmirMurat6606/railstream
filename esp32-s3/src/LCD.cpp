#include "LCD.h"

bool LCD::i2CAddrTest(uint8_t addr){
    Wire.beginTransmission(addr);
    return (Wire.endTransmission() == 0);
}

LCD::LCD(u_int8_t interface_addr, u_int8_t cols, u_int8_t rows, uint8_t sda_pin, uint8_t scl_pin) : SDA_PIN(sda_pin), SCL_PIN(scl_pin), lcd(interface_addr, cols, rows)
{   
     // Start I2C 
    Wire.begin();
    
    uint8_t final_addr = interface_addr;

    if (!i2CAddrTest(interface_addr))
    {
        if (i2CAddrTest(0x27))
            final_addr = 0x27;
        else if (i2CAddrTest(0x3F))
            final_addr = 0x3F;
    }

    // Reconstruct 
    lcd = LiquidCrystal_I2C(final_addr, cols, rows);

    // Setup the lcd
    lcd.init();
    lcd.backlight();

    // Test
    lcd.print("Initialized");
}

void LCD::clear(){
    lcd.clear();
}

void LCD::display(const char * text){
    lcd.setCursor(0, 0);
    lcd.print(text);
}
