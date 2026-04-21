#include "DisplayController.h"

const char* shortenTopic(const char* topic);

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

    auto *sub = mqtt_subscriber.getSubscription(current_index);
    const char *topic = sub->topic;
    const char* newData = sub->bufferData;

    // shorter representation
    const char* shortTopic = shortenTopic(topic);
    char line[17]; // 16 chars + null

    // Line contains data to display
    snprintf(line, sizeof(line), "%s: %s", shortTopic, newData);

    // Update only if necessary
    if (prev_index != current_index || strcmp(line, lastDisplayedData) != 0){
        lcd.display(line);
        strncpy(lastDisplayedData, line, MQTT::Subscription::MAX_BUFLEN);
    }
}

// helper function to map topic to a shorter representation
const char* shortenTopic(const char* topic)
{
    if (strcmp(topic, "rail/tm/sn") == 0) return "tm-sn";
    if (strcmp(topic, "rail/tm/pu") == 0) return "tm-pu";
    if (strcmp(topic, "rail/tm/me") == 0) return "tm-me";
    return "unk";
}