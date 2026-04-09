#ifndef JOY_STICK_H
#define JOY_STICK_H

#include <Arduino.h>

/**
 * @brief Describes the joysticks z-axis state
 */
enum ZState
{
    PRESSED,
    NOT_PRESSED
};

/**
 * @brief Describes the joysticks x-axis state
 * 
 * Note that the x-axis is the vertical axis in the setup logic
 */
enum XState
{   
    NEUTRALX,
    UP,
    DOWN
};

/**
 * @brief Describes the joysticks y-axis state
 * 
 * Note that the y-axis is the horizontal axis in the setup logic
 */
enum YState{
    NEUTRALY,
    LEFT,
    RIGHT
};

struct JoyStickState
{
    XState x;
    YState y;
    ZState button;

    String print();
};

class JoyStick
{

private:
    const uint8_t x_axis_pin;
    const uint8_t y_axis_pin;
    const uint8_t z_axis_pin;

    JoyStickState state;

public:
    /**
     * @brief Constructor of the JoyStick class
     *
     * Make sure the z_axis_pin on your device has the ability for pullup resistance!!!
     */
    JoyStick(const uint8_t x_axis_pin, const uint8_t y_axis_pin, const uint8_t z_axis_pin);

    /**
     * @brief Updates the x, y and z positions (states) of the joystick
     *
     * @return Returns a boolean indicating if the update was succesfull or not
     */
    bool update();

    /**
     * @brief Returns the state of the joystick (e.g. the x, y and z values)
     *
     * @return The positions (state) of the joystick
     */
    JoyStickState getState();
};

#endif