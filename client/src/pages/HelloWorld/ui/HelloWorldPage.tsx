"use client";

import { useVisitCount } from "@/pages/HelloWorld/model/useVisitCount";
import { FireWorks } from "@/widgets/FireWorks";

function getOrdinalSuffix(num: number): string {
	if (num === 0) return "th";
	const lastDigit = num % 10;
	const lastTwoDigits = num % 100;

	if (lastTwoDigits >= 11 && lastTwoDigits <= 13) {
		return "th";
	}

	switch (lastDigit) {
		case 1:
			return "st";
		case 2:
			return "nd";
		case 3:
			return "rd";
		default:
			return "th";
	}
}

export function HelloWorldPage() {
	const count = useVisitCount();

	return (
		<div className="min-h-screen p-8 flex flex-col items-center gap-8">
			<h1 className="text-4xl font-bold">
				Hello world ({count}
				{getOrdinalSuffix(count)} time)
			</h1>
			<FireWorks />
		</div>
	);
}
