"use client";

import { useEffect, useRef, useState } from "react";
import dynamic from "next/dynamic";
import { PageHeader } from "@/components/PageHeader";
import { useCampsiteStore } from "@/stores/useCampsiteStore";
import type { CreateCampsiteData, Campsite } from "@/stores/useCampsiteStore";

const CampsiteMap = dynamic(() => import("@/components/CampsiteMap"), {
  ssr: false,
  loading: () => (
    <div className="h-64 rounded-xl bg-gray-100 flex items-center justify-center text-sm text-gray-400">
      地図を読み込み中...
    </div>
  ),
});

const emptyForm: CreateCampsiteData = {
  name: "",
  address: "",
  latitude: null,
  longitude: null,
  notes: "",
};

export default function CampsitesPage() {
  const { campsites, loading, fetchCampsites, createCampsite, updateCampsite, deleteCampsite } =
    useCampsiteStore();
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState<CreateCampsiteData>(emptyForm);
  const [searchQuery, setSearchQuery] = useState("");
  const [gpsLoading, setGpsLoading] = useState(false);
  const [gpsError, setGpsError] = useState("");
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const mapRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    fetchCampsites();
  }, [fetchCampsites]);

  const resetForm = () => {
    setForm(emptyForm);
    setEditingId(null);
    setShowForm(false);
    setGpsError("");
  };

  const geocodeAddress = async (address: string): Promise<{ lat: number; lng: number } | null> => {
    try {
      const res = await fetch(
        `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(address)}&format=json&limit=1`,
        { headers: { "Accept-Language": "ja" } }
      );
      const data = await res.json();
      if (data.length > 0) {
        return { lat: parseFloat(data[0].lat), lng: parseFloat(data[0].lon) };
      }
    } catch {}
    return null;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) return;

    let submitForm = { ...form };
    if (submitForm.latitude == null && submitForm.longitude == null && submitForm.address.trim()) {
      const coords = await geocodeAddress(submitForm.address);
      if (coords) {
        submitForm = { ...submitForm, latitude: coords.lat, longitude: coords.lng };
      }
    }

    if (editingId) {
      await updateCampsite(editingId, submitForm);
    } else {
      await createCampsite(submitForm);
    }
    resetForm();
  };

  const handleEdit = (campsite: Campsite) => {
    setForm({
      name: campsite.name,
      address: campsite.address,
      latitude: campsite.latitude,
      longitude: campsite.longitude,
      notes: campsite.notes,
    });
    setEditingId(campsite.id);
    setShowForm(true);
    setGpsError("");
  };

  const handleDelete = (id: number) => {
    if (confirm("このキャンプ場を削除しますか？")) {
      deleteCampsite(id);
    }
  };

  const handleGetGPS = () => {
    if (!navigator.geolocation) {
      setGpsError("このブラウザは位置情報に対応していません");
      return;
    }
    setGpsLoading(true);
    setGpsError("");
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setForm((f) => ({
          ...f,
          latitude: pos.coords.latitude,
          longitude: pos.coords.longitude,
        }));
        setGpsLoading(false);
      },
      () => {
        setGpsError("位置情報の取得に失敗しました");
        setGpsLoading(false);
      }
    );
  };

  const buildMapsUrl = (campsite: Campsite) => {
    if (campsite.latitude != null && campsite.longitude != null) {
      return `https://www.google.com/maps?q=${campsite.latitude},${campsite.longitude}`;
    }
    if (campsite.address) {
      return `https://www.google.com/maps/search/${encodeURIComponent(campsite.address)}`;
    }
    return null;
  };

  const filteredCampsites = campsites.filter((c) => {
    if (!searchQuery.trim()) return true;
    const q = searchQuery.toLowerCase();
    return (
      c.name.toLowerCase().includes(q) ||
      (c.address ?? "").toLowerCase().includes(q)
    );
  });

  const mappableCampsites = campsites.filter(
    (c) => c.latitude != null && c.longitude != null
  );

  return (
    <div>
      <PageHeader
        title="キャンプ場"
        description="お気に入りのキャンプ場を管理しましょう"
      />

      {mappableCampsites.length > 0 && (
        <div className="mt-4" ref={mapRef}>
          <CampsiteMap campsites={mappableCampsites} selectedId={selectedId} />
        </div>
      )}

      <div className="mt-4">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="キャンプ場名・住所で検索..."
          className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
        />
      </div>

      <div className="mt-4">
        {!showForm ? (
          <button
            onClick={() => { setForm(emptyForm); setEditingId(null); setShowForm(true); }}
            className="w-full rounded-lg bg-[var(--camp-orange)] px-4 py-2 text-sm font-medium text-white hover:opacity-90"
          >
            + 新しいキャンプ場
          </button>
        ) : (
          <form onSubmit={handleSubmit} className="rounded-xl bg-white p-4 shadow-sm space-y-3">
            <div>
              <label className="block text-xs text-gray-500 mb-1">名前 *</label>
              <input
                type="text"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="キャンプ場名..."
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                required
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">住所</label>
              <input
                type="text"
                value={form.address}
                onChange={(e) => setForm({ ...form, address: e.target.value })}
                placeholder="住所を入力..."
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
              />
            </div>
            <div>
              <div className="flex items-center justify-between mb-1">
                <label className="block text-xs text-gray-500">緯度・経度</label>
                <button
                  type="button"
                  onClick={handleGetGPS}
                  disabled={gpsLoading}
                  className="flex items-center gap-1 text-xs text-[var(--camp-orange)] hover:opacity-80 disabled:opacity-50"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <circle cx="12" cy="12" r="3" /><line x1="12" y1="2" x2="12" y2="6" /><line x1="12" y1="18" x2="12" y2="22" /><line x1="2" y1="12" x2="6" y2="12" /><line x1="18" y1="12" x2="22" y2="12" />
                  </svg>
                  {gpsLoading ? "取得中..." : "現在地を取得"}
                </button>
              </div>
              {gpsError && (
                <p className="text-xs text-red-500 mb-1">{gpsError}</p>
              )}
              <div className="grid grid-cols-2 gap-3">
                <input
                  type="number"
                  value={form.latitude ?? ""}
                  onChange={(e) => setForm({ ...form, latitude: e.target.value ? Number(e.target.value) : null })}
                  placeholder="35.6762"
                  step="any"
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                />
                <input
                  type="number"
                  value={form.longitude ?? ""}
                  onChange={(e) => setForm({ ...form, longitude: e.target.value ? Number(e.target.value) : null })}
                  placeholder="139.6503"
                  step="any"
                  className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
                />
              </div>
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">メモ</label>
              <textarea
                value={form.notes}
                onChange={(e) => setForm({ ...form, notes: e.target.value })}
                placeholder="キャンプ場の特徴、設備など..."
                rows={2}
                className="w-full rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--camp-orange)]"
              />
            </div>
            <div className="flex gap-2">
              <button
                type="submit"
                className="flex-1 rounded-lg bg-[var(--camp-orange)] px-4 py-2 text-sm font-medium text-white hover:opacity-90"
              >
                {editingId ? "更新" : "登録する"}
              </button>
              <button
                type="button"
                onClick={resetForm}
                className="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-600 hover:bg-gray-50"
              >
                キャンセル
              </button>
            </div>
          </form>
        )}
      </div>

      <div className="mt-6 space-y-3">
        {loading && campsites.length === 0 && (
          <p className="text-center text-gray-400 text-sm">読み込み中...</p>
        )}
        {!loading && campsites.length === 0 && (
          <p className="text-center text-gray-400 text-sm mt-8">
            キャンプ場が登録されていません。上のボタンから追加しましょう！
          </p>
        )}
        {!loading && campsites.length > 0 && filteredCampsites.length === 0 && (
          <p className="text-center text-gray-400 text-sm mt-8">
            「{searchQuery}」に一致するキャンプ場が見つかりません。
          </p>
        )}
        {filteredCampsites.map((campsite) => {
          const mapsUrl = buildMapsUrl(campsite);
          const hasCords = campsite.latitude != null && campsite.longitude != null;
          const isSelected = selectedId === campsite.id;
          return (
            <div
              key={campsite.id}
              className={`rounded-xl bg-white p-4 shadow-sm transition-all ${hasCords ? "cursor-pointer" : ""} ${isSelected ? "ring-2 ring-[var(--camp-orange)]" : ""}`}
              onClick={() => {
                if (!hasCords) return;
                setSelectedId(campsite.id);
                mapRef.current?.scrollIntoView({ behavior: "smooth", block: "center" });
              }}
            >
              <div className="flex items-start justify-between">
                <div className="flex-1 min-w-0">
                  <span className="text-sm font-medium text-gray-800">
                    {campsite.name}
                  </span>
                  {campsite.address && (
                    <p className="mt-0.5 text-xs text-gray-500">{campsite.address}</p>
                  )}
                  <div className="mt-1 flex flex-wrap gap-2">
                    {campsite.latitude != null && campsite.longitude != null && (
                      <span className="inline-block rounded-full bg-blue-100 text-blue-600 px-2 py-0.5 text-xs">
                        {campsite.latitude.toFixed(4)}, {campsite.longitude.toFixed(4)}
                      </span>
                    )}
                    {mapsUrl && (
                      <a
                        href={mapsUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1 rounded-full bg-green-100 text-green-700 px-2 py-0.5 text-xs hover:bg-green-200"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z" /><circle cx="12" cy="10" r="3" />
                        </svg>
                        Google Maps
                      </a>
                    )}
                  </div>
                  {campsite.notes && (
                    <p className="mt-2 text-xs text-gray-600 whitespace-pre-wrap">{campsite.notes}</p>
                  )}
                </div>
                <div className="flex gap-1 ml-2">
                  <button
                    onClick={() => handleEdit(campsite)}
                    className="text-gray-400 hover:text-gray-600 p-1"
                    aria-label="編集"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z" />
                    </svg>
                  </button>
                  <button
                    onClick={() => handleDelete(campsite.id)}
                    className="text-gray-400 hover:text-red-500 p-1"
                    aria-label="削除"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M3 6h18" /><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" /><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
