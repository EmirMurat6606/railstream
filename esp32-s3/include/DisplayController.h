#ifndef DISPLAY_CONTROLLER_H
#define DISPLAY_CONTROLLER_H

#include "LCD.h"
#include "JoyStick.h"

class DisplayController{

private:

    LCD &lcd;
    JoyStick &joy_stick;

    int current_index;

    XState previousX; // Used for edge detection


public:

    /**
     * @brief This is the constructor of the DisplayController class
     */
    DisplayController(LCD &lcd, JoyStick &joy_stick);

    /**
     * @brief Updates the LCD display based on the joystick positions (state)
     */
    void update();

};

#endif