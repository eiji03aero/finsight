import { useEffect } from "react";
import { useLocalStorage } from "@/shared/utils";

const STORAGE_KEY = "hello_world_visit_count";

export function useVisitCount() {
	const [count, setCount] = useLocalStorage<number>(STORAGE_KEY, 0);

	// 初回レンダリング時にカウントをインクリメント
	useEffect(() => {
		setCount((prev) => prev + 1);
	}, [setCount]);

	return count;
}
