#include "MQTTSubscriber.h"
#include "Secrets.h"

MQTT::Subscriber *MQTT::Subscriber::instance = nullptr;

MQTT::Credentials::Credentials() : mqtt_url(MQTTSecrets::MQTT_URL), mqtt_username(MQTTSecrets::MQTT_USERNAME), mqtt_password(MQTTSecrets::MQTT_PASSWORD) {}

void MQTT::Subscriber::callback(char *topic, byte *payload, unsigned int length)
{
    if (instance != nullptr)
    {
        instance->handleMessage(topic, payload, length);
    }
}

void MQTT::Subscriber::handleMessage(char *topic, byte *payload, unsigned int length)
{

    for (uint8_t i = 0; i < subscription_counter; i++)
    {

        if (strcmp(topic, subscriptions[i].topic) == 0)
        {

            uint8_t copyLen = length;
            if (copyLen >= Subscription::MAX_BUFLEN)
            {
                copyLen = Subscription::MAX_BUFLEN - 1;
            }

            memcpy(subscriptions[i].bufferData, payload, copyLen);
            subscriptions[i].bufferData[copyLen] = '\0';

            return;
        }
    }

}

MQTT::Subscriber::Subscriber(const char *client_id) : client_id(client_id), wifiClient(WiFiSSLClient()), client(wifiClient)
{
    instance = this; // for the callback function
    this->client.setServer(this->credentials.mqtt_url, 8883);
    this->client.setCallback(this->callback);
}

bool MQTT::Subscriber::connect()
{
    // cleanSession (last argument) is false
    bool result = client.connect(client_id, this->credentials.mqtt_username, this->credentials.mqtt_password, 0, 0, false, 0, false);
    return result;
}

bool MQTT::Subscriber::subscribe(const char *topic, uint8_t qos)
{
    if (subscription_counter >= MAX_SUBSCRIPTIONS)
    {
        return false;
    }

    bool result = this->client.subscribe(topic, qos);

    if (result)
    {
        this->subscriptions[subscription_counter].topic = topic;
        this->subscriptions[subscription_counter].slotIndex = subscription_counter;
        this->subscription_counter++;

    }

    return result;
}

PubSubClient &MQTT::Subscriber::getClient() { return this->client; }

uint8_t MQTT::Subscriber::getSubscriptionCount() const { return subscription_counter; }

MQTT::Subscription *MQTT::Subscriber::getSubscription(uint8_t index)
{
    if (index < subscription_counter)
    {
        return &subscriptions[index];
    }
    return nullptr;
}
