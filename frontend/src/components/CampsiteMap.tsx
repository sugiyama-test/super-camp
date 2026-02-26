"use client";

import { useEffect } from "react";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import type { Campsite } from "@/stores/useCampsiteStore";

// Next.jsでleafletのデフォルトマーカーアイコンが壊れる問題を修正
delete (L.Icon.Default.prototype as unknown as Record<string, unknown>)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
  iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
  shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
});

interface Props {
  campsites: Campsite[];
}

export default function CampsiteMap({ campsites }: Props) {
  const validCampsites = campsites.filter(
    (c): c is Campsite & { latitude: number; longitude: number } =>
      c.latitude != null && c.longitude != null
  );

  if (validCampsites.length === 0) return null;

  const center: [number, number] = [
    validCampsites.reduce((sum, c) => sum + c.latitude, 0) / validCampsites.length,
    validCampsites.reduce((sum, c) => sum + c.longitude, 0) / validCampsites.length,
  ];

  return (
    <MapContainer
      center={center}
      zoom={validCampsites.length === 1 ? 12 : 7}
      style={{ height: "256px", borderRadius: "12px" }}
      scrollWheelZoom={false}
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      {validCampsites.map((campsite) => (
        <Marker key={campsite.id} position={[campsite.latitude, campsite.longitude]}>
          <Popup>
            <div>
              <p className="font-medium text-sm">{campsite.name}</p>
              {campsite.address && (
                <p className="text-xs text-gray-500 mt-0.5">{campsite.address}</p>
              )}
            </div>
          </Popup>
        </Marker>
      ))}
    </MapContainer>
  );
}
