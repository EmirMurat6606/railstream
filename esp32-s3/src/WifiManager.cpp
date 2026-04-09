#include "WifiManager.h"
#include "Secrets.h"

Wifi::Credentials::Credentials() : ssid_Router(WIFISecrets::SSID_ROUTER), password_Router(WIFISecrets::PASSWORD_ROUTER) {}

void Wifi::Manager::connect()
{
    while (!tryConnect(20000))
    {   
        // wait a little and try again
        delay(2000);
    }
}

bool Wifi::Manager::tryConnect(const unsigned long timeout){
     WiFi.begin(credentials.ssid_Router, credentials.password_Router);

    unsigned long startTime = millis();

    while (WiFi.status() != WL_CONNECTED)
    {
        if (millis() - startTime > timeout)
        {
            return false;
        }

        delay(500);
    }

    return true; 
}

bool Wifi::Manager::status(){
    return WiFi.status() == WL_CONNECTED;
}