#include "MQTTClient.h"
#include "Secrets.h"


MQTT::Credentials::Credentials() : mqtt_url(MQTTSecrets::MQTT_URL), mqtt_username(MQTTSecrets::MQTT_USERNAME), mqtt_password(MQTTSecrets::MQTT_PASSWORD) {}

void MQTT::Client::callback(char *topic, byte *payload, unsigned int length)
{   
    Serial.println("Topic ");
    Serial.print(topic);
    Serial.print(" Message: ");


    for (int i = 0; i < length; i++)
    {   
        Serial.print((char)payload[i]);
    }
}

MQTT::Client::Client(const char *client_id) : client_id(client_id), wifiClient(WiFiSSLClient()), client(wifiClient) 
{   
    
    this->client.setServer(this->credentials.mqtt_url, 8883);
    this->client.setCallback(this->callback);

}

bool MQTT::Client::connect()
{
    // cleanSession (last argument) is false
    Serial.println("Trying MQTT connect...");
    bool result = client.connect(client_id, this->credentials.mqtt_username, this->credentials.mqtt_password);
    if (result) Serial.println("MQTT connected!");
    else Serial.println("MQTT connection failed!");
    return result;
}

bool MQTT::Client::subscribe(const char *topic, uint8_t qos)
{
    return this->client.subscribe(topic, qos);
}

PubSubClient &MQTT::Client::getClient()
{
    return this->client;
}
