// src/pages/Home.tsx
import { useEffect, useState, useRef } from 'react';
import { useLocation } from '../hooks/useLocation.js';
import { restaurantApi } from '../api/index.js';
import { RestaurantCard } from '../components/RestaurantCard.jsx';
import { RestaurantForm } from '../components/RestaurantForm.js';
import { MapPanel } from '../components/MapPanel';
import { MapContainer, TileLayer, useMapEvents } from 'react-leaflet';
import React from 'react';

interface Restaurant {
  id: number;
  name: string;
  latitude: number;
  longitude: number;
  avg_score: number;
  review_count: number;
  distance?: number;
  final_score?: number;
}

interface MapClickHandlerProps {
  onMapClick: (lat: number, lng: number) => void;
}

const MapClickHandler = ({ onMapClick }: MapClickHandlerProps) => {
  useMapEvents({
    click(event) {
      onMapClick(event.latlng.lat, event.latlng.lng);
    },
  });
  return null;
};

type SortOption = 'distance' | 'score' | 'reviews' | 'recommended';

const Home = () => {
  const { coords, accuracy, loading: locationLoading, error: locationError, requestLocation } = useLocation();
  const [restaurants, setRestaurants] = useState<Restaurant[]>([]);
  const [recommended, setRecommended] = useState<Restaurant[]>([]);
  const [showRestaurantForm, setShowRestaurantForm] = useState(false);
  const [isSelectingLocationForNew, setIsSelectingLocationForNew] = useState(false);
  const [newRestaurantLocation, setNewRestaurantLocation] = useState<{ lat: number; lng: number } | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState('');
  const [sortBy, setSortBy] = useState<SortOption>('distance');
  const [manualCoordsById, setManualCoordsById] = useState<Record<number, { lat: number; lng: number }>>({});
  const [pickingRestaurantId, setPickingRestaurantId] = useState<number | null>(null);
  const searchTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    requestLocation();
  }, [requestLocation]);

  const fetchRestaurants = async (searchTerm: string, sort: SortOption) => {
    try {
      setLoading(true);
      setError(null);
      const lat = coords?.lat;
      const lng = coords?.lng;
      
      const res = await restaurantApi.getAll(lat, lng, searchTerm, sort);
      setRestaurants(res.data || []); 
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取餐厅列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchRecommended = async () => {
    try {
      const lat = coords?.lat;
      const lng = coords?.lng;
      
      const res = await restaurantApi.getAll(lat, lng, '', 'recommended');
      setRecommended(res.data?.slice(0, 5) || []); 
    } catch (err) {
      console.error('获取推荐失败:', err);
    }
  };

  useEffect(() => {
    if (searchTimeoutRef.current) {
      clearTimeout(searchTimeoutRef.current);
    }
    searchTimeoutRef.current = setTimeout(() => {
      fetchRestaurants(search, sortBy);
    }, 3000);
    return () => {
      if (searchTimeoutRef.current) {
        clearTimeout(searchTimeoutRef.current);
      }
    };
  }, [search, sortBy, coords]);

  useEffect(() => {
    fetchRecommended();
  }, [coords]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchRestaurants(search, sortBy);
  };

  const handleSortChange = (newSort: SortOption) => {
    setSortBy(newSort);
  };

  const isLoading = loading || locationLoading;

  if (isLoading && restaurants.length === 0) {
    return (
      <div className="grid gap-4 md:grid-cols-2">
        {[1, 2, 3, 4].map(i => (
          <div key={i} className="h-32 bg-gray-200 rounded-xl animate-pulse"></div>
        ))}
      </div>
    );
  }

  if (locationError && restaurants.length === 0) {
    return (
      <div className="text-center py-10">
        <div className="text-yellow-600 mb-2">⚠️ {locationError}</div>
        <div className="text-gray-500 text-sm">将使用默认位置展示餐厅</div>
      </div>
    );
  }

  if (error && restaurants.length === 0) {
    return (
      <div className="text-center py-10">
        <div className="text-red-500 mb-4">{error}</div>
      </div>
    );
  }

  const displayRestaurants = restaurants.map((item) => {
    const manual = manualCoordsById[item.id];
    if (!manual) return item;
    return { ...item, latitude: manual.lat, longitude: manual.lng };
  });

  const displayRecommended = recommended.map((item) => {
    const manual = manualCoordsById[item.id];
    if (!manual) return item;
    return { ...item, latitude: manual.lat, longitude: manual.lng };
  });

  const handleManualLocate = (restaurantId: number) => {
    const target = displayRestaurants.find((item) => item.id === restaurantId);
    if (!target) return;

    const latInput = window.prompt(`为「${target.name}」输入纬度（例如 30.5928）`, String(target.latitude));
    if (latInput === null) return;
    const lngInput = window.prompt(`为「${target.name}」输入经度（例如 114.3055）`, String(target.longitude));
    if (lngInput === null) return;

    const lat = Number(latInput.trim());
    const lng = Number(lngInput.trim());

    if (!Number.isFinite(lat) || lat < -90 || lat > 90 || !Number.isFinite(lng) || lng < -180 || lng > 180) {
      window.alert('经纬度格式不正确，请输入有效数字。');
      return;
    }

    setManualCoordsById((prev) => ({ ...prev, [restaurantId]: { lat, lng } }));
  };

  const handleMapLocate = (restaurantId: number) => {
    setPickingRestaurantId(restaurantId);
  };

  const handlePickCoordinates = (lat: number, lng: number) => {
    if (!pickingRestaurantId) return;
    setManualCoordsById((prev) => ({ ...prev, [pickingRestaurantId]: { lat, lng } }));
    setPickingRestaurantId(null);
  };

  const pickingRestaurant = displayRestaurants.find((item) => item.id === pickingRestaurantId) || null;

  return (
    <div className="space-y-6">
      {/* 搜索和筛选区域 */}
      <div className="bg-white rounded-xl p-4 shadow-sm border border-gray-100">
        <div className="flex flex-col md:flex-row gap-4">
          {/* 搜索框 */}
          <form onSubmit={handleSearch} className="flex-1 flex gap-2">
            <div className="flex-1 relative">
              <input
                type="text"
                placeholder="搜索店名..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100"
              />
              <button
                type="submit"
                className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-blue-500"
              >
                🔍
              </button>
            </div>
          </form>

          {/* 筛选下拉框 */}
          <div className="flex gap-2">
            <select
              value={sortBy}
              onChange={(e) => handleSortChange(e.target.value as SortOption)}
              className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100 bg-white"
            >
              <option value="distance">按距离</option>
              <option value="score">按评分</option>
              <option value="reviews">按评论数</option>
              <option value="recommended">智能推荐</option>
            </select>
            <button
              type="button"
              onClick={requestLocation}
              disabled={locationLoading}
              className="px-4 py-2 border border-blue-300 text-blue-600 rounded-lg hover:bg-blue-50 disabled:opacity-50"
            >
              {locationLoading ? '定位中...' : '📍 定位'}
            </button>
          </div>
        </div>
      </div>

      {/* 系统推荐 */}
      {recommended.length > 0 && (
        <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl p-4 border border-blue-100">
          <h3 className="text-lg font-bold text-gray-800 mb-3">🏆 系统推荐</h3>
          <div className="grid grid-cols-1 md:grid-cols-5 gap-3">
            {displayRecommended.map((r, index) => (
              <a
                key={r.id}
                href={`/restaurant/${r.id}`}
                className="bg-white rounded-lg p-3 shadow-sm hover:shadow-md transition-shadow border border-gray-100"
              >
                <div className="flex items-center gap-2 mb-1">
                  <span className={`text-xs font-bold px-2 py-0.5 rounded ${
                    index === 0 ? 'bg-yellow-100 text-yellow-700' :
                    index === 1 ? 'bg-gray-100 text-gray-700' :
                    index === 2 ? 'bg-orange-100 text-orange-700' :
                    'bg-blue-100 text-blue-700'
                  }`}>
                    TOP{index + 1}
                  </span>
                  <span className="font-medium text-gray-800 truncate flex-1">{r.name}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-blue-500 font-bold">⭐ {r.avg_score.toFixed(1)}</span>
                  <span className="text-gray-400 text-xs">
                    {r.distance ? r.distance.toFixed(1) + 'km' : ''}
                  </span>
                </div>
              </a>
            ))}
          </div>
        </div>
      )}

      {pickingRestaurant && (
        <div className="flex items-center justify-between rounded-lg border border-blue-200 bg-blue-50 px-3 py-2 text-sm text-blue-700">
          <span>正在为「{pickingRestaurant.name}」地图点选坐标，请在右侧地图单击。</span>
          <button type="button" onClick={() => setPickingRestaurantId(null)} className="text-xs font-medium hover:underline">
            取消
          </button>
        </div>
      )}

      {locationLoading && (
        <div className="rounded-lg bg-blue-50 px-3 py-2 text-sm text-blue-700">
          正在获取当前位置...
        </div>
      )}

      {/* 餐厅列表 */}
      <div>
        <div className="flex items-center justify-between mb-2">
          <h2 className="text-xl font-bold text-gray-800">
            餐厅列表
            <span className="text-gray-400 text-sm font-normal ml-2">({restaurants.length} 家)</span>
          </h2>
          <div className="flex items-center gap-3">
            {typeof accuracy === 'number' && (
              <span className="text-xs text-gray-500">定位精度: 约 {Math.round(accuracy)} 米</span>
            )}
            <button
              type="button"
              onClick={() => setShowRestaurantForm(!showRestaurantForm)}
              className="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 font-medium"
            >
              + 新增店铺
            </button>
          </div>
        </div>

        {displayRestaurants.length > 0 ? (
          <div className="grid gap-5 lg:grid-cols-12">
            <section className="grid gap-4 md:grid-cols-2 lg:col-span-8">
              {displayRestaurants.map((r) => (
                <RestaurantCard
                  key={r.id}
                  restaurant={r}
                  onManualLocate={handleManualLocate}
                  onMapLocate={handleMapLocate}
                />
              ))}
            </section>
            <aside className="lg:col-span-4">
              <MapPanel
                coords={coords}
                accuracy={accuracy}
                restaurants={displayRestaurants.slice(0, 20)}
                pickingRestaurantName={pickingRestaurant?.name ?? null}
                onPickCoordinates={handlePickCoordinates}
              />
            </aside>
          </div>
        ) : (
          <div className="col-span-2 text-center py-10 text-gray-400">
            {search ? '没有找到匹配的餐厅' : '暂无数据'}
          </div>
        )}
      </div>

      {/* 悬浮弹窗 - 新增店铺 */}
      {showRestaurantForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl p-6 w-full max-w-lg shadow-xl relative">
            <button
              type="button"
              onClick={() => {
                setShowRestaurantForm(false);
                setIsSelectingLocationForNew(false);
                setNewRestaurantLocation(null);
              }}
              className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 text-2xl z-10"
            >
              ×
            </button>

            {isSelectingLocationForNew ? (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold text-[#202124]">在地图上选择位置</h3>
                <p className="text-sm text-[#5f6368]">点击地图上的位置来设置新店铺的位置</p>
                <div className="h-[400px] rounded-xl overflow-hidden border border-[#dadce0]">
                  <MapContainer
                    center={coords ? [coords.lat, coords.lng] : [39.9042, 116.4074]}
                    zoom={13}
                    scrollWheelZoom
                    className="h-full w-full"
                  >
                    <TileLayer
                      attribution="&copy; 高德地图"
                      url="https://webrd0{s}.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&x={x}&y={y}&z={z}"
                      subdomains={['1', '2', '3', '4']}
                    />
                    <MapClickHandler onMapClick={(lat, lng) => {
                      setNewRestaurantLocation({ lat, lng });
                      setIsSelectingLocationForNew(false);
                    }} />
                  </MapContainer>
                </div>
                <button
                  type="button"
                  onClick={() => setIsSelectingLocationForNew(false)}
                  className="w-full rounded-xl border border-[#dadce0] py-3 font-medium text-[#5f6368] hover:bg-[#f1f3f4]"
                >
                  取消
                </button>
              </div>
            ) : (
              <RestaurantForm
                coords={coords}
                onSelectLocation={() => setIsSelectingLocationForNew(true)}
                selectedLocation={newRestaurantLocation}
                onSuccess={() => {
                  setShowRestaurantForm(false);
                  setNewRestaurantLocation(null);
                  fetchRestaurants(search, sortBy);
                }}
                onCancel={() => {
                  setShowRestaurantForm(false);
                  setNewRestaurantLocation(null);
                }}
              />
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default Home;
