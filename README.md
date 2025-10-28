# 🌐 ESP32 IoT Simulation Platform 

Dieses Projekt ist ein **Proof of Concept (PoC)** für eine skalierbare IoT-Architektur,  
bei der virtuelle ESP32-Geräte Messwerte (z. B. Temperatur, Luftfeuchtigkeit, Luftdruck)  
an einen MQTT-Broker senden.  
Ein **Go-Backend** empfängt diese Daten, leitet sie über Kafka weiter, speichert sie in einer Datenbank  
und stellt sie später über eine Weboberfläche und ein Monitoring-Dashboard bereit.

---

## 🚀 Ziel des Projekts

Das Ziel ist es, den kompletten Datenfluss von  
**IoT-Sensor → MQTT → Kafka → Go-Backend → PostgreSQL → Web/Grafana**  
zu demonstrieren und eine Architektur zu entwerfen, die auf **über 1 000 Kunden** skalierbar bleibt.

---

## 🧩 Systemübersicht

```text
         ┌────────────────────────────────────────────┐
         │                Simulationsebene            │
         │ (ESP32 virtuell in Wokwi oder Go-Simulator)│
         │  → sendet Metriken (Temp, Feuchte, Druck)  │
         └─────────────────────┬──────────────────────┘
                               │ MQTT (Topic: iot/metrics)
                               ▼
                 ┌──────────────────────────┐
                 │   MQTT-Broker (HiveMQ)    │
                 │   • empfängt Sensordaten  │
                 │   • verteilt an Kafka     │
                 └──────────────┬────────────┘
                                │ MQTT → Kafka Bridge
                                ▼
                ┌──────────────────────────────┐
                │   Apache Kafka                │
                │   • Topic: iot-metrics        │
                │   • Topic: iot-control        │
                │   • Message Backbone für      │
                │     Skalierung & Services     │
                └──────────────┬───────────────┘
                               │
               ┌───────────────┼────────────────────┐
               │               │                    │
               ▼               ▼                    ▼
     ┌────────────────┐  ┌────────────────┐  ┌────────────────┐
     │ Go Backend API │  │ Control Service│  │  Redis Cache    │
     │  (Consumer)    │  │  (Producer)    │  │  • Schneller     │
     │  • verarbeitet │  │  • sendet      │  │    Zugriff auf   │
     │    Metriken    │  │    Befehle     │  │    letzte Werte  │
     │  • speichert in│  │    an MQTT     │  └────────────────┘
     │    PostgreSQL  │  └────────────────┘
     └──────┬─────────┘
            │ REST / WebSocket API
            ▼
      ┌─────────────────────────────┐
      │ Web Dashboard (Frontend)    │
      │ • Visualisiert Live-Daten   │
      │ • Sendet Steuerkommandos    │
      └──────────────┬──────────────┘
                     ▼
        ┌──────────────────────────┐
        │ Monitoring (Grafana)     │
        │ • greift auf PostgreSQL  │
        │   oder Prometheus zu     │
        │ • zeigt historische &    │
        │   Live-Metriken          │
        └──────────────────────────┘
