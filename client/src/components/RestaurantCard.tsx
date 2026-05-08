// src/components/RestaurantCard.tsx
import { Link } from 'react-router-dom';
import React from 'react';

interface Restaurant {
  id: number;
  name: string;
  latitude: number;
  longitude: number;
  avg_score: number;
  review_count: number;
  final_score?: number;
}

interface Props {
  restaurant: Restaurant;
  onManualLocate: (restaurantId: number) => void;
  onMapLocate: (restaurantId: number) => void;
}

export const RestaurantCard = ({ restaurant, onManualLocate, onMapLocate }: Props) => {
  return (
    <div className="group rounded-2xl border border-[#e8eaed] bg-white px-4 py-3 shadow-sm transition hover:border-[#d2e3fc] hover:shadow-md">
      <div className="flex items-start justify-between gap-2">
        <div className="min-w-0 flex-1">
          <h2 className="truncate text-base font-semibold text-[#202124]">{restaurant.name}</h2>
          <p className="text-xs text-[#5f6368]">
            坐标: {restaurant.latitude.toFixed(4)}, {restaurant.longitude.toFixed(4)}
          </p>
        </div>
        <div className="shrink-0 text-right">
          <div className="text-lg font-bold text-[#1a73e8]">
            {restaurant.avg_score.toFixed(1)} <span className="text-xs font-medium text-[#5f6368]">分</span>
          </div>
          <div className="text-xs text-[#5f6368]">{restaurant.review_count} 条</div>
        </div>
      </div>
      
      <div className="mt-2 flex items-center justify-between gap-2 text-xs">
        {restaurant.final_score ? (
          <span className="rounded-full bg-[#e8f0fe] px-2 py-0.5 font-medium text-[#1a73e8] whitespace-nowrap">
            推荐指数: {restaurant.final_score.toFixed(1)}
          </span>
        ) : <div />}

        <div className="flex items-center gap-2 flex-wrap justify-end">
          <button
            type="button"
            onClick={() => onMapLocate(restaurant.id)}
            className="rounded-full border border-[#d2e3fc] bg-white px-2 py-0.5 font-medium text-[#1a73e8] hover:bg-[#e8f0fe] whitespace-nowrap"
          >
            地图点选
          </button>
          <button
            type="button"
            onClick={() => onManualLocate(restaurant.id)}
            className="rounded-full border border-[#d2e3fc] bg-white px-2 py-0.5 font-medium text-[#1a73e8] hover:bg-[#e8f0fe] whitespace-nowrap"
          >
            手动定位
          </button>
          <Link to={`/restaurant/${restaurant.id}`} className="font-medium text-[#1a73e8] group-hover:underline whitespace-nowrap">
            查看详情
          </Link>
        </div>
      </div>
    </div>
  );
};