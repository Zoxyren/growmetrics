# main.py - MicroPython MQTT Client für HiveMQ Cloud (TLS/SSL)
# Diese Datei sollte im Ordner 'src/' platziert werden.

import time
import ssl
import json
import network
import urandom # Importiert für die Generierung von Zufallszahlen
from umqtt.simple import MQTTClient

# --- KONSTANTEN ---
# Wir erwarten, dass die Konfigurationsdateien im Ordner 'config/' liegen
CONFIG_FILE = "config/config.json"
WIFI_CONFIG_FILE = "config/wifi.json" 

MQTT_PORT = 8883  # WICHTIG: TLS-Port für HiveMQ Cloud
CLIENT_ID = "ESP32_MicroPython_Client_A"
TOPIC_PUBLISH = b"wowki/esp32/daten" # MQTT-Topics müssen Bytes sein

# Skalierung für die Zufallstemperatur (20.0°C bis 30.0°C mit 1 Dezimalstelle)
TEMP_MIN_SCALED = 200 # 20.0 * 10
TEMP_RANGE_SCALED = 100 # 10.0 * 10

# Globale Konfigurationsvariablen
config = {}
wifi_config = {}

def load_config():
    """Loads MQTT and WiFi settings from respective files."""
    global config, wifi_config
    try:
        # Load MQTT configuration
        with open(CONFIG_FILE, 'r') as f:
            config = json.load(f)
        
        # Load WiFi configuration (for Wokwi simulation)
        with open(WIFI_CONFIG_FILE, 'r') as f:
            wifi_config = json.load(f)

        print("Konfiguration erfolgreich geladen.")
        
        required_keys = ["mqtt_server", "mqtt_user", "mqtt_password"]
        if not all(k in config for k in required_keys):
            print("ERROR: config/config.json is missing necessary MQTT keys.")
            return False
            
        return True
    except OSError as e:
        # Error: Configuration file not found.
        print(f"FEHLER: Eine Konfigurationsdatei nicht gefunden: {e}. Prüfen Sie die Pfade und Dateien.")
        return False
    except ValueError as e:
        print(f"FEHLER: Fehler beim Parsen einer JSON-Datei: {e}")
        return False


def connect_wifi():
    """Connects to the WiFi network (using data from config/wifi.json for Wokwi)."""
    wlan = network.WLAN(network.STA_IF)
    wlan.active(True)
    
    ssid = wifi_config.get('ssid')
    password = wifi_config.get('password')
    
    if not wlan.isconnected():
        print(f"Verbinde mit WLAN '{ssid}' (Wokwi Simulation)...")
        wlan.connect(ssid, password)
        
        timeout = 20
        while not wlan.isconnected() and timeout > 0:
            print(".", end="")
            time.sleep(1)
            timeout -= 1
        
        if wlan.isconnected():
            print("\nWiFi verbunden. IP:", wlan.ifconfig()[0])
            return True
        else:
            print("\nFEHLER: WLAN-Verbindung fehlgeschlagen. Prüfen Sie config/wifi.json.")
            return False
    return True


def connect_mqtt():
    """Establishes a secure connection to the HiveMQ Broker."""
    mqtt_server = config['mqtt_server']
    
    # --- TLS/SSL Connection Configuration ---
    # HiveMQ Cloud requires a TLS connection
    ssl_params = {
        "server_hostname": mqtt_server,
        "ssl_version": ssl.PROTOCOL_TLS_CLIENT
    }

    print(f"Verbinde mit HiveMQ Broker ({mqtt_server})...")
    
    # Initialize the MQTTClient with TLS parameters
    client = MQTTClient(
        client_id=CLIENT_ID,
        server=mqtt_server,
        port=MQTT_PORT,
        user=config['mqtt_user'],
        password=config['mqtt_password'],
        ssl=True,  # Enable SSL/TLS
        ssl_params=ssl_params
    )
    
    try:
        client.connect()
        print("MQTT-Verbindung erfolgreich hergestellt.")
        return client
    except Exception as e:
        print(f"FEHLER: Verbindung zu MQTT fehlgeschlagen: {e}")
        return None

# Main application loop
if __name__ == '__main__':
    if not load_config():
        # Exits if configuration cannot be loaded (missing files or JSON error)
        exit() 
        
    if not connect_wifi():
        # Exits if WiFi connection fails
        exit()

    mqtt_client = connect_mqtt()

    if mqtt_client:
        counter = 0
        while True:
            try:
                # 1. Generiere einen Zufallswert im skalierten Bereich (z.B. 200 bis 300)
                # urandom.getrandbits(8) gibt eine Zahl zwischen 0 und 255
                # Modulo TEMP_RANGE_SCALED + 1 (101) gibt eine Zahl zwischen 0 und 100
                random_offset = urandom.getrandbits(8) % (TEMP_RANGE_SCALED + 1)
                random_temp_int = random_offset + TEMP_MIN_SCALED
                
                # 2. Skaliere auf den Dezimalwert (z.B. 254 -> 25.4)
                temperature = random_temp_int / 10.0
                
                # 3. Erstelle das Daten-Payload als JSON
                data = json.dumps({
                    "device_id": CLIENT_ID,
                    "timestamp": time.time(),
                    "uptime_seconds": counter,
                    "temperature_c": temperature # Neu: Zufallstemperatur
                })
                
                # Publish the data
                mqtt_client.publish(TOPIC_PUBLISH, data.encode('utf-8'))
                
                print(f"Gesendet an {TOPIC_PUBLISH.decode()}: {data}")
                
                counter += 1
                time.sleep(5) # Sende alle 5 Sekunden
                
            except OSError as e:
                print(f"Verbindungsfehler: {e}. Versuche Neuverbindung...")
                mqtt_client.disconnect()
                mqtt_client = connect_mqtt()
            
            except Exception as e:
                print(f"Unerwarteter Fehler: {e}")
                time.sleep(10)




