// src/types/index.ts

export type Restaurant = {
  id: number;
  name: string;
  latitude: number;
  longitude: number;
  avg_score: number;
  review_count: number;
  final_score?: number; 
};

export type Rating = {
  id?: number;
  user_id: number;
  restaurant_id: number;
  restaurant_name: string;
  stars: number;
  comment: string;
  created_at?: string;
};
