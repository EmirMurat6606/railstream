#ifndef MQTT_SUBSCRIBER_H
#define MQTT_SUBSCRIBER_H

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

    /**
     * @brief represents a MQTT subscription, which couples a topic to an index and stores the latest message payload
     */
    struct Subscription {

        static const uint8_t MAX_BUFLEN = 17;

        const char* topic;
        char bufferData[MAX_BUFLEN];
        
        uint8_t slotIndex;
    };

    class Subscriber
    {

    public:
    
        static const uint8_t MAX_SUBSCRIPTIONS = 3;

    private:
        // needed for static callback 
        // Note: this design is not ideal, but PubSubClient does not allow binding non-static functions as callbacks
        static Subscriber* instance;


        Credentials credentials;

        WiFiSSLClient wifiClient;
        
        PubSubClient client;

        Subscription subscriptions[MAX_SUBSCRIPTIONS];

        uint8_t subscription_counter = 0;

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

        /**
         * @brief Handles the received message and stores it in the right subscription buffer
         * 
         * @param topic the topic from which a message was received
         * @param payload the payload of the received message
         * @param length the length of the received message payload
         */
        void handleMessage(char* topic, byte* payload, unsigned int length);

    public:

        /**
         * @brief Constructor of the Subscriber class 
         * 
         * @param client_id unique id of the subscriber
         */

        Subscriber(const char* client_id);

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
         * Note: if the amount of subscriptions exceeds MAX_SUBSCRIPTIONS, the operation fails automatically
         * 
         * @return Boolean indicating if the subscription was successfull or not
         */
        bool subscribe(const char* topic,  uint8_t qos);

        /**
         * @brief returns the PubSubClient  
         */ 
        PubSubClient& getClient();

        /**
         * @brief returns the number of current subscriptions
          * 
          * @return the number of current subscriptions
         */
        uint8_t getSubscriptionCount() const;

         /**
          * @brief returns the subscription at the given index
          * 
          * @param index the index of the subscription to return
          * @return Subscription* pointer to the subscription, or nullptr if index is out of bounds
          */
         Subscription* getSubscription(uint8_t index);
    };
};

#endif