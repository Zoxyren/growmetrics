# 🌐 ESP32 IoT Simulation Platform

Dieses Projekt ist ein **Proof of Concept (PoC)** für eine hochskalierbare IoT-Infrastruktur. Es demonstriert, wie Messwerte von virtuellen ESP32-Geräten über moderne Messaging-Systeme verarbeitet, sicher gespeichert und effizient bereitgestellt werden. 

Das Herzstück ist ein **Go-Backend**, das nach den Prinzipien der **Hexagonalen Architektur (Ports & Adapters)** entwickelt wurde, um maximale Testbarkeit und technologische Unabhängigkeit zu gewährleisten.

---

## 🚀 Architektur-Ziel

Das Ziel ist ein entkoppelter und extrem performanter Datenfluss:  
**IoT-Sensor → MQTT → NATS.IO → Go-Backend (Core) → PostgreSQL & Redis → Grafana**

Durch den Einsatz von **NATS** als Message-Backbone und **Redis** als Cache ist das System darauf ausgelegt, die Datenströme von **über 1.000 Kunden** und deren Geräten mit minimaler Latenz zu verarbeiten.

---

## 🧩 Systemübersicht

```text
         ┌────────────────────────────────────────────┐
         │              Simulationsebene              │
         │ (ESP32 virtuell in Wokwi oder Go-Simulator)│
         │  → sendet Metriken (Temp, Feuchte, Druck)  │
         └─────────────────────┬──────────────────────┘
                               │ MQTT
                               ▼
                 ┌──────────────────────────┐
                 │   MQTT-Broker (HiveMQ)   │
                 └──────────────┬───────────┘
                                │ 
                                ▼
                 ┌──────────────────────────┐
                 │         NATS.IO          │
                 │   • Message Backbone     │
                 │   • Low-Latency Pub/Sub  │
                 └──────────────┬───────────┘
                                │
                                ▼
         ┌────────────────────────────────────────────┐
         │             Go Backend (Core)              │
         │  ┌──────────────────────────────────────┐  │
         │  │        Business Logic Service        │  │
         │  │ • Auto-Registration von Geräten      │  │
         │  │ • Datenanreicherung (Kunden-Mapping) │  │
         │  └──────┬───────────────────────┬───────┘  │
         └─────────┼───────────────────────┼──────────┘
                   │                       │
         ┌─────────▼────────┐    ┌──────────▼─────────┐
         │ PostgreSQL DB    │    │    Redis Cache     │
         │ • Langzeitarchiv │    │ • Echtzeit-Werte   │
         │ • Kundenverwaltung│    │ • Schneller Zugriff│
         └─────────┬────────┘    └────────────────────┘
                   │
                   ▼
         ┌──────────────────────────┐
         │   Monitoring (Grafana)   │
         │   • Historische Analysen │
         │   • Live-Dashboards      │
         └──────────────────────────┘
