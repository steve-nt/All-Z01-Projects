# 🗺️ Groupie Tracker Geolocalization

A geolocation enhancement for the Groupie Tracker project that maps concert locations of selected artists or bands using geographic coordinates.  
Developed in Go using only standard libraries.

---

## 🧭 Project Overview

This project expands the **Groupie Tracker** by integrating geolocation functionality. It takes concert addresses (e.g., "Germany, Mainz") and converts them into geographic coordinates (e.g., `49.59380, 8.15052`) to place markers on a map. The goal is to visualize all concert locations for an artist/band retrieved from the existing API.

---

## 🔑 Key Features

**Geocoding Integration:**

Translates concert locations into latitude/longitude using the Google Maps API.

Handles address formatting variations (e.g., "Paris-France" → "Paris, France").

**Interactive Map:**

Displays markers for all tour locations of a selected artist.

**Seamless Integration:**

Works alongside existing Groupie Tracker features (search, filters, artist profiles).

Zero external dependencies—pure Go standard library.

**Error Resilience:**

Gracefully handles API failures, invalid locations, and edge cases.

## 🛫 Usage

1. Clone the repository:

   ```bash
   git clone https://platform.zone01.gr/git/ttarara/groupie-tracker-geolocalization

   cd groupie-tracker-geolocalization

   go run .

   ```

 Access the app at http://localhost:8080.


---

## ✍️ Authors

Theocharoula Tarara 🎵

🎵 Sofia Busho

---

## 💃 Enjoy exploring the world of music with Groupie Tracker! 🕺
