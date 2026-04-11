#ifndef WIFI_MANAGER_H
#define WIFI_MANAGER_H

#include <WiFiS3.h>

namespace Wifi
{

    /**
     * @brief stores the WiFi credentials (router id and password)
     *
     * Uses the id and password from the WiFiSecrets namespace in Secrets.h
     */
    struct Credentials
    {
        const char *ssid_Router;
        const char *password_Router;

        Credentials();
    };

    class Manager
    {

    private:
        Credentials credentials;

    public:
        /**
         * @brief Establishes a wifi connection with your local router
         *
         * This function only terminates if the WiFi connection is established
         */
        void connect();

        /**
         * @brief Tries to connect to WiFi within a specified time interval
         *
         * Note: the connection either fails or succeeds, the function ends in both cases
         *
         * @returns Boolean indicating if the connection was established or not
         */
        bool tryConnect(const unsigned long timeout);

        /**
         * @brief Returns the connection status
         *
         * @return A boolean indicating if the device is connected to WiFi or not
         */
        bool status();
    };
};

#endif