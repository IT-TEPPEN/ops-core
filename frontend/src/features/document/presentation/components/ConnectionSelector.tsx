import { useNavigate } from "react-router-dom";

/**
 * ConnectionSelector Component
 *
 * OAuth接続の選択コンポーネント。
 * 既存の接続一覧を表示し、選択できる。
 * 新しい接続を追加するボタンも提供する。
 *
 * 責任:
 * - 接続一覧の表示
 * - 接続の選択
 * - 新しい接続追加のトリガー
 */
export function ConnectionSelector() {
  const navigate = useNavigate();
  const handleAddConnection = () => {
    navigate("/oauth/authorize");
  };

  return (
    <div className="flex items-center justify-between">
      <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300">
        Connection
      </h3>
      <button
        onClick={handleAddConnection}
        className="text-xs px-2 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
      >
        + Add
      </button>
    </div>
  );
}
