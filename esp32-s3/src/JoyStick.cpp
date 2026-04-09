#include "JoyStick.h"

String JoyStickState::print()
{
    String output = "";

    output += "x-direction: " + String(x) + ", ";
    output += "y-direction: " + String(y) + ", ";
    output += "button: ";
    output += (button == PRESSED) ? "PRESSED" : "NOT PRESSED";

    return output;
}

JoyStick::JoyStick(const uint8_t x_axis_pin, const uint8_t y_axis_pin, const uint8_t z_axis_pin) : x_axis_pin(x_axis_pin), y_axis_pin(y_axis_pin), z_axis_pin(z_axis_pin)
{
    pinMode(z_axis_pin, INPUT_PULLUP);
    this->state.button = NOT_PRESSED;
}

bool JoyStick::update()
{

    uint16_t x_val = analogRead(x_axis_pin);
    uint16_t y_val = analogRead(y_axis_pin);

    if (x_val <= 300)
        this->state.x = UP;
    else if (x_val >= 800)
        this->state.x = DOWN;
    else
        this->state.x = NEUTRALX;

    if (y_val <= 300)
        this->state.y = RIGHT;
    else if (y_val >= 800)
        this->state.y = LEFT;
    else
        this->state.y = NEUTRALY;

    // If statements here to check values (test this out)
    uint8_t zVal = digitalRead(z_axis_pin);
    this->state.button = (zVal == LOW) ? PRESSED : NOT_PRESSED;
}

JoyStickState JoyStick::getState()
{
    return this->state;
}