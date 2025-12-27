import { Link } from "react-router-dom";
import { RepositoryList } from "../features/repository/components/RepositoryList";

/**
 * リポジトリ一覧ページ
 * 登録済みリポジトリの一覧を表示する
 */
function RepositoryListPage() {
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Repository Management</h1>
        <Link
          to="/repositories/new"
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          Register New Repository
        </Link>
      </div>

      <RepositoryList />
    </div>
  );
}

export default RepositoryListPage;
