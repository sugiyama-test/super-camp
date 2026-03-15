"use client";

import { useEffect, useRef } from "react";
import { MapContainer, TileLayer, Marker, Popup, useMap } from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import type { Campsite } from "@/stores/useCampsiteStore";

delete (L.Icon.Default.prototype as unknown as Record<string, unknown>)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
  iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
  shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
});

type ValidCampsite = Campsite & { latitude: number; longitude: number };

interface ControllerProps {
  selectedId: number | null;
  campsites: ValidCampsite[];
  markerRefs: React.MutableRefObject<Map<number, L.Marker>>;
}

function MapController({ selectedId, campsites, markerRefs }: ControllerProps) {
  const map = useMap();

  useEffect(() => {
    if (selectedId == null) return;
    const campsite = campsites.find((c) => c.id === selectedId);
    if (!campsite) return;
    map.flyTo([campsite.latitude, campsite.longitude], 13, { duration: 0.8 });
    setTimeout(() => {
      markerRefs.current.get(selectedId)?.openPopup();
    }, 850);
  }, [selectedId]);

  return null;
}

interface Props {
  campsites: Campsite[];
  selectedId?: number | null;
}

export default function CampsiteMap({ campsites, selectedId = null }: Props) {
  const markerRefs = useRef<Map<number, L.Marker>>(new Map());

  const validCampsites = campsites.filter(
    (c): c is ValidCampsite => c.latitude != null && c.longitude != null
  );

  if (validCampsites.length === 0) return null;

  const JAPAN_CENTER: [number, number] = [36.5, 137.5];

  return (
    <MapContainer
      center={JAPAN_CENTER}
      zoom={5}
      style={{ height: "256px", borderRadius: "12px" }}
      scrollWheelZoom={false}
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      <MapController selectedId={selectedId} campsites={validCampsites} markerRefs={markerRefs} />
      {validCampsites.map((campsite) => (
        <Marker
          key={campsite.id}
          position={[campsite.latitude, campsite.longitude]}
          ref={(marker) => {
            if (marker) markerRefs.current.set(campsite.id, marker);
            else markerRefs.current.delete(campsite.id);
          }}
        >
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
