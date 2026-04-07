#ifndef JOY_STICK_H
#define JOY_STICK_H

#include <Arduino.h>

/**
 * Describes the Joystick z-axis state
 */
enum ZState
{
    PRESSED,
    NOT_PRESSED
};

struct Coordinates
{
    uint16_t x;
    uint16_t y;
    ZState button;

    String print();
};

class JoyStick
{

private:
    const uint8_t x_axis_pin;
    const uint8_t y_axis_pin;
    const uint8_t z_axis_pin;

    Coordinates coordinates;

public:
    /**
     * @brief Constructor of the JoyStick class
     *
     * Make sure the z_axis_pin on your device has the ability for pullup resistance!!!
     */
    JoyStick(const uint8_t x_axis_pin, const uint8_t y_axis_pin, const uint8_t z_axis_pin);

    /**
     * @brief Updates the x, y and z coordinates of the joystick
     *
     * @return Returns a boolean indicating if the update was succesfull or not
     */
    bool update();

    /**
     * @brief Returns the state of the joystick (e.g. the x, y and z values)
     *
     * @return The coordinates of the joystick
     */
    Coordinates getState();
};

#endif