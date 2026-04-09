#ifndef DISPLAY_CONTROLLER_H
#define DISPLAY_CONTROLLER_H

#include "LCD.h"
#include "JoyStick.h"
#include "MQTTSubscriber.h"

class DisplayController
{

private:
    LCD &lcd;

    JoyStick &joy_stick;

    MQTT::Subscriber &mqtt_subscriber;

    int current_index;

    char lastDisplayedData[MQTT::Subscription::MAX_BUFLEN];

    XState previousX; // Used for edge detection

public:
    /**
     * @brief This is the constructor of the DisplayController class
     *
     * @param lcd the lcd display to update
     * @param joy_stick the joystick component to get position data from
     * @param mqtt_subscriber the mqtt subscriber to get message data from
     */
    DisplayController(LCD &lcd, JoyStick &joy_stick, MQTT::Subscriber &mqtt_subscriber);

    /**
     * @brief Updates the LCD display based on the joystick positions (state)
     */
    void update();
};

#endif