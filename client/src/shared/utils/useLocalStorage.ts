import { useState, useEffect, useCallback } from "react";

/**
 * LocalStorageを使用した状態管理フック
 * @param key - LocalStorageのキー
 * @param initialValue - 初期値
 * @returns [value, setValue] - 値とセッター関数のタプル
 */
export function useLocalStorage<T>(
	key: string,
	initialValue: T,
): [T, (value: T | ((prev: T) => T)) => void] {
	// LocalStorageから初期値を取得
	const [storedValue, setStoredValue] = useState<T>(() => {
		try {
			const item = localStorage.getItem(key);
			return item ? (JSON.parse(item) as T) : initialValue;
		} catch (error) {
			console.warn(`Error reading localStorage key "${key}":`, error);
			return initialValue;
		}
	});

	// 値をLocalStorageに保存する関数
	const setValue = useCallback(
		(value: T | ((prev: T) => T)) => {
			try {
				// setStoredValueを使って最新の値を取得
				setStoredValue((currentValue) => {
					const valueToStore =
						value instanceof Function ? value(currentValue) : value;
					localStorage.setItem(key, JSON.stringify(valueToStore));
					return valueToStore;
				});
			} catch (error) {
				console.warn(`Error setting localStorage key "${key}":`, error);
			}
		},
		[key],
	);

	// 他のタブでの変更を監視
	useEffect(() => {
		const handleStorageChange = (e: StorageEvent) => {
			if (e.key === key && e.newValue) {
				try {
					setStoredValue(JSON.parse(e.newValue) as T);
				} catch (error) {
					console.warn(
						`Error parsing localStorage value for key "${key}":`,
						error,
					);
				}
			}
		};

		window.addEventListener("storage", handleStorageChange);
		return () => window.removeEventListener("storage", handleStorageChange);
	}, [key]);

	return [storedValue, setValue];
}
