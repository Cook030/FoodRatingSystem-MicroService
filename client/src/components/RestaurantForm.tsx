import { useState } from 'react';
import { restaurantApi } from '../api';
import React from 'react';

interface Props {
  onSuccess: () => void;
  onCancel?: () => void;
  coords: { lat: number; lng: number } | null;
  onSelectLocation: () => void;
  selectedLocation: { lat: number; lng: number } | null;
}

export const RestaurantForm = ({ onSuccess, onCancel, coords, onSelectLocation, selectedLocation }: Props) => {
  const [name, setName] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isSubmitting) return;

    if (!name) {
      setMessage({ type: 'error', text: '请填写店铺名称' });
      return;
    }

    if (!selectedLocation) {
      setMessage({ type: 'error', text: '请先在地图上选择店铺位置' });
      return;
    }

    try {
      setIsSubmitting(true);
      setMessage(null);

      await restaurantApi.create({
        name,
        latitude: selectedLocation.lat,
        longitude: selectedLocation.lng,
      });

      setName('');
      setMessage({ type: 'success', text: '店铺创建成功！🎉' });
      onSuccess();

      setTimeout(() => setMessage(null), 3000);
    } catch (err) {
      setMessage({
        type: 'error',
        text: err instanceof Error ? err.message : '创建失败，请检查网络'
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const displayCoords = selectedLocation
    ? `${selectedLocation.lat.toFixed(6)}, ${selectedLocation.lng.toFixed(6)}`
    : coords
    ? `${coords.lat.toFixed(6)}, ${coords.lng.toFixed(6)}`
    : '未选择';

  return (
    <>
      <form onSubmit={handleSubmit} className="space-y-4">
        <h3 className="text-lg font-semibold text-[#202124]">创建新店铺</h3>

        {message && (
          <div className={`rounded-xl p-3 text-sm ${
            message.type === 'success'
              ? 'bg-green-50 text-green-700'
              : 'bg-red-50 text-red-700'
          }`}>
            {message.text}
          </div>
        )}

        <div>
          <label className="mb-2 block text-sm font-medium text-[#3c4043]">店铺名称</label>
          <input
            type="text"
            className="w-full resize-none rounded-xl border border-[#dadce0] p-3 outline-none transition-all focus:border-[#1a73e8] focus:ring-2 focus:ring-[#e8f0fe]"
            placeholder="请输入店铺名称"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>

        <div>
          <label className="mb-2 block text-sm font-medium text-[#3c4043]">店铺位置</label>
          <div className="flex items-center gap-2">
            <div className="flex-1 rounded-xl border border-[#dadce0] p-3 text-sm text-[#5f6368] bg-gray-50">
              {displayCoords}
            </div>
            <button
              type="button"
              onClick={onSelectLocation}
              className="rounded-xl border border-[#1a73e8] px-4 py-3 text-sm font-medium text-[#1a73e8] hover:bg-[#e8f0fe]"
            >
              {selectedLocation ? '重新选择' : '选择位置'}
            </button>
          </div>
        </div>

        <div className="flex gap-3 pt-2">
          <button
            type="submit"
            disabled={isSubmitting}
            className={`flex-1 rounded-xl py-3 font-medium text-white transition-colors ${
              isSubmitting
                ? 'cursor-not-allowed bg-[#8ab4f8]'
                : 'bg-[#1a73e8] hover:bg-[#1765cc] active:bg-[#1558b0]'
            }`}
          >
            {isSubmitting ? '正在创建...' : '创建店铺'}
          </button>
          {onCancel && (
            <button
              type="button"
              onClick={onCancel}
              className="flex-1 rounded-xl border border-[#dadce0] py-3 font-medium text-[#5f6368] transition-colors hover:bg-[#f1f3f4]"
            >
              取消
            </button>
          )}
        </div>
      </form>

      {message?.type === 'success' && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[60]">
          <div className="bg-white rounded-2xl p-8 w-full max-w-sm shadow-xl text-center">
            <div className="text-6xl mb-4">🎉</div>
            <h3 className="text-xl font-bold text-[#202124] mb-2">添加成功！</h3>
            <p className="text-[#5f6368] mb-6">店铺「{name}」已成功创建</p>
            <button
              type="button"
              onClick={() => {
                setMessage(null);
                onSuccess();
              }}
              className="w-full rounded-xl py-3 font-medium bg-[#1a73e8] text-white hover:bg-[#1765cc]"
            >
              确定
            </button>
          </div>
        </div>
      )}
    </>
  );
};