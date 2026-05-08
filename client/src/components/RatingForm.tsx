// src/components/RatingForm.tsx
import { useState } from 'react';
import { ratingApi } from '../api';
import React from 'react';

interface Props {
  restaurantId: number;
  onSuccess: () => void;
}

export const RatingForm = ({ restaurantId, onSuccess }: Props) => {
  const [stars, setStars] = useState(0);
  const [hoverStars, setHoverStars] = useState(0);
  const [comment, setComment] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isSubmitting) return;

    const userData = localStorage.getItem('user');
    if (!userData) {
      setMessage({ type: 'error', text: '请先登录' });
      return;
    }

    const user = JSON.parse(userData);

    if (stars === 0) {
      setMessage({ type: 'error', text: '请先选择评分' });
      return;
    }

    try {
      setIsSubmitting(true);
      setMessage(null);
      
      await ratingApi.submit({
        username: user.user_name,
        restaurant_id: restaurantId,
        restaurant_name: "",
        stars,
        comment
      });
      
      setComment('');
      setStars(0);
      setMessage({ type: 'success', text: '评价提交成功！🎉' });
      onSuccess(); 
      
      setTimeout(() => setMessage(null), 3000);
    } catch (err) {
      setMessage({ 
        type: 'error', 
        text: err instanceof Error ? err.message : '提交评价失败，请检查网络'
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="mt-6 rounded-2xl border border-[#e8eaed] bg-white p-6 shadow-sm">
      <h3 className="mb-4 text-lg font-semibold text-[#202124]">写下你的真实评价</h3>
      
      {message && (
        <div className={`mb-4 rounded-xl p-3 text-sm ${
          message.type === 'success' 
            ? 'bg-green-50 text-green-700' 
            : 'bg-red-50 text-red-700'
        }`}>
          {message.text}
        </div>
      )}
      
      <div className="mb-4">
        <label className="mb-2 block text-sm font-medium text-[#3c4043]">打个分吧</label>
        <div className="flex gap-2">
          {[1, 2, 3, 4, 5].map((s) => (
            <button 
              key={s} 
              type="button"
              onMouseEnter={() => setHoverStars(s)}
              onMouseLeave={() => setHoverStars(0)}
              onClick={() => setStars(s)}
              className="text-3xl transition-transform hover:scale-110 focus:outline-none"
            >
              <span className={(hoverStars || stars) >= s ? 'text-[#f9ab00]' : 'text-[#dadce0]'}>
                ★
              </span>
            </button>
          ))}
        </div>
        {stars > 0 && (
          <div className="mt-2 text-sm text-[#5f6368]">
            {stars} 星评分
          </div>
        )}
      </div>

      <div className="mb-4">
        <label className="mb-2 block text-sm font-medium text-[#3c4043]">详细感受</label>
        <textarea 
          className="w-full resize-none rounded-xl border border-[#dadce0] p-3 outline-none transition-all focus:border-[#1a73e8] focus:ring-2 focus:ring-[#e8f0fe]"
          rows={4}
          placeholder="这家店的味道、环境、服务怎么样？"
          value={comment}
          onChange={(e) => setComment(e.target.value)}
          required
        />
      </div>

      <button 
        type="submit"
        disabled={isSubmitting}
        className={`w-full rounded-xl py-3 font-medium text-white transition-colors ${
          isSubmitting 
            ? 'cursor-not-allowed bg-[#8ab4f8]' 
            : 'bg-[#1a73e8] hover:bg-[#1765cc] active:bg-[#1558b0]'
        }`}
      >
        {isSubmitting ? '正在提交...' : '提交评价'}
      </button>
    </form>
  );
};