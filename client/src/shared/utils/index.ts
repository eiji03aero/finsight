export const cn = (...cns: (string | undefined)[]) =>
	cns.filter((c) => !!c).join(" ");

export { useLocalStorage } from "./useLocalStorage";
