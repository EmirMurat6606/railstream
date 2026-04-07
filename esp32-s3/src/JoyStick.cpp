#include "JoyStick.h"

String Coordinates::print()
{
    String output = "";

    output += "x-coordinate: " + String(x) + ", ";
    output += "y-coordinate: " + String(y) + ", ";
    output += "button: ";
    output += (button == PRESSED) ? "PRESSED" : "NOT PRESSED";

    return output;
}

JoyStick::JoyStick(const uint8_t x_axis_pin, const uint8_t y_axis_pin, const uint8_t z_axis_pin) : x_axis_pin(x_axis_pin), y_axis_pin(y_axis_pin), z_axis_pin(z_axis_pin)
{
    pinMode(z_axis_pin, INPUT_PULLUP);
    this->coordinates.button = NOT_PRESSED;
}

bool JoyStick::update()
{

    this->coordinates.x = analogRead(x_axis_pin);
    this->coordinates.y = analogRead(y_axis_pin);

    uint8_t zVal = digitalRead(z_axis_pin);
    this->coordinates.button = (zVal == LOW) ? PRESSED : NOT_PRESSED;
}

Coordinates JoyStick::getState()
{
    return coordinates;
}