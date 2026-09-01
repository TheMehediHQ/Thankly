import { useState } from 'react';

interface Gratitude {
  id: string;
  content: string;
  gratitude_date: string;
}

export default function GratitudeEditor() {
  const [content, setContent] = useState('');
  const [gratitudes, setGratitudes] = useState<Gratitude[]>([]);
  const [loading, setLoading] = useState(false);

  const remaining = 3 - gratitudes.length;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim() || remaining <= 0) return;

    setLoading(true);
    try {
      const res = await fetch('/api/gratitudes', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({ content }),
      });

      if (res.ok) {
        const data = await res.json();
        setGratitudes([...gratitudes, data]);
        setContent('');
      }
    } catch (err) {
      console.error('Failed to create gratitude:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <form onSubmit={handleSubmit} className="space-y-4">
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          rows={3}
          maxLength={500}
          className="w-full px-4 py-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition resize-none"
          placeholder="I'm grateful for..."
          disabled={remaining <= 0}
        />
        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-500">{remaining} remaining today</span>
          <button
            type="submit"
            disabled={loading || remaining <= 0}
            className="px-6 py-2 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 transition disabled:opacity-50"
          >
            {loading ? 'Adding...' : 'Add Gratitude'}
          </button>
        </div>
      </form>
    </div>
  );
}
