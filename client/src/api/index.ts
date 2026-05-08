// src/api/index.ts
import axios, { AxiosError } from 'axios';

interface Restaurant {
  id: number;
  name: string;
  latitude: number;
  longitude: number;
  avg_score: number;
  review_count: number;
  final_score?: number; 
}

export interface Rating {
  id?: number;
  user_id: string;
  user_name?: string;
  restaurant_id: number;
  restaurant_name: string;
  stars: number;
  comment: string;
  created_at?: string;
}

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 5000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error: AxiosError<{ error?: string; message?: string }>) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    const message = error.response?.data?.error || error.response?.data?.message || error.message || '请求失败';
    console.error('API Error:', message);
    return Promise.reject(new Error(message));
  }
);

export const restaurantApi = {
  getRecommended: (lat: number, lng: number) => 
    api.get<Restaurant[]>('/restaurants/nearby', { params: { lat, lng } }),
  
  getDetail: (id: number) => 
    api.get<Restaurant>(`/restaurants/${id}`),

  getRatings: (id: number) => 
    api.get<Rating[]>(`/restaurants/${id}/ratings`),

  getAll: (lat?: number, lng?: number, search?: string, sort?: string) =>
    api.get<Restaurant[]>('/restaurants', { params: { lat, lng, search, sort } }),

  create:(data:{name:string,latitude:number,longitude:number})=>
    api.post<Restaurant>('/restaurants',data)
};

export interface User {
  user_id: string;
  user_name: string;
  created_at?: string;
}

export interface AuthResponse {
  message: string;
  user: User;
  token: string;
}

export const userApi = {
  register: (username: string, password: string) =>
    api.post<AuthResponse>('/user/register', { username, password }),

  login: (username: string, password: string) =>
    api.post<AuthResponse>('/user/login', { username, password }),
};

export const ratingApi = {
  submit: (data: { username: string; restaurant_id: number; restaurant_name: string; stars: number; comment: string }) => 
    api.post('/rating', data),
};

export default api;