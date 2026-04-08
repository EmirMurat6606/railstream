#ifndef WIFI_MANAGER_H
#define WIFI_MANAGER_H

#include <WiFiS3.h>

/**
 * @brief stores the WiFi credentials
 */
struct Credentials{
    String ssid_Router;
    String password_Router;
};

class WifiManager{

private:
    Credentials credentials;

public:
    /**
     * @brief Constructor of the WifiManager class
     */
    WifiManager();

};

#endif