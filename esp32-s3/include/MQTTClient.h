#ifndef MQTT_CLIENT_H
#define MQTT_CLIENT_H

#include <PubSubClient.h>
#include <WiFiSSLClient.h>



namespace MQTT
{

    /**
     * @brief stores the WiFi credentials (router id and password)
     *
     * Uses the id and password from the WiFiSecrets namespace in Secrets.h
     */
    struct Credentials
    {
        const char *mqtt_url;
        const char *mqtt_username;
        const char *mqtt_password;

        Credentials();
    };

    class Client
    {

    private:
        Credentials credentials;

        WiFiSSLClient wifiClient;
        
        PubSubClient client;

        const char* client_id;

        /**
         * @brief Callback function that is executed when the client receives a message from the specified topic
         *
         * Note: Callback function needs to be static and cannot be bound
         *
         * @param topic the topic from which a message was received
         * @param payload the payload of the received message
         * @param length the length of the received message payload
         */
        static void callback(char *topic, byte *payload, unsigned int length);

    public:

        /**
         * @brief Constructor of the Client class 
         * 
         * @param client_id unique id of the client
         */

        Client(const char* client_id);

        /**
         * @brief connects to the MQTT broker
         *
         * @return Boolean indicating if the connection was successfull or not
         */
        bool connect();

        /**
         * @brief subscribes to a specific topic
         * 
         * @param topic the name of the topic
         * @param qos the Quality of Serives level (0 or 1)
         * 
         * @return Boolean indicating if the subscription was successfull or not
         */
        bool subscribe(const char* topic,  uint8_t qos);

        /**
         * @brief returns the PubSubClient  
         */ 
        PubSubClient& getClient();


    };
};

#endif