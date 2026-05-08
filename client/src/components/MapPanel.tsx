import React, { useEffect, useMemo } from 'react';
import { Circle, CircleMarker, MapContainer, Popup, TileLayer, useMap, useMapEvents } from 'react-leaflet';

interface Restaurant {
  id: number;
  name: string;
  latitude: number;
  longitude: number;
  avg_score: number;
  review_count: number;
}

interface Props {
  coords: { lat: number; lng: number } | null;
  accuracy?: number | null;
  restaurants: Restaurant[];
  pickingRestaurantName?: string | null;
  onPickCoordinates?: (lat: number, lng: number) => void;
}

const PI = Math.PI;
const A = 6378245.0;
const EE = 0.00669342162296594323;

const outOfChina = (lat: number, lng: number) =>
  lng < 72.004 || lng > 137.8347 || lat < 0.8293 || lat > 55.8271;

const transformLat = (lng: number, lat: number) => {
  let ret = -100.0 + 2.0 * lng + 3.0 * lat + 0.2 * lat * lat + 0.1 * lng * lat + 0.2 * Math.sqrt(Math.abs(lng));
  ret += (20.0 * Math.sin(6.0 * lng * PI) + 20.0 * Math.sin(2.0 * lng * PI)) * 2.0 / 3.0;
  ret += (20.0 * Math.sin(lat * PI) + 40.0 * Math.sin((lat / 3.0) * PI)) * 2.0 / 3.0;
  ret += (160.0 * Math.sin((lat / 12.0) * PI) + 320 * Math.sin((lat * PI) / 30.0)) * 2.0 / 3.0;
  return ret;
};

const transformLng = (lng: number, lat: number) => {
  let ret = 300.0 + lng + 2.0 * lat + 0.1 * lng * lng + 0.1 * lng * lat + 0.1 * Math.sqrt(Math.abs(lng));
  ret += (20.0 * Math.sin(6.0 * lng * PI) + 20.0 * Math.sin(2.0 * lng * PI)) * 2.0 / 3.0;
  ret += (20.0 * Math.sin(lng * PI) + 40.0 * Math.sin((lng / 3.0) * PI)) * 2.0 / 3.0;
  ret += (150.0 * Math.sin((lng / 12.0) * PI) + 300.0 * Math.sin((lng / 30.0) * PI)) * 2.0 / 3.0;
  return ret;
};

const wgs84ToGcj02 = (lat: number, lng: number): [number, number] => {
  if (outOfChina(lat, lng)) return [lat, lng];

  let dLat = transformLat(lng - 105.0, lat - 35.0);
  let dLng = transformLng(lng - 105.0, lat - 35.0);
  const radLat = (lat / 180.0) * PI;
  let magic = Math.sin(radLat);
  magic = 1 - EE * magic * magic;
  const sqrtMagic = Math.sqrt(magic);
  dLat = (dLat * 180.0) / (((A * (1 - EE)) / (magic * sqrtMagic)) * PI);
  dLng = (dLng * 180.0) / ((A / sqrtMagic) * Math.cos(radLat) * PI);
  return [lat + dLat, lng + dLng];
};

const gcj02ToWgs84 = (lat: number, lng: number): [number, number] => {
  if (outOfChina(lat, lng)) return [lat, lng];
  const [mgLat, mgLng] = wgs84ToGcj02(lat, lng);
  return [lat * 2 - mgLat, lng * 2 - mgLng];
};

const RecenterMap = ({ center }: { center: [number, number] }) => {
  const map = useMap();

  useEffect(() => {
    map.setView(center, map.getZoom());
  }, [center, map]);

  return null;
};

const MapClickPicker = ({ enabled, onPick }: { enabled: boolean; onPick?: (lat: number, lng: number) => void }) => {
  useMapEvents({
    click(event) {
      if (!enabled || !onPick) return;
      const [wgsLat, wgsLng] = gcj02ToWgs84(event.latlng.lat, event.latlng.lng);
      onPick(wgsLat, wgsLng);
    },
  });
  return null;
};

export const MapPanel = ({ coords, accuracy, restaurants, pickingRestaurantName, onPickCoordinates }: Props) => {
  const mappedUserCoords = useMemo<[number, number] | null>(() => {
    if (!coords) return null;
    return wgs84ToGcj02(coords.lat, coords.lng);
  }, [coords]);

  const mappedRestaurants = useMemo(
    () => restaurants.map((item) => ({ ...item, mapped: wgs84ToGcj02(item.latitude, item.longitude) })),
    [restaurants]
  );

  const center = useMemo<[number, number]>(() => {
    if (mappedUserCoords) {
      return mappedUserCoords;
    }
    if (mappedRestaurants.length > 0) {
      return mappedRestaurants[0].mapped;
    }
    return [39.9042, 116.4074];
  }, [mappedUserCoords, mappedRestaurants]);

  return (
    <div className="sticky top-24 space-y-3 rounded-2xl border border-[#e8eaed] bg-white p-4 shadow-sm">
      <div className="flex items-center justify-between">
        <h3 className="text-base font-semibold text-[#202124]">地图视图</h3>
        <span className="text-xs text-[#5f6368]">{restaurants.length} 个标记</span>
      </div>
      <div className="h-[460px] overflow-hidden rounded-xl border border-[#e8eaed]">
        <MapContainer center={center} zoom={13} scrollWheelZoom className="h-full w-full">
          <RecenterMap center={center} />
          <MapClickPicker enabled={Boolean(pickingRestaurantName)} onPick={onPickCoordinates} />
          <TileLayer
            attribution="&copy; 高德地图"
            url="https://webrd0{s}.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&x={x}&y={y}&z={z}"
            subdomains={['1', '2', '3', '4']}
          />
          {mappedUserCoords && (
            <>
              {typeof accuracy === 'number' && accuracy > 0 && (
                <Circle
                  center={mappedUserCoords}
                  radius={Math.max(accuracy, 20)}
                  pathOptions={{ color: '#1a73e8', fillColor: '#1a73e8', fillOpacity: 0.12 }}
                />
              )}
              <CircleMarker
                center={mappedUserCoords}
                radius={8}
                pathOptions={{ color: '#1a73e8', fillColor: '#1a73e8', fillOpacity: 0.9 }}
              >
                <Popup>
                  你的位置
                  {typeof accuracy === 'number' ? `（精度约 ${Math.round(accuracy)} 米）` : ''}
                </Popup>
              </CircleMarker>
            </>
          )}
          {mappedRestaurants.map((item) => (
            <CircleMarker
              key={item.id}
              center={item.mapped}
              radius={6}
              pathOptions={{ color: '#ea4335', fillColor: '#ea4335', fillOpacity: 0.85 }}
            >
              <Popup>
                <div className="text-sm">
                  <div className="font-semibold">{item.name}</div>
                  <div>
                    {item.avg_score.toFixed(1)} 分 · {item.review_count} 条评价
                  </div>
                </div>
              </Popup>
            </CircleMarker>
          ))}
        </MapContainer>
      </div>
      {pickingRestaurantName && (
        <div className="rounded-lg bg-[#e8f0fe] px-3 py-2 text-xs text-[#1a73e8]">
          正在为「{pickingRestaurantName}」选点：请在地图上单击目标位置
        </div>
      )}
      <p className="text-xs text-[#5f6368]">
        中心坐标：{center[0].toFixed(4)}, {center[1].toFixed(4)}
      </p>
    </div>
  );
};
