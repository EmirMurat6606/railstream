#include "DisplayController.h"

DisplayController::DisplayController(LCD &lcd, JoyStick &joy_stick, MQTT::Subscriber &mqtt_subscriber) : lcd(lcd), joy_stick(joy_stick), mqtt_subscriber(mqtt_subscriber), current_index(0), previousX(NEUTRALX) {
     lastDisplayedData[0] = '\0'; // Initialize with empty string
}

void DisplayController::update()
{
    // Update the joystick
    this->joy_stick.update();
    const JoyStickState &state = this->joy_stick.getState();

    // Update the display
    uint8_t prev_index = this->current_index;

    switch (state.x)
    {
    case UP:
        if (this->previousX != UP)
            this->current_index = (this->current_index >= this->mqtt_subscriber.getSubscriptionCount() - 1) ? this->mqtt_subscriber.getSubscriptionCount() - 1 : this->current_index + 1;
        this->previousX = UP;
        break;
    case DOWN:
        if (this->previousX != DOWN)
            this->current_index = (this->current_index <= 0) ? 0 : this->current_index - 1;
        this->previousX = DOWN;
        break;
    default:
        this->previousX = NEUTRALX;
        break;
    }

    // Update only if necessary
    const char* newData = mqtt_subscriber.getSubscription(current_index)->bufferData;
    if (prev_index != current_index || strcmp(newData, lastDisplayedData) != 0){
        lcd.display(newData);
        strncpy(lastDisplayedData, newData, MQTT::Subscription::MAX_BUFLEN);
    }
}