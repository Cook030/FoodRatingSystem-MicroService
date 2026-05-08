// src/hooks/useLocation.ts
import { useState, useCallback } from 'react';

interface LocationState {
  coords: { lat: number; lng: number } | null;
  accuracy: number | null;
  loading: boolean;
  error: string | null;
}

export const useLocation = (): LocationState & { requestLocation: () => void } => {
  const [state, setState] = useState<LocationState>({
    coords: null,
    accuracy: null,
    loading: false,
    error: null,
  });

  const requestLocation = useCallback(() => {
    if (!navigator.geolocation) {
      setState({
        coords: null,
        accuracy: null,
        loading: false,
        error: '浏览器不支持定位功能',
      });
      return;
    }

    setState(prev => ({ ...prev, loading: true, error: null }));

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setState({
          coords: {
            lat: pos.coords.latitude,
            lng: pos.coords.longitude,
          },
          accuracy: pos.coords.accuracy ?? null,
          loading: false,
          error: null,
        });
      },
      (err) => {
        let errorMessage = '无法获取位置';
        switch (err.code) {
          case err.PERMISSION_DENIED:
            errorMessage = '用户拒绝定位权限';
            break;
          case err.POSITION_UNAVAILABLE:
            errorMessage = '定位信息不可用';
            break;
          case err.TIMEOUT:
            errorMessage = '定位请求超时';
            break;
        }
        setState({
          coords: null,
          accuracy: null,
          loading: false,
          error: errorMessage,
        });
      },
      {
        enableHighAccuracy: false,
        timeout: 12000,
        maximumAge: 0,
      }
    );
  }, []);

  return { ...state, requestLocation };
};
