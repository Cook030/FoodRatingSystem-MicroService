// src/pages/Detail.tsx
import React from 'react';
import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { restaurantApi } from '../api/index.js';
import { Rating } from '../api/index.js';
import { RatingForm } from '../components/RatingForm.jsx';

interface Restaurant {
  id: number;
  name: string;
  latitude: number;
  longitude: number;
  avg_score: number;
  review_count: number;
  final_score?: number;
}

const Detail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [restaurant, setRestaurant] = useState<Restaurant | null>(null);
  const [ratings, setRatings] = useState<Rating[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const loadDetail = async () => {
    if (!id) return;
    try {
      setLoading(true);
      setError(null);
      const res = await restaurantApi.getDetail(Number(id));
      setRestaurant(res.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取餐厅详情失败');
    } finally {
      setLoading(false);
    }
  };

  const loadRatings = async () => {
    if (!id) return;
    try {
      const res = await restaurantApi.getRatings(Number(id));
      setRatings(res.data || []);
    } catch (err) {
      console.error('获取评论失败:', err);
    }
  };

  useEffect(() => {
    loadDetail();
    loadRatings();
  }, [id]);

  const handleRatingSuccess = () => {
    loadDetail();
    loadRatings();
    setSuccessMessage('评价提交成功！🎉');
    setTimeout(() => setSuccessMessage(null), 3000);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[200px]">
        <div className="animate-pulse text-gray-400">加载中...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-10">
        <div className="text-red-500 mb-4">{error}</div>
        <button
          onClick={() => navigate(-1)}
          className="text-blue-600 hover:underline"
        >
          返回上一页
        </button>
      </div>
    );
  }

  if (!restaurant) {
    return (
      <div className="text-center py-10 text-gray-400">
        找不到该餐厅
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {successMessage && (
        <div className="fixed top-4 left-1/2 -translate-x-1/2 z-50 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg animate-fade-in">
          {successMessage}
        </div>
      )}
      
      <button 
        onClick={() => navigate(-1)} 
        className="text-gray-500 hover:text-blue-600 text-sm font-medium transition-colors"
      >
        ← 返回列表
      </button>
      
      <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
        <h1 className="text-3xl font-bold text-gray-900">{restaurant.name}</h1>
        <div className="flex gap-6 mt-4 items-center">
          <div className="text-center">
            <div className="text-gray-400 text-xs mb-1">综合评分</div>
            <div className="text-2xl font-black text-yellow-500">{restaurant.avg_score.toFixed(1)}</div>
          </div>
          <div className="w-px h-8 bg-gray-200"></div>
          <div className="text-center">
            <div className="text-gray-400 text-xs mb-1">评价人数</div>
            <div className="text-xl font-bold text-gray-700">{restaurant.review_count}</div>
          </div>
        </div>
      </div>

      <RatingForm 
        restaurantId={restaurant.id} 
        onSuccess={handleRatingSuccess} 
      />

      <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
        <h2 className="text-xl font-bold text-gray-800 mb-4">
          用户评价 
          <span className="text-gray-400 text-sm font-normal ml-2">({ratings.length} 条)</span>
        </h2>
        
        {ratings.length > 0 ? (
          <div className="space-y-4">
            {ratings.map((rating, index) => (
              <div key={rating.id ?? index} className="border-b border-gray-100 pb-4 last:border-b-0 last:pb-0">
                <div className="flex items-center gap-3 mb-2">
                  <div className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center text-lg">
                    👤
                  </div>
                  <div className="flex-1">
                    <div className="font-medium text-gray-800">
                      {rating.user_name || rating.user_id || '匿名用户'}
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-blue-500 font-bold">
                        {'★'.repeat(Math.round(rating.stars))}
                        {'☆'.repeat(5 - Math.round(rating.stars))}
                      </span>
                      <span className="text-yellow-500 font-bold">{rating.stars.toFixed(1)}</span>
                    </div>
                  </div>
                  {rating.created_at && (
                    <span className="text-xs text-gray-400">
                      {new Date(rating.created_at).toLocaleDateString('zh-CN')}
                    </span>
                  )}
                </div>
                {rating.comment && (
                  <p className="text-gray-600 text-sm ml-13">{rating.comment}</p>
                )}
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-8 text-gray-400">
            暂无评价，快来抢沙发吧！
          </div>
        )}
      </div>
    </div>
  );
};

export default Detail;