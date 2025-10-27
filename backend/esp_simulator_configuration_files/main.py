import network
import time
from umqtt.simple import MQTTClient
import ujson as json
import ubinascii
import random
import ssl

# --- 1. WLAN (Wokwi-spezifisch) ---
WIFI_SSID = "Wokwi-GUEST"
WIFI_PASS = ""

# --- 2. MQTT-Broker (HiveMQ Cloud) ---
MQTT_SERVER = "8b758ea22b9f4c0f94ac43c9b09a254f.s1.eu.hivemq.cloud"
MQTT_USER = "oliver1"
MQTT_PASS = "tESTUSER1234"
MQTT_PORT = 8883
MQTT_TOPIC = b"esp32/oliver1/metrics"

# --- 3. Eindeutige Client-ID ---
wlan = network.WLAN(network.STA_IF)
wlan.active(True)
CLIENT_ID = b"esp32-wokwi-" + ubinascii.hexlify(wlan.config("mac"))


# --- WLAN-Verbindung ---
def connect_wifi():
    print(f"Verbinde mit WLAN '{WIFI_SSID}' ...")
    wlan.connect(WIFI_SSID, WIFI_PASS)

    for _ in range(15):
        if wlan.isconnected():
            print("✅ WLAN verbunden:", wlan.ifconfig())
            return True
        time.sleep(1)
        print(".", end="")

    print("❌ WLAN-Verbindung fehlgeschlagen")
    return False


# --- MQTT-Verbindung ---
def connect_mqtt():
    print("🔗 Verbinde mit MQTT-Broker (HiveMQ Cloud über SSL)...")
    try:
        client = MQTTClient(
            client_id=CLIENT_ID,
            server=MQTT_SERVER,
            port=MQTT_PORT,
            user=MQTT_USER,
            password=MQTT_PASS,
            ssl=True,
            ssl_params={"server_hostname": MQTT_SERVER},
        )
        client.connect()
        print("✅ Erfolgreich mit HiveMQ verbunden!")
        return client
    except Exception as e:
        print("❌ MQTT-Verbindung fehlgeschlagen:", e)
        return None


# --- Zufällige Metriken senden ---
def send_random_metrics(client):
    temp = round(random.uniform(19.5, 27.5), 2)
    humidity = round(random.uniform(35.0, 65.0), 2)
    pressure = round(random.uniform(990, 1035), 1)

    data = {
        "device": CLIENT_ID.decode(),
        "temperature": temp,
        "humidity": humidity,
        "pressure": pressure,
    }

    msg = json.dumps(data)
    try:
        client.publish(MQTT_TOPIC, msg.encode())
        print(f"📤 Gesendet: {msg}")
    except Exception as e:
        print("❌ Fehler beim Senden:", e)


# --- Hauptprogramm ---
if connect_wifi():
    client = connect_mqtt()
    if client:
        try:
            while True:
                send_random_metrics(client)
                time.sleep(10)
        except KeyboardInterrupt:
            print("\n⏹️ Programm beendet.")
        finally:
            client.disconnect()
            print("MQTT getrennt.")
    else:
        print("MQTT-Verbindung nicht möglich.")
else:
    print("Keine WLAN-Verbindung. Programm abgebrochen.")
